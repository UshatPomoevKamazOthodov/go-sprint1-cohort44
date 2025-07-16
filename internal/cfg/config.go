package cfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

var cfg *Config

func init() {
	cfg = &Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Println("no environment variables found")
		flag.Parse()

		flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "адрес сервера")
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080/", "базовый адрес")
		log.Println(cfg.ServerAddr)
	}
}

func GetConfigData() *Config {
	return cfg
}
