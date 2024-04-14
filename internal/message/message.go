package message

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
)

type Cmd string

const (
	VersionCmd    Cmd = "version"
	VersionAckCmd Cmd = "verack"
)

type MsgHeader struct {
	magic    uint32 // 4 bytes
	command  Cmd    // 12 bytes
	length   uint32 // 4 bytes
	checksum []byte // 4 bytes
}

func newMsgHeader(magic uint32, command Cmd, length uint32, checksum []byte) MsgHeader {
	return MsgHeader{
		magic:    magic,
		command:  command,
		length:   length,
		checksum: checksum,
	}
}

func (msg *MsgHeader) toBytes() (buf []byte) {
	message := new(bytes.Buffer)
	_ = binary.Write(message, binary.LittleEndian, msg.magic)
	command := []byte(msg.command)
	command = append(command, make([]byte, 12-len(command))...)
	message.Write(command)
	_ = binary.Write(message, binary.LittleEndian, msg.length)
	message.Write(msg.checksum[0:4])
	buf = message.Bytes()
	return
}

func computeChecksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:]
}

func randUint64() uint64 {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf) // no need to check error
	return binary.LittleEndian.Uint64(buf)
}
