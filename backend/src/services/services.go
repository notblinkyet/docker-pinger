package services

import "github.com/notblinkyet/docker-pinger/backend/src/storage"

type Services struct {
	Ping      IPingService
	Container IContainerService
}

func NewServices(storage *storage.Storage) *Services {
	return &Services{
		Ping:      NewPingService(storage.Ping),
		Container: NewContainerService(storage.Container),
	}
}
