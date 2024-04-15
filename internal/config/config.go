package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Loglevel  string    `mapstructure:"loglevel"`
		Node      Node      `mapstructure:"node"`
		Handshake Handshake `mapstructure:"handshake"`
	}
	Node struct {
		Host string `mapstructure:"host"`
		Port uint16 `mapstructure:"port"`
		Btcd Btcd   `mapstructure:"btcd"`
	}
	Btcd struct {
		Agent           string `mapstructure:"agent"`
		ProtocolVersion uint32 `mapstructure:"p_version"`
		RelayTx         bool   `mapstructure:"relay_tx"`
		Services        uint32 `mapstructure:"services"`
		Network         uint32 `mapstructure:"network"`
	}
	Handshake struct {
		Timeout time.Duration `mapstructure:"timeout"`
	}
)

func (node *Node) Address() string {
	return fmt.Sprintf("%s:%d", node.Host, node.Port)
}

func Load(file string) (c Config, err error) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BTCD_NODE_HANDSHAKE")
	err = viper.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("cannot read config: %w", err)
		return
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal config: %w", err)
		return
	}
	return
}
