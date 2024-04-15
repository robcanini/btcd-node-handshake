package node

import (
	"github.com/robcanini/btcd-node-handshake/internal/message"
	"net"
)

type Node interface {
	Connect() (func(), error)
	IsConnected() bool
	SendVer() error
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
