package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storagePath"`
	RedisPath   string `yaml:"redisPath"`
	Server      Server `yaml:"server"`
}

type Server struct {
	Address string `yaml:"address"`
}

// LoadConfig loads the configuration from YAML file
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(`D:\zed\golang\message\config\config.yaml`)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
		return nil, err
	}

	return &config, nil
}
