package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DBConfig struct {
	User     string
	Password string
	DbName   string
	Host     string
	Port     string
}

func LoadConfig() (*DBConfig, error) {
	viper.SetConfigFile("/app/.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error to load env %v", err)
	}

	return &DBConfig{
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASS"),
		DbName:   viper.GetString("DB_NAME"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
	}, nil
}
