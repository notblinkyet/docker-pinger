package services

import (
	"log"
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
	logger    *log.Logger
}

func NewContainerService(storage container.IContainerStorage, api *pinger.PingerApi, logger *log.Logger) *ContainerService {
	return &ContainerService{
		Storage:   storage,
		PingerApi: api,
		logger:    logger,
	}
}

func (service *ContainerService) Create(ip string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := service.PingerApi.Post(ip); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := service.Storage.Create(ip); err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	errors := make([]error, 0, 2)
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		service.logger.Println(errors)
		return errors[0]
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

	go func() {
		wg.Wait()
		close(errChan)
	}()

	errors := make([]error, 0, 2)
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		service.logger.Println(errors)
		return errors[0]
	}

	return nil
}

func (service *ContainerService) GetAll() ([]models.Container, error) {
	containers, err := service.Storage.GetAll()
	if err != nil {
		service.logger.Println(err)
	}
	service.logger.Println(containers)
	return containers, err
}
