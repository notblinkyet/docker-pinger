package services

import (
	"log"

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
	logger  *log.Logger
}

func NewPingService(storage ping.IPingStorage, logger *log.Logger) *PingService {
	return &PingService{
		Storage: storage,
		logger:  logger,
	}
}

func (service *PingService) GetAll() ([]models.Ping, error) {
	pings, err := service.Storage.GetAll()
	if err != nil {
		service.logger.Println(err)
		return pings, err
	}
	service.logger.Println(pings)
	return pings, err
}

func (service *PingService) GetLast() ([]models.Ping, error) {
	pings, err := service.Storage.GetLast()
	if err != nil {
		service.logger.Println(err)
		return pings, err
	}
	service.logger.Println(pings)
	return pings, err
}

func (service *PingService) Create(pings []models.Ping) {
	errs := service.Storage.Create(pings)
	if len(errs) != 0 {
		service.logger.Println(errs)
	}
}
