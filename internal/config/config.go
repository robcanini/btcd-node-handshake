package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Loglevel string `mapstructure:"loglevel"`
		BtcdNode string `mapstructure:"btcd_node"`
	}
)

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
