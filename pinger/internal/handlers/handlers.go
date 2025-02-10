package handlers

import "github.com/notblinkyet/docker-pinger/pinger/internal/services"

type Handlers struct {
	Ping IPingHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Ping: NewPingHandler(services.Ping),
	}
}
