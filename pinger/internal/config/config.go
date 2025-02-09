package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrFoundConfigPath error = errors.New("can't found config path in env")
	ErrReadConfigFile  error = errors.New("can't read config file")
)

type Config struct {
	Env    string `yaml:"env"`
	Server Server `yaml:"server"`
	Api    Api    `yaml:"api"`
}

type Server struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Api struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	GetEndpoint  string `yaml:"get_endpoint"`
	PostEndpoint string `yaml:"post_endpoint"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrReadConfigFile
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func MustLoad() *Config {
	path := os.Getenv("PING_CONFIG_PATH")
	if path == "" {
		panic(ErrFoundConfigPath)
	}
	config, err := Load(path)
	if err != nil {
		panic(err)
	}
	return config
}
