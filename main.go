package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/robcanini/btcd-node-handshake/internal/config"
	"github.com/robcanini/btcd-node-handshake/internal/handshake"
	"github.com/robcanini/btcd-node-handshake/internal/node"

	"github.com/rs/zerolog"
)

var GitVersion string

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	log := createLogger()

	var err error
	defer func() {
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error")
		}
	}()

	// flags
	var (
		configPath   string
		printVersion bool
	)

	flag.StringVar(&configPath, "config", "config/config.yml", "config file path")
	flag.BoolVar(&printVersion, "version", false, "print version")
	flag.Parse()

	if printVersion {
		fmt.Printf("btcd-node-handshake %s\n", GitVersion)
		return
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return
	}
	log.Debug().
		Str("path", configPath).
		Msg("loaded configuration file")

	level, err := zerolog.ParseLevel(cfg.Loglevel)
	if err != nil {
		return
	}
	zerolog.SetGlobalLevel(level)

	// btcd node connection
	btcd := node.NewBtcdTcpClient(log, ctx, cfg.Node)
	dispose, err := btcd.Connect()
	if err != nil {
		return
	}
	defer dispose()

	// handshake
	hProgr := handshake.Run(log, btcd)
	for event := range hProgr {
		log.Info().Str("code", string(event)).Msg("h_event received")
		if event == handshake.Error {
			err = errors.New("handshake failed")
			return
		}
		if event == handshake.Done {
			log.Info().Msg("handshake successfully completed")
		}
	}
}
