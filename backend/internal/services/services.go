package services

import (
	"github.com/notblinkyet/docker-pinger/backend/internal/api/pinger"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage"
)

type Services struct {
	Ping      IPingService
	Container IContainerService
}

func NewServices(storage *storage.Storage, api *pinger.PingerApi) *Services {
	return &Services{
		Ping:      NewPingService(storage.Ping),
		Container: NewContainerService(storage.Container, api),
	}
}
