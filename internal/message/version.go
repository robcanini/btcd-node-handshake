package message

import (
	"net"
	"time"
)

type NetAddress struct {
	// Last time the address was seen.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.  This field is
	// not present in the bitcoin version message (MsgVersion) nor was it
	// added until protocol version >= NetAddressTimeVersion.
	Timestamp time.Time

	// Bitfield which identifies the services supported by the address.
	Services uint64

	// IP address of the peer.
	IP net.IP

	// Port the peer is using.  This is encoded in big endian on the wire
	// which differs from most everything else.
	Port uint16
}

type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32

	// Bitfield which identifies the enabled services.
	Services uint64

	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time

	// Address of the remote peer.
	AddrYou NetAddress

	// Address of the local peer.
	AddrMe NetAddress

	// Unique value associated with message that is used to detect self
	// connections.
	Nonce uint64

	// The user agent that generated message.  This is a encoded as a varString
	// on the wire.  This has a max length of MaxUserAgentLen.
	UserAgent string

	// Last block seen by the generator of the version message.
	LastBlock int32

	// Don't announce transactions to peer.
	DisableRelayTx bool
}
