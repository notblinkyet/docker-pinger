package handlers

import "github.com/notblinkyet/docker-pinger/backend/src/services"

type Handlers struct {
	Containers IContainerHandler
	Pings      IPingHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Containers: NewContainerHandler(services.Container),
		Pings:      NewPingHandler(services.Ping),
	}
}
