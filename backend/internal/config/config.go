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

type Storage struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Database string `yaml:"dbname"`
}

type Server struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

type PingerApi struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	PostEndpoint   string `yaml:"post_endpoint"`
	DeleteEndpoint string `yaml:"delete_endpoint"`
}

type Config struct {
	Env           string    `yaml:"env"`
	MigrationPath string    `yaml:"migration_path"`
	Storage       Storage   `yaml:"storage"`
	Server        Server    `yaml:"server"`
	PingerApi     PingerApi `yaml:"pinger_api"`
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
	path := os.Getenv("BACKEND_CONFIG_PATH")
	if path == "" {
		panic(ErrFoundConfigPath)
	}
	config, err := Load(path)
	if err != nil {
		panic(err)
	}
	return config
}

func New(env, storageHost, ServerHost, StorageDB, StorageUsername,
	MigrationPath, PingerApiHost, PostEndpoint, DeleteEndpoint string,
	StoragePort, ServerPort, PingerApiPort int, TimeOut time.Duration) *Config {
	return &Config{
		Env:           env,
		MigrationPath: MigrationPath,
		Storage: Storage{
			Host:     storageHost,
			Port:     ServerPort,
			Username: StorageUsername,
			Database: StorageDB,
		},
		Server: Server{
			Host:    ServerHost,
			Port:    ServerPort,
			TimeOut: TimeOut,
		},
		PingerApi: PingerApi{
			Host:           PingerApiHost,
			Port:           PingerApiPort,
			PostEndpoint:   PostEndpoint,
			DeleteEndpoint: DeleteEndpoint,
		},
	}
}
