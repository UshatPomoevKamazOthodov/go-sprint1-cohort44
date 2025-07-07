package cfg

import "flag"

type Config struct {
	ServerAddr string
	BaseURL    string
}

var cfg *Config

func init() {
	cfg = &Config{}
	flag.Parse()

	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "адрес сервера")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080/", "базовый адрес")
}

func GetConfigData() *Config {
	return cfg
}
