package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/robcanini/btcd-node-handshake/internal/config"

	"github.com/rs/zerolog"
)

var GitVersion string

func main() {
	_, ctxCancel := context.WithCancel(context.Background())
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

	// bitcoin node connection
	conn, err := net.Dial("tcp", cfg.BtcNode)
	if err != nil {
		fmt.Println("error while connecting to the bitcoin node:", err)
		return
	}
	defer conn.Close()

	// handshake here

	fmt.Println("handshake accomplished")
}
