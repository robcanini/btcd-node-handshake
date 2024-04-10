package node

import (
	"net"
)

type TcpConnection struct {
	tcp net.Conn
}

// todo: add marshalling/unmarshalling

func NewTcpConnection(tcpConnection net.Conn) *TcpConnection {
	return &TcpConnection{
		tcp: tcpConnection,
	}
}

func (conn *TcpConnection) Dispose() (err error) {
	err = conn.tcp.Close()
	if err != nil {
		return
	}
	return
}

func (conn *TcpConnection) Read(bytes []byte) (err error) {
	_, err = conn.tcp.Read(bytes)
	if err != nil {
		return
	}
	return
}

func (conn *TcpConnection) Write(bytes []byte) (err error) {
	_, err = conn.tcp.Write(bytes)
	if err != nil {
		return
	}
	return
}
