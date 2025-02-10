package handlers

import "github.com/notblinkyet/docker-pinger/backend/internal/services"

type Handlers struct {
	Container IContainerHandler
	Ping      IPingHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Container: NewContainerHandler(services.Container),
		Ping:      NewPingHandler(services.Ping),
	}
}
