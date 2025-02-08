package config

import "time"

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
