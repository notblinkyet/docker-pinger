package services

import (
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/backend/src/storage"
)

type IPingService interface {
	GetAll() ([]models.Ping, error)
	GetLast() ([]models.Ping, error)
	Create([]models.Ping)
}

type PingService struct {
	Storage storage.IPingStorage
}

func NewPingService(storage storage.IPingStorage) *PingService {
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
