package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string           `yaml:"env" env-default:"local"`
	DbPath           string           `yaml:"database-path" env-required:"true"`
	TelegramToken    string           `yaml:"telegram-token" env-required:"true"`
	ContentProviders []ProviderConfig `yaml:"content-providers"`
}

type ProviderConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout"`
	Name    string        `yaml:"name" env-required:"true"`
}

func LoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "config.yaml", "path to config file")
	flag.Parse()

	return res
}
