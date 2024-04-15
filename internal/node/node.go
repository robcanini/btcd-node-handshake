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
	IsConnected() bool
	SendVer(callback func(message.Message)) error
	SendVerAck() error
	onMessage(msg message.Message)
	VerAck() error
}

type Network uint32

type connection interface {
	dispose()
	read([]byte) (int, error)
	write([]byte) error
	tcpConn() net.Conn
}
