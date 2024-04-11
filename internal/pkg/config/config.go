package config

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"

	"github.com/vuquang23/poseidon/internal/pkg/repository"
	"github.com/vuquang23/poseidon/internal/pkg/server"
	"github.com/vuquang23/poseidon/pkg/eth"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
	"github.com/vuquang23/poseidon/pkg/redis"
)

type Config struct {
	Common     CommonConfig
	Eth        eth.Config
	Http       server.Config
	Log        logger.Config
	Postgres   postgres.Config
	Redis      redis.Config
	Repository repository.Config
}

func New() Config {
	return Config{}
}

func (c *Config) Load(cfgFile string) error {
	// Default config values
	defaults.SetDefaults(c)

	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Read config file failed (%s)\n", err)

		configBuffer, err := json.Marshal(c)
		if err != nil {
			return err
		}

		err = viper.ReadConfig(bytes.NewBuffer(configBuffer))
		if err != nil {
			return err
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()
	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
