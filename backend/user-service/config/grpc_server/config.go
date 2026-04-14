package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env               string        `mapstructure:"env"`
	GRPCPort          int           `mapstructure:"grpc_port"`
	InternalGRPCPort  int           `mapstructure:"internal_grpc_port"`
	HTTPPort          int           `mapstructure:"http_port"`
	DBHost            string        `mapstructure:"db_host"`
	DBPort            string        `mapstructure:"db_port"`
	DBUser            string        `mapstructure:"db_user"`
	DBPassword        string        `mapstructure:"db_password"`
	DBName            string        `mapstructure:"db_name"`
	MaxConns          int32         `mapstructure:"max_conns"`
	MinConns          int32         `mapstructure:"min_conns"`
	MaxConnIdle       time.Duration `mapstructure:"max_conn_idle"`
	ConnTimeout       time.Duration `mapstructure:"conn_timeout"`
	RedisHost         string        `mapstructure:"redis_host"`
	RedisPort         string        `mapstructure:"redis_port"`
	JWTSecret         string        `mapstructure:"jwt_secret"`
	AccessExpiration  time.Duration `mapstructure:"access_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
	TokenIssuer       string        `mapstructure:"token_issuer"`
}

// LoadConfig loads the configuration from YAML file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/home/andrey/projects/music/backend/user-service/config/grpc_server")

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
