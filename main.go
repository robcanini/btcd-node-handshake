package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/robcanini/btcd-node-handshake/internal/config"
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
	btcd := node.NewBtcdTcpClient(log, ctx, cfg.Node, cfg.Host)
	btcdDispose, err := btcd.Connect()
	if err != nil {
		return
	}
	defer btcdDispose()

	// init handshake sending version msg to btcd node
	ch := make(chan node.HandshakeCode)
	err = btcd.StartHandshake(ch)
	if err != nil {
		return
	}

	// wait for the handshake to be completed or timeout
	hcode := waitForHandshake(ch, cfg, log)
	if hcode == node.HDone {
		log.Info().Msg("handshake completed")
	} else {
		log.Error().Int("code", int(hcode)).Msg("handshake failed")
	}
}

const (
	ExitError = 1
)

func waitForHandshake(ch chan node.HandshakeCode, cfg config.Config, log zerolog.Logger) node.HandshakeCode {
	for {
		select {
		// code
		case code := <-ch:
			return code
		// connection timeout
		case <-time.After(cfg.Handshake.Timeout):
			log.Error().Msg("timeout reached before handshake")
			os.Exit(ExitError)
		}
	}
}
