package node

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"time"

	"github.com/robcanini/btcd-node-handshake/internal/config"

	"github.com/rs/zerolog"
)

type Btcd struct {
	log zerolog.Logger
	ctx context.Context
	cfg config.Btcd

	conn connection
}

func NewBtcdTcpClient(log zerolog.Logger, ctx context.Context, config config.Btcd) *Btcd {
	return &Btcd{
		log: log,
		ctx: ctx,
		cfg: config,
	}
}

func (b *Btcd) Connect() (disposeFun func(), err error) {
	dialConn, err := net.Dial("tcp", b.cfg.Node)
	if err != nil {
		return
	}
	b.log.Info().Str("node", b.cfg.Node).Msg("connected to btcd node")
	b.conn = newTcpConnection(b.log, dialConn)
	disposeFun = b.close
	return
}

func (b *Btcd) close() {
	b.conn.dispose()
	b.conn = nil
}

func (b *Btcd) IsConnected() bool {
	return b.conn != nil
}

func (b *Btcd) SendVersion() (err error) {
	// TODO: poc

	version := uint32(70016)
	timestamp := time.Now().UnixMilli()
	services := uint64(1)
	addrYourNode := "127.0.0.1"
	portYourNode := uint16(8333)
	addrNodeToConnect := "127.0.0.1"
	portNodeToConnect := uint16(8333)
	nonce := uint64(123456)
	userAgent := "/Satoshi:26.0.0/"
	startBlockHeight := uint32(838611)
	relay := uint8(1)

	payload := new(bytes.Buffer)
	binary.Write(payload, binary.LittleEndian, version)
	binary.Write(payload, binary.LittleEndian, timestamp)
	binary.Write(payload, binary.LittleEndian, services)
	binary.Write(payload, binary.LittleEndian, net.ParseIP(addrNodeToConnect).To4())
	binary.Write(payload, binary.BigEndian, portNodeToConnect)
	binary.Write(payload, binary.LittleEndian, services)
	binary.Write(payload, binary.LittleEndian, net.ParseIP(addrYourNode).To4())
	binary.Write(payload, binary.BigEndian, portYourNode)
	binary.Write(payload, binary.LittleEndian, nonce)
	binary.Write(payload, binary.LittleEndian, uint8(len(userAgent)))
	payload.Write([]byte(userAgent))
	binary.Write(payload, binary.LittleEndian, startBlockHeight)
	binary.Write(payload, binary.LittleEndian, relay)

	message := new(bytes.Buffer)
	binary.Write(message, binary.LittleEndian, uint32(0xf9beb4d9))
	command := []byte("version")
	command = append(command, make([]byte, 12-len(command))...)
	message.Write(command)
	binary.Write(message, binary.LittleEndian, uint32(payload.Len()))
	message.Write(payload.Bytes())

	err = b.conn.write(message.Bytes())
	if err != nil {
		return
	}
	b.log.Debug().Msg("version command sent to btcd node")
	return
}

func (b *Btcd) VerAck() error {
	//TODO implement me
	panic("implement me")
}
