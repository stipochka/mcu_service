package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mcu_service/internal/models"
)

type Config struct {
	Env           string          `yaml:"env" env-default:"local"`
	PortName      string          `yaml:"port_name" env-required:"true"`
	IntervalToAsk time.Duration   `yaml:"interval_to_ask" env-default:"60s"`
	BrokerAddress string          `yaml:"broker_address" env-required:"true"`
	Devices       []models.Record `yaml:"devices" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("file %s does not exists", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config %s", err.Error())
	}

	return &cfg
}
