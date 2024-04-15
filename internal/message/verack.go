package message

import (
	"bytes"
)

type MsgVerAck struct {
	Network uint32
}

func NewMsgVerAck(network uint32) *MsgVerAck {
	return &MsgVerAck{
		Network: network,
	}
}

func (msg *MsgVerAck) ToBytes() (buf []byte, err error) {
	payload := new(bytes.Buffer)
	message := new(bytes.Buffer)

	// header
	header := newMsgHeader(msg.Network, CmdVersionAck, uint32(payload.Len()), computeChecksum(payload.Bytes()))
	message.Write(header.toBytes())

	buf = message.Bytes()
	return
}
