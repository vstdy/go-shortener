package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	BaseURL         string `mapstructure:"BASE_URL"`
	FileStoragePath string `mapstructure:"FILE_STORAGE_PATH"`
	SecretKey       string `mapstructure:"SECRET_KEY"`
}

func LoadEnvs() (Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)

	return cfg, err
}
