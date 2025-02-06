package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	TestApiBybit        string `yaml:"testApiBybit"`
	TestKey             string `yaml:"testKey"`
	TestSecret          string `yaml:"testSecret"`
	TestWebSocketClient string `yaml:"testWebSocketClientBybit"`
	PostgresUser        string `yaml:"postgresUser"`
	PostgresPassword    string `yaml:"postgresPassword"`
	PostgresDB          string `yaml:"postgresDB"`
	SslMode             string `yaml:"sslMode"`
}

func MustLoad() *Config {
	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %s", configPath)
	}
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config file %s: %s", configPath, err)
	}
	return &cfg
}
