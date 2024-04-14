package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"net"

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
	b.log.Debug().Msg(toHexString(msgBytes))
	err = b.conn.write(msgBytes)
	if err != nil {
		return
	}
	return
}

func toHexString(data []byte) (formatted string) {
	str := hex.EncodeToString(data)
	var buffer bytes.Buffer
	buffer.WriteString("\nMessage header\n")
	buffer.WriteString(str[0:8])
	buffer.WriteString("\n")
	buffer.WriteString(str[8:32])
	buffer.WriteString("\n")
	buffer.WriteString(str[32:40])
	buffer.WriteString("\n")
	buffer.WriteString(str[40:48])
	buffer.WriteString("\n")
	buffer.WriteString("\nVersion message \n")
	buffer.WriteString(str[48:56])
	buffer.WriteString("\n")
	buffer.WriteString(str[56:72])
	buffer.WriteString("\n")
	buffer.WriteString(str[72:88])
	buffer.WriteString("\n")
	buffer.WriteString(str[88:140])
	buffer.WriteString("\n")
	buffer.WriteString(str[140:192])
	buffer.WriteString("\n")
	buffer.WriteString(str[193:209])
	buffer.WriteString("\n")
	/*
		buffer.WriteString(str[209:241])
		buffer.WriteString("\n")
		buffer.WriteString(str[241:250])
		buffer.WriteString("\n")

	*/
	formatted = buffer.String()
	return
}

func (b *Btcd) VerAck() error {
	//TODO implement me
	panic("implement me")
}
