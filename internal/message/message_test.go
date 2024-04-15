package message

import (
	"bytes"
	"testing"
)

func TestNewMsgHeader(t *testing.T) {
	magic := uint32(123456)
	command := CmdVersion
	length := uint32(100)
	checksum := make([]byte, 4)
	header := newMsgHeader(magic, command, length, checksum)

	if header.Magic != magic {
		t.Errorf("NewMsgHeader failed: expected Magic %d, got %d", magic, header.Magic)
	}
	if header.Command != command {
		t.Errorf("NewMsgHeader failed: expected Command %s, got %s", command, header.Command)
	}
	if header.Length != length {
		t.Errorf("NewMsgHeader failed: expected Length %d, got %d", length, header.Length)
	}
	if !bytes.Equal(header.Checksum, checksum) {
		t.Errorf("NewMsgHeader failed: expected Checksum %v, got %v", checksum, header.Checksum)
	}
}

func TestMsgHeaderToBytes(t *testing.T) {
	header := MsgHeader{
		Magic:    123456,
		Command:  "version",
		Length:   100,
		Checksum: []byte{0x00, 0x01, 0x02, 0x03},
	}

	expectedBytes := []byte{
		0x40, 0xe2, 0x01, 0x00, // Magic
		0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, // Command
		0x64, 0x00, 0x00, 0x00, // Length
		0x00, 0x01, 0x02, 0x03, // Checksum
	}

	buf := header.toBytes()

	if !bytes.Equal(buf, expectedBytes) {
		t.Errorf("MsgHeader ToBytes failed: expected %x, got %x", expectedBytes, buf)
	}
}

func TestComputeChecksum(t *testing.T) {
	payload := []byte{0x01, 0x02, 0x03, 0x04}
	expectedChecksum := []byte{
		0x8d, 0xe4, 0x72, 0xe2, 0x39, 0x96, 0x10, 0xba,
		0xaa, 0x7f, 0x84, 0x84, 0x05, 0x47, 0xcd, 0x40,
		0x94, 0x34, 0xe3, 0x1f, 0x5d, 0x3b, 0xd7, 0x1e,
		0x4d, 0x94, 0x7f, 0x28, 0x38, 0x74, 0xf9, 0xc0,
	}
	checksum := computeChecksum(payload)

	if !bytes.Equal(checksum, expectedChecksum) {
		t.Errorf("ComputeChecksum failed: expected %x, got %x", expectedChecksum, checksum)
	}
}

func TestRandUint64(t *testing.T) {
	value := randUint64()
	if value == 0 {
		t.Errorf("RandUint64 failed: expected non-zero value, got %d", value)
	}
}
