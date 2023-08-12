package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Dsn         string `mapstructure:"dsn"`
	HttpAddress string `mapstructure:"http_address"`
}

func ParseConfig(filename string) (*Config, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannout be empty")
	}
	var config *Config
	v := viper.New()
	v.SetConfigFile(filename)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
