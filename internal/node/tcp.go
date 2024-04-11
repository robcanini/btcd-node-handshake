package node

import (
	"net"

	"github.com/rs/zerolog"
)

type tcpConnection struct {
	log zerolog.Logger

	tcp net.Conn
}

func newTcpConnection(log zerolog.Logger, tcpConn net.Conn) *tcpConnection {
	return &tcpConnection{
		log: log,
		tcp: tcpConn,
	}
}

func (conn *tcpConnection) dispose() {
	err := conn.tcp.Close()
	if err != nil {
		conn.log.Warn().Err(err).Msg("error while closing tcp connection")
		return
	}
	conn.log.Debug().Msg("tcp connection closed")
}

func (conn *tcpConnection) read(bytes []byte) (err error) {
	_, err = conn.tcp.Read(bytes)
	if err != nil {
		return
	}
	return
}

func (conn *tcpConnection) write(bytes []byte) (err error) {
	_, err = conn.tcp.Write(bytes)
	if err != nil {
		return
	}
	return
}
