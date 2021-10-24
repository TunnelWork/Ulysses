package conf

import (
	"errors"

	"github.com/TunnelWork/Ulysses/src/internal/db"
	"github.com/TunnelWork/Ulysses/src/internal/logger"
	"gopkg.in/yaml.v2"
)

var (
	ErrBadConfig = errors.New("internal/conf: bad config content")
)

type Config struct {
	Sys SystemConfig        `yaml:"sys"`
	Log logger.LoggerConfig `yaml:"log"`
	DB  db.DatabaseConfig   `yaml:"db"`
}

func LoadUlyssesConfig(content []byte) (Config, error) {
	newConfig := Config{
		Sys: defaultSystemConfig(),
		Log: defaultLoggerConfig(),
		DB:  defaultDatabaseConfig(),
	}
	err := yaml.Unmarshal(content, &newConfig)

	return newConfig, err
}
