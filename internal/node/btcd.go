package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/robcanini/btcd-node-handshake/internal/config"
	"github.com/robcanini/btcd-node-handshake/internal/message"

	"github.com/rs/zerolog"
)

type Btcd struct {
	log  zerolog.Logger
	ctx  context.Context
	cfg  config.Node
	host string

	conn   connection
	stopCh chan HandshakeCode
	stop   bool

	paramsMux sync.Mutex
}

func NewBtcdTcpClient(log zerolog.Logger, ctx context.Context, config config.Node, host string) *Btcd {
	return &Btcd{
		log:  log,
		ctx:  ctx,
		cfg:  config,
		host: host,
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
	defer close(errorCh)
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
			if b.stop {
				// tcp channel has been closed
				return
			}
			b.log.Error().Err(err).Msg("error reading data")
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
	if b.conn != nil {
		b.conn.dispose()
		b.conn = nil
	}
}

func (b *Btcd) StartHandshake(stopCh chan HandshakeCode, lastBlock uint64) (err error) {
	btcdCfg := b.cfg.Btcd
	sourceAddr := message.NetAddress{
		IP:   net.ParseIP(b.host).To4(),
		Port: 8443,
	}
	targetAddr := message.NetAddress{
		IP:   net.ParseIP(b.cfg.Host).To4(),
		Port: b.cfg.Port,
	}
	msg := message.NewMsgVersion(
		b.cfg.Btcd.Network,
		btcdCfg.ProtocolVersion,
		uint64(btcdCfg.Services),
		sourceAddr,
		targetAddr,
		btcdCfg.Agent,
		int32(lastBlock),
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
	b.stopCh = stopCh
	return
}

func (b *Btcd) SendVerAck() (err error) {
	msg := message.NewMsgVerAck(
		b.cfg.Btcd.Network,
	)
	msgBytes, err := msg.ToBytes()
	if err != nil {
		b.closeHandshake(HError)
		return
	}
	err = b.conn.write(msgBytes)
	if err != nil {
		b.closeHandshake(HError)
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
		b.updateNegotiationParams(msg)
	case message.CmdSendAddrV2:
		b.updateNegotiationParams(msg)
	case message.CmdVersionAck:
		b.versionAcknowledgeHandler(msg)
	default:
		b.log.Warn().
			Str("command", string(msg.Header.Command)).
			Msg("unsupported command")
	}
}

func (b *Btcd) updateNegotiationParams(_ message.Message) {
	b.paramsMux.Lock()
	defer b.paramsMux.Unlock()
	b.log.Info().
		Msg("updating handshake negotiation params accordingly to node state")

	// todo: we can safely update the negotiation client params accordingly to the btcd node
}

func (b *Btcd) versionAcknowledgeHandler(_ message.Message) {
	b.log.Info().
		Msg("received version ack from node. sending our ack")
	err := b.SendVerAck()
	if err != nil {
		b.closeHandshake(HError)
		return
	}
	b.closeHandshake(HDone)
}

func (b *Btcd) closeHandshake(code HandshakeCode) {
	b.stop = true
	b.stopCh <- code
}
