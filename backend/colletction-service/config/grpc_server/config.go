package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string `mapstructure:"env"`
	GRPCPort   int    `mapstructure:"grpc_port"`
	HTTPPort   int    `mapstructure:"http_port"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	RedisHost  string `mapstructure:"redis_host"`
	RedisPort  int    `mapstructure:"redis_port"`
}

// LoadConfig loads the configuration from YAML file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/home/andrey/projects/music/backend/track-service/config/grpc_server")

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
