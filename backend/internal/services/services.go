package services

import (
	"log"

	"github.com/notblinkyet/docker-pinger/backend/internal/api/pinger"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage"
)

type Services struct {
	Ping      IPingService
	Container IContainerService
}

func NewServices(storage *storage.Storage, api *pinger.PingerApi, logger *log.Logger) *Services {
	return &Services{
		Ping:      NewPingService(storage.Ping, logger),
		Container: NewContainerService(storage.Container, api, logger),
	}
}
