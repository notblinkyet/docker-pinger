package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/notblinkyet/docker-pinger/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	Config1 *config.Config = config.New("local", "localhost", "localhost", "docker-pinger", "postgres",
		"/home/hobonail/go_projects/docker-pinger/backend/migrations", 5432, 9999, time.Second)
)

func TestLoad(t *testing.T) {
	var testCases = []*config.Config{Config1}

	for _, testConfig := range testCases {
		ConfigLoading(t, testConfig)
	}
}

func ConfigLoading(t *testing.T, testConfig *config.Config) {
	temp, err := os.CreateTemp("", "test_config.yaml")
	assert.NoError(t, err, "Can't create test config file")
	defer os.Remove(temp.Name())

	out, err := yaml.Marshal(testConfig)
	assert.NoError(t, err, "Can't marshal data")
	_, err = temp.Write(out)
	assert.NoError(t, err, "Can't write data in file")

	config, err := config.Load(temp.Name())
	assert.NoError(t, err)
	assert.Equal(t, testConfig.Env, config.Env)
	assert.Equal(t, testConfig.Storage, config.Storage)
	assert.Equal(t, testConfig.Server, config.Server)
}
