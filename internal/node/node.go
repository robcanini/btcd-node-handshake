package node

type Node interface {
	// Connect to the node returning the conn handler
	Connect() (*Connection, error)
}

type Connection interface {
	Dispose() error
	Read([]byte) error
	Write([]byte) error
}
