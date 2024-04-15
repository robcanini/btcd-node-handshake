package node

import (
	"net"
	"testing"

	"github.com/rs/zerolog"
)

const (
	testTcpServer = "tcpbin.com:4242"
)

func TestTcpConnectionDispose(t *testing.T) {
	testConn, err := net.Dial("tcp", testTcpServer)
	if err != nil {
		return
	}
	defer testConn.Close()
	// Create a tcpConnection with a mock logger and the mock net.Conn
	conn := newTcpConnection(zerolog.Nop(), testConn)
	// Call dispose
	conn.dispose()
	// Check if the TCP connection is closed
	if err := testConn.Close(); err == nil {
		t.Error("Expected error while closing TCP connection, got nil")
	}
}

func TestTcpConnectionWrite(t *testing.T) {
	testConn, err := net.Dial("tcp", testTcpServer)
	if err != nil {
		return
	}
	defer testConn.Close()
	// Create a tcpConnection with a mock logger and the mock net.Conn
	conn := newTcpConnection(zerolog.Nop(), testConn)
	// Define a test message
	testMessage := []byte("Hello server!")
	// Write to the connection
	err = conn.write(testMessage)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
