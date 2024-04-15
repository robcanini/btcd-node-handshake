package message

import (
	"bytes"
	"testing"
)

const (
	testNetwork = 12345
)

func TestNewMsgVerAck(t *testing.T) {
	msg := NewMsgVerAck(testNetwork)

	if msg.Network != testNetwork {
		t.Errorf("NewMsgVerAck failed: expected Network %d, got %d", testNetwork, msg.Network)
	}
}

func TestMsgVerAckToBytes(t *testing.T) {
	msg := NewMsgVerAck(testNetwork)

	// Calculate expected checksum
	payload := []byte{}
	expectedChecksum := computeChecksum(payload)

	// Generate expected message bytes
	expectedHeader := newMsgHeader(testNetwork, CmdVersionAck, uint32(len(payload)), expectedChecksum)
	expectedBytes := append(expectedHeader.toBytes(), payload...)

	// Get bytes from message
	buf, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("MsgVerAck ToBytes failed: %v", err)
	}

	// Compare bytes
	if !bytes.Equal(buf, expectedBytes) {
		t.Errorf("MsgVerAck ToBytes failed: expected %v, got %v", expectedBytes, buf)
	}
}
