package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env               string        `yaml:"env"`
	GRPCPort          string        `yaml:"grpc_port"`
	HTTPPort          string        `yaml:"http_port"`
	DBHost            string        `yaml:"db_host"`
	DBPort            string        `yaml:"db_port"`
	DBUser            string        `yaml:"db_user"`
	DBPassword        string        `yaml:"db_password"`
	DBName            string        `yaml:"db_name"`
	RedisHost         string        `yaml:"redis_host"`
	RedisPort         string        `yaml:"redis_port"`
	JWTSecret         string        `yaml:"jwt_secret"`
	AccessExpiration  time.Duration `yaml:"access_expiration"`
	RefreshExpiration time.Duration `yaml:"refresh_expiration"`
	TokenIssuer       string        `yaml:"token_issuer"`
}

// LoadConfig loads the configuration from YAML file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

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
