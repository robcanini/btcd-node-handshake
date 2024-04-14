package message

const (
	VersionCommand string = "version"
)

type MsgHeader struct {
	magic    uint32  // 4 bytes
	command  string  // 12 bytes
	length   uint32  // 4 bytes
	checksum [4]byte // 4 bytes
}
