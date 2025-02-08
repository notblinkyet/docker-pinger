package services

import (
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/backend/src/storage"
)

type IContainerService interface {
	Create(ip string) error
	Delete(ip string) error
	GetAll() ([]models.Container, error)
}

type ContainerService struct {
	Storage storage.IContainerStorage
}

func NewContainerService(storage storage.IContainerStorage) *ContainerService {
	return &ContainerService{
		Storage: storage,
	}
}

func (service *ContainerService) Create(ip string) error {
	return service.Storage.Create(ip)
}

func (service *ContainerService) Delete(ip string) error {
	return service.Storage.Delete(ip)
}

func (service *ContainerService) GetAll() ([]models.Container, error) {
	return service.Storage.GetAll()
}
