package config

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	// Define test configuration data
	configData := `
loglevel: debug
node:
  host: example.com
  port: 8080
  btcd:
    agent: test-agent
    p_version: 70015
    relay_tx: true
    services: 0
    network: 0
handshake:
  timeout: 10s
`
	// Set up viper with test configuration data
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(configData))
	if err != nil {
		t.Fatalf("failed to set up viper: %v", err)
	}

	// Load configuration
	c, err := Load("testdata/test.yml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check if configuration fields are correctly loaded
	expectedLogLevel := "debug"
	if c.Loglevel != expectedLogLevel {
		t.Errorf("Expected loglevel %s, got %s", expectedLogLevel, c.Loglevel)
	}

	expectedNodeHost := "example.com"
	if c.Node.Host != expectedNodeHost {
		t.Errorf("Expected node host %s, got %s", expectedNodeHost, c.Node.Host)
	}

	expectedNodePort := uint16(8080)
	if c.Node.Port != expectedNodePort {
		t.Errorf("Expected node port %d, got %d", expectedNodePort, c.Node.Port)
	}

	expectedAgent := "test-agent"
	if c.Node.Btcd.Agent != expectedAgent {
		t.Errorf("Expected agent %s, got %s", expectedAgent, c.Node.Btcd.Agent)
	}

	expectedPVersion := uint32(70015)
	if c.Node.Btcd.ProtocolVersion != expectedPVersion {
		t.Errorf("Expected protocol version %d, got %d", expectedPVersion, c.Node.Btcd.ProtocolVersion)
	}

	expectedRelayTx := true
	if c.Node.Btcd.RelayTx != expectedRelayTx {
		t.Errorf("Expected relay_tx %t, got %t", expectedRelayTx, c.Node.Btcd.RelayTx)
	}

	expectedServices := uint32(0)
	if c.Node.Btcd.Services != expectedServices {
		t.Errorf("Expected services %d, got %d", expectedServices, c.Node.Btcd.Services)
	}

	expectedNetwork := uint32(0)
	if c.Node.Btcd.Network != expectedNetwork {
		t.Errorf("Expected network %d, got %d", expectedNetwork, c.Node.Btcd.Network)
	}

	expectedTimeout := 10 * time.Second
	if c.Handshake.Timeout != expectedTimeout {
		t.Errorf("Expected handshake timeout %v, got %v", expectedTimeout, c.Handshake.Timeout)
	}
}
