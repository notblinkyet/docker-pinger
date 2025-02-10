package services

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyTracked error = errors.New("ip is already being tracked")
	ErrNotTracked     error = errors.New("ip is not being tracked")
)

type IPingService interface {
	Create(ip string) (err error)
	Delete(ip string) (err error)
}

type PingService struct {
	Ips *sync.Map
}

func NewPingService(ips *sync.Map) *PingService {
	return &PingService{
		Ips: ips,
	}
}

func (service *PingService) Create(ip string) (err error) {
	if _, ok := service.Ips.Load(ip); ok {
		return ErrAlreadyTracked
	}
	service.Ips.Store(ip, struct{}{})
	return nil
}

func (service *PingService) Delete(ip string) (err error) {
	if _, ok := service.Ips.Load(ip); !ok {
		return ErrNotTracked
	}
	service.Ips.Delete(ip)
	return nil
}
