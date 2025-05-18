package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress      string
	MaxConnections     int
	LogLevel           string
	MaxRequestsPerConn int

	ClientReadTimeoutSeconds         int
	ClientWriteTimeoutSeconds        int
	ClientMaxIdleConnDurationSeconds int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("ServerAddress", ":8080")
	viper.SetDefault("MaxConnections", 10000)
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("MaxRequestsPerConn", 5000)
	viper.SetDefault("ClientReadTimeoutSeconds", 15)
	viper.SetDefault("ClientWriteTimeoutSeconds", 15)
	viper.SetDefault("ClientMaxIdleConnDurationSeconds", 60)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate the configuration
	if config.ServerAddress == "" {
		return nil, fmt.Errorf("ServerAddress cannot be empty")
	}
	if config.MaxConnections <= 0 {
		return nil, fmt.Errorf("MaxConnections must be greater than 0")
	}
	if config.MaxRequestsPerConn <= 0 {
		return nil, fmt.Errorf("MaxRequestsPerConn must be greater than 0")
	}
	if config.ClientReadTimeoutSeconds <= 0 {
		return nil, fmt.Errorf("ClientReadTimeoutSeconds must be greater than 0")
	}
	if config.ClientWriteTimeoutSeconds <= 0 {
		return nil, fmt.Errorf("ClientWriteTimeoutSeconds must be greater than 0")
	}
	if config.ClientMaxIdleConnDurationSeconds <= 0 {
		return nil, fmt.Errorf("ClientMaxIdleConnDurationSeconds must be greater than 0")
	}
	return &config, nil
}
