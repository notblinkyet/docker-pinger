package services

import (
	"sync"

	"github.com/notblinkyet/docker-pinger/backend/internal/api/pinger"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage/container"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type IContainerService interface {
	Create(ip string) error
	Delete(ip string) error
	GetAll() ([]models.Container, error)
}

type ContainerService struct {
	Storage   container.IContainerStorage
	PingerApi *pinger.PingerApi
}

func NewContainerService(storage container.IContainerStorage, api *pinger.PingerApi) *ContainerService {
	return &ContainerService{
		Storage:   storage,
		PingerApi: api,
	}
}

func (service *ContainerService) Create(ip string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := service.PingerApi.Post(ip)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := service.Storage.Create(ip)
		if err != nil {
			errChan <- err
		}
	}()

	defer close(errChan)
	wg.Wait()
	for err := range errChan {
		return err
	}
	return nil

}

func (service *ContainerService) Delete(ip string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := service.PingerApi.Delete(ip)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := service.Storage.Delete(ip)
		if err != nil {
			errChan <- err
		}
	}()

	defer close(errChan)
	wg.Wait()
	for err := range errChan {
		return err
	}
	return nil
}

func (service *ContainerService) GetAll() ([]models.Container, error) {
	return service.Storage.GetAll()
}
