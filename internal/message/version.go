package message

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"net"
	"time"
)

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
	// todo: isolate each attr to serializable comp behaviour
	
	sourceAddr := NetAddress{
		IP:   net.ParseIP("127.0.0.1").To4(),
		Port: 8443,
	}
	targetAddr := NetAddress{
		IP:   net.ParseIP("127.0.0.1").To4(),
		Port: 8443,
	}
	Timestamp := time.Unix(time.Now().Unix(), 0).Unix()
	LastBlock := uint32(212672)

	payload := new(bytes.Buffer)
	_ = binary.Write(payload, binary.LittleEndian, msg.ProtocolVersion)
	_ = binary.Write(payload, binary.LittleEndian, msg.Services)
	_ = binary.Write(payload, binary.LittleEndian, uint64(Timestamp))
	_ = binary.Write(payload, binary.LittleEndian, sourceAddr.ToBytes())
	_ = binary.Write(payload, binary.LittleEndian, targetAddr.ToBytes())
	_ = binary.Write(payload, binary.LittleEndian, randUint64())
	_ = binary.Write(payload, binary.LittleEndian, uint8(len(msg.UserAgent)))
	payload.Write([]byte(msg.UserAgent))
	_ = binary.Write(payload, binary.LittleEndian, LastBlock)
	_ = binary.Write(payload, binary.LittleEndian, !msg.RelayTx) // protocol wants DisableRelayTx

	message := new(bytes.Buffer)

	// header
	_ = binary.Write(message, binary.LittleEndian, msg.Network)
	command := []byte(VersionCommand)
	command = append(command, make([]byte, 12-len(command))...)
	message.Write(command)
	_ = binary.Write(message, binary.LittleEndian, uint32(payload.Len()))
	message.Write(computeChecksum(payload.Bytes())[0:4])

	// payload
	message.Write(payload.Bytes())

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
	rand.Read(buf) // no need to check error
	return binary.LittleEndian.Uint64(buf)
}
