package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Storage  string        `yaml:"storage" required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" required:"true"`
	Server   ServerConfig
}

type ServerConfig struct {
	Port    string `yaml:"port" env-required:"true"`
	Timeout string `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		path = "./config/local.yaml"
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return path
}
