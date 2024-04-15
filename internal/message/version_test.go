package message

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestNetAddressToBytes(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	netAddr := NetAddress{
		Timestamp: time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
		Services:  12345,
		IP:        ip,
		Port:      8333,
	}

	expectedBytes := []byte{
		0x20, 0x8d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0x7f, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	}
	buf := netAddr.ToBytes()

	if !bytes.Equal(buf, expectedBytes) {
		t.Errorf("NetAddress ToBytes failed: expected %x, got %x", expectedBytes, buf)
	}
}
