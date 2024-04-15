package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/robcanini/btcd-node-handshake/internal/config"
	"github.com/robcanini/btcd-node-handshake/internal/message"

	"github.com/rs/zerolog"
)

const (
	MainNet  Network = 0xd9b4bef9
	TestNet  Network = 0xdab5bffa
	TestNet3 Network = 0x0709110b
	SimNet   Network = 0x12141c16
)

type Btcd struct {
	log zerolog.Logger
	ctx context.Context
	cfg config.Node

	conn connection
}

func NewBtcdTcpClient(log zerolog.Logger, ctx context.Context, config config.Node) *Btcd {
	return &Btcd{
		log: log,
		ctx: ctx,
		cfg: config,
	}
}

func (b *Btcd) Connect() (disposeFun func(), err error) {
	dialConn, err := net.Dial("tcp", b.cfg.Address())
	if err != nil {
		return
	}
	b.log.Info().Str("node", b.cfg.Address()).Msg("connected to btcd node")
	b.conn = newTcpConnection(b.log, dialConn)
	disposeFun = b.close
	// notify client about incoming message from remote btcd tcp
	go b.listenInbound()
	return
}

func (b *Btcd) readFromRemote(dataCh chan []byte, errorCh chan error) {
	defer close(dataCh)
	for {
		headerData := make([]byte, message.HeaderSize)
		_, err := b.conn.read(headerData)
		if err != nil {
			errorCh <- err
			return
		}

		var command [message.CommandSize]byte
		var hdr message.MsgHeader
		if err = readElements(bytes.NewReader(headerData), &hdr.Magic, &command, &hdr.Length, &hdr.Checksum); err != nil {
			errorCh <- err
			return
		}

		payloadData := make([]byte, hdr.Length)
		_, err = io.ReadFull(b.conn.tcpConn(), payloadData)
		if err != nil {
			errorCh <- err
			return
		}
		messageData := append(headerData, payloadData...)
		dataCh <- messageData
	}
}

func (b *Btcd) listenInbound() {
	dataCh := make(chan []byte)
	errorCh := make(chan error)

	// start reading from btcd node
	go b.readFromRemote(dataCh, errorCh)

	for {
		select {
		// data
		case data := <-dataCh:
			b.log.Debug().Str("data", string(data)).Msg("received data")
			b.handleReadData(errorCh, data)
		// errors
		case err := <-errorCh:
			b.log.Error().Err(err).Msg("error reading data")
			return
		// connection timeout
		case <-time.After(b.cfg.ConnTimeout):
			b.log.Error().Msg("connection timed out")
			return
		}
	}
}

func (b *Btcd) handleReadData(errorCh chan error, data []byte) {
	var hdr message.MsgHeader
	var command [message.CommandSize]byte
	if err := readElements(bytes.NewReader(data), &hdr.Magic, &command, &hdr.Length, &hdr.Checksum); err != nil {
		errorCh <- err
	}
	hdr.Command = message.Cmd(bytes.TrimRight(command[:], "\x00"))
	payload := data[message.HeaderSize:]
	if len(payload) != int(hdr.Length) {
		errorCh <- errors.New("payload length mismatch")
	}
	b.onMessage(message.Message{
		Header:  hdr,
		Payload: payload,
	})
}

func readElements(r io.Reader, elements ...interface{}) (err error) {
	for _, element := range elements {
		if err = binary.Read(r, binary.LittleEndian, element); err != nil {
			return
		}
	}
	return
}

func (b *Btcd) close() {
	b.conn.dispose()
	b.conn = nil
}

func (b *Btcd) IsConnected() bool {
	return b.conn != nil
}

func (b *Btcd) SendVer() (err error) {
	btcdCfg := b.cfg.Btcd
	sourceAddr := message.NetAddress{
		IP:   net.ParseIP("127.0.0.1").To4(),
		Port: 8443,
	}
	targetAddr := message.NetAddress{
		IP:   net.ParseIP(b.cfg.Host).To4(),
		Port: b.cfg.Port,
	}
	msg := message.NewMsgVersion(
		uint32(MainNet),
		btcdCfg.ProtocolVersion,
		uint64(btcdCfg.Services),
		sourceAddr,
		targetAddr,
		btcdCfg.Agent,
		212672,
		!btcdCfg.RelayTx,
	)

	msgBytes, err := msg.ToBytes()
	if err != nil {
		return
	}
	err = b.conn.write(msgBytes)
	if err != nil {
		return
	}
	return
}

func (b *Btcd) SendVerAck() (err error) {
	msg := message.NewMsgVerAck(
		uint32(MainNet),
	)
	msgBytes, err := msg.ToBytes()
	if err != nil {
		return
	}
	err = b.conn.write(msgBytes)
	if err != nil {
		return
	}
	return
}

func (b *Btcd) onMessage(msg message.Message) {
	if len(msg.Header.Command) == 0 {
		b.log.Error().
			Msg("command is empty. discarding message")
	}
	switch msg.Header.Command {
	case message.CmdVersion:
		b.log.Info().Msg("version received")
	case message.CmdSendAddrV2:
		b.log.Info().Msg("sendaddrv2 received")
	case message.CmdVersionAck:
		b.log.Info().Msg("verack received")
	default:
		b.log.Warn().
			Str("command", string(msg.Header.Command)).
			Msg("unsupported command")
	}
}

func (b *Btcd) VerAck() (err error) {
	response := make([]byte, 200)
	_, err = b.conn.read(response)
	if err != nil {
		return
	}

	var version uint32
	var services uint64
	var timestamp uint32
	var addrNodeToConnect [4]byte
	var portNodeToConnect uint16

	payload := bytes.NewReader(response)
	_ = binary.Read(payload, binary.LittleEndian, &version)
	_ = binary.Read(payload, binary.LittleEndian, &services)
	_ = binary.Read(payload, binary.LittleEndian, &timestamp)
	_ = binary.Read(payload, binary.LittleEndian, &addrNodeToConnect)
	_ = binary.Read(payload, binary.BigEndian, &portNodeToConnect)

	fmt.Printf("Versione: %d\n", version)
	fmt.Printf("Servizi: %d\n", services)
	fmt.Printf("Timestamp: %d\n", timestamp)
	fmt.Printf("Indirizzo IP del nodo a cui ti stai connettendo: %s\n", net.IP(addrNodeToConnect[:]))
	fmt.Printf("Porta del nodo a cui ti stai connettendo: %d\n", portNodeToConnect)

	return
}
