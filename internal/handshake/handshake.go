package handshake

import (
	"github.com/robcanini/btcd-node-handshake/internal/node"

	"github.com/rs/zerolog"
)

type HEvent string

const (
	Started HEvent = "H_STARTED"
	Error   HEvent = "H_ERROR"
	VerSent HEvent = "H_VERSION_SENT"
	VerAck  HEvent = "H_VER_ACK"
	Done    HEvent = "H_DONE"
)

type Handshake struct {
	log zerolog.Logger

	node     node.Node
	progress chan HEvent
}

func Run(log zerolog.Logger, node node.Node) chan HEvent {
	progrCh := make(chan HEvent)
	h := Handshake{
		log:      log,
		node:     node,
		progress: progrCh,
	}
	go h.initHandshake()
	return h.progress
}

func (h *Handshake) sendHEvent(event HEvent) {
	h.progress <- event
}

func (h *Handshake) initHandshake() {
	defer close(h.progress)
	n := h.node
	if !n.IsConnected() {
		h.log.Debug().Msg("node not connected, quitting handshake")
		h.sendHEvent(Error)
		return
	}
	h.sendHEvent(Started)
	h.log.Debug().Msg("started handshake with the target node")

	err := h.node.SendVersion()
	if err != nil {
		h.log.Error().Err(err).Msg("error sending version to the target node")
		h.sendHEvent(Error)
		return
	}
	h.sendHEvent(VerSent)
	h.log.Debug().Msg("version sent to the target node")
}
