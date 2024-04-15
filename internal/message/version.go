package message

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

// NetAddress implements the addrv2 standard
type NetAddress struct {
	Timestamp time.Time
	Services  uint64
	IP        net.IP
	Port      uint16
}

func (net *NetAddress) ToBytes() (buf []byte) {
	buf = make([]byte, 26)
	binary.LittleEndian.PutUint32(buf[:4], uint32(net.Timestamp.Unix()))
	binary.LittleEndian.PutUint64(buf, uint64(net.Services))
	copy(buf[:16], net.IP.To16())
	binary.BigEndian.PutUint16(buf[:2], net.Port)
	return
}

type MsgVersion struct {
	Network         uint32
	ProtocolVersion uint32
	Services        uint64
	Timestamp       time.Time
	AddrYou         NetAddress
	AddrMe          NetAddress
	Nonce           uint64
	UserAgent       string
	LastBlock       int32
	RelayTx         bool
}

func NewMsgVersion(
	network uint32,
	protocolVersion uint32,
	services uint64,
	addrCur NetAddress,
	addrTar NetAddress,
	userAger string,
	lastBlock int32,
	relayTx bool,
) *MsgVersion {
	return &MsgVersion{
		Network:         network,
		ProtocolVersion: protocolVersion,
		Services:        services,
		AddrYou:         addrCur,
		AddrMe:          addrTar,
		UserAgent:       userAger,
		LastBlock:       lastBlock,
		RelayTx:         relayTx,
	}
}

func (msg *MsgVersion) ToBytes() (buf []byte, err error) {
	Timestamp := time.Unix(time.Now().Unix(), 0).Unix()
	LastBlock := uint32(212672)

	payload := new(bytes.Buffer)
	_ = binary.Write(payload, binary.LittleEndian, msg.ProtocolVersion)
	_ = binary.Write(payload, binary.LittleEndian, msg.Services)
	_ = binary.Write(payload, binary.LittleEndian, uint64(Timestamp))
	_ = binary.Write(payload, binary.LittleEndian, msg.AddrMe.ToBytes())
	_ = binary.Write(payload, binary.LittleEndian, msg.AddrYou.ToBytes())
	_ = binary.Write(payload, binary.LittleEndian, randUint64())
	_ = binary.Write(payload, binary.LittleEndian, uint8(len(msg.UserAgent)))
	payload.Write([]byte(msg.UserAgent))
	_ = binary.Write(payload, binary.LittleEndian, LastBlock)
	_ = binary.Write(payload, binary.LittleEndian, !msg.RelayTx) // protocol wants DisableRelayTx

	message := new(bytes.Buffer)

	// header
	header := newMsgHeader(msg.Network, CmdVersion, uint32(payload.Len()), computeChecksum(payload.Bytes()))
	message.Write(header.toBytes())

	// payload
	message.Write(payload.Bytes())

	buf = message.Bytes()
	return
}
