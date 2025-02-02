package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	StoragePath string         `yaml:"storage_path" env-required:"true"`
	TokenTTL    TokenTTLConfig `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig     `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type TokenTTLConfig struct {
	Auth    time.Duration `yaml:"auth" env-required:"true"`
	Refresh time.Duration `yaml:"refresh" env-required:"true"`
}

// MustLoad trying to read config in yaml format.
// Priority of loading: flag->env->default.
// If not loaded panic.
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Error while reading config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetching path of config and returned item.
// Priority of loading: flag->env->default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "Path to config file in yaml format")

	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
