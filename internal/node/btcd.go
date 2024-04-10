package node

import (
	"context"

	"github.com/robcanini/btcd-node-handshake/internal/config"

	"github.com/rs/zerolog"
)

type Btcd struct {
	log    zerolog.Logger
	ctx    context.Context
	config config.Config
}

func NewBtcdTcpClient(log zerolog.Logger, ctx context.Context, config config.Config) *Btcd {
	return &Btcd{
		log:    log,
		ctx:    ctx,
		config: config,
	}
}

func (b *Btcd) Connect() (err error) {
	// todo
	return
}
