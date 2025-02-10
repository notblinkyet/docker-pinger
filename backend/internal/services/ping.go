package services

import (
	"github.com/notblinkyet/docker-pinger/backend/internal/storage/ping"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type IPingService interface {
	GetAll() ([]models.Ping, error)
	GetLast() ([]models.Ping, error)
	Create([]models.Ping)
}

type PingService struct {
	Storage ping.IPingStorage
}

func NewPingService(storage ping.IPingStorage) *PingService {
	return &PingService{
		Storage: storage,
	}
}

func (service *PingService) GetAll() ([]models.Ping, error) {
	return service.Storage.GetAll()
}

func (service *PingService) GetLast() ([]models.Ping, error) {
	return service.Storage.GetLast()
}

func (service *PingService) Create(pings []models.Ping) {
	service.Storage.Create(pings)
}
