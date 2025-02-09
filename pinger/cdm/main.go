package main

import (
	"log"
	"os"

	"github.com/notblinkyet/docker-pinger/pinger/internal/api"
	"github.com/notblinkyet/docker-pinger/pinger/internal/config"
)

func main() {
	cfg := config.MustLoad()
	api := api.NewApi(&cfg.Api)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_, _ = api, logger
}
