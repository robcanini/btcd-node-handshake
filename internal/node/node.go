package node

import (
	"net"

	"github.com/robcanini/btcd-node-handshake/internal/message"
)

type HandshakeCode int

const (
	HDone  HandshakeCode = 0
	HError HandshakeCode = 1
)

type Node interface {
	Connect() (func(), error)
	StartHandshake(stopCh chan HandshakeCode, lastBlock uint64) (err error)
	SendVerAck() error
	onMessage(msg message.Message)
}

type Network uint32

type connection interface {
	dispose()
	read([]byte) (int, error)
	write([]byte) error
	tcpConn() net.Conn
}
