package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	} `mapstructure:"kafka"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	Cache struct {
		TTLSeconds int `mapstructure:"ttl_seconds"`
	} `mapstructure:"cache"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфига")
	}

	var cnfg Config
	if err := viper.Unmarshal(&cnfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига")
	}

	return &cnfg, nil
}
