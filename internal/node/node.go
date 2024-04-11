package node

type Node interface {
	Connect() (func(), error)
	IsConnected() bool
	SendVersion() error
	VerAck() error
}

type connection interface {
	dispose()
	read([]byte) error
	write([]byte) error
}
