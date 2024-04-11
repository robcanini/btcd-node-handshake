package message

import (
	"io"
	"net"
	"time"
)

type NetAddress struct {
	Timestamp time.Time
	Services  uint64
	IP        net.IP
	Port      uint16
}

type MsgVersion struct {
	ProtocolVersion int32
	Services        uint64
	Timestamp       time.Time
	AddrYou         NetAddress
	AddrMe          NetAddress
	Nonce           uint64
	UserAgent       string
	LastBlock       int32
	DisableRelayTx  bool
}

func NewMsgVersion() *MsgVersion {
	return &MsgVersion{
		ProtocolVersion: 0,
		Services:        0,
		Timestamp:       time.Time{},
		AddrYou:         NetAddress{},
		AddrMe:          NetAddress{},
		Nonce:           0,
		UserAgent:       "",
		LastBlock:       0,
		DisableRelayTx:  false,
	}
}

func (msg *MsgVersion) Encode(w io.Writer) error {

	return nil
}
