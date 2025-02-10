package services

import "sync"

type Services struct {
	Ping IPingService
}

func NewServices(ips *sync.Map) *Services {
	return &Services{
		Ping: NewPingService(ips),
	}
}
