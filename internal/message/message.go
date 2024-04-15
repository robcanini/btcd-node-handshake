package message

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
)

type Cmd string

const (
	CmdVersion    Cmd = "version"
	CmdVersionAck Cmd = "verack"
	CmdSendAddrV2 Cmd = "sendaddrv2"
)

const (
	// HeaderSize is a fixed size of: Bitcoin network (Magic) 4 bytes + command 12 bytes + payload Length 4 bytes + Checksum 4 bytes.
	HeaderSize = 24
	// CommandSize is the fixed size of all commands in the common bitcoin message header.
	CommandSize = 12
)

type (
	Message struct {
		Header  MsgHeader
		Payload []byte
	}
	MsgHeader struct {
		Magic    uint32 // 4 bytes
		Command  Cmd    // 12 bytes
		Length   uint32 // 4 bytes
		Checksum []byte // 4 bytes
	}
)

func newMsgHeader(magic uint32, command Cmd, length uint32, checksum []byte) MsgHeader {
	return MsgHeader{
		Magic:    magic,
		Command:  command,
		Length:   length,
		Checksum: checksum,
	}
}

func (msg *MsgHeader) toBytes() (buf []byte) {
	message := new(bytes.Buffer)
	_ = binary.Write(message, binary.LittleEndian, msg.Magic)
	command := []byte(msg.Command)
	command = append(command, make([]byte, 12-len(command))...)
	message.Write(command)
	_ = binary.Write(message, binary.LittleEndian, msg.Length)
	message.Write(msg.Checksum[0:4])
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
