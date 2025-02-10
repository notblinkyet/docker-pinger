package models

import (
	"time"
)

type Ping struct {
	Id               int           `json:"id"`
	ContainerId      int           `json:"container_id"`
	Ip               string        `json:"ip"`
	Latency          time.Duration `json:"latency"`
	LastSuccess      time.Time     `json:"last_success"`
	PingAt           time.Time     `json:"ping_at"`
	Success          bool          `json:"success"`
	WasSuccessBefore bool          `json:"was_success_before"`
}

type Container struct {
	Id        int    `json:"id"`
	Ip        string `json:"ip"`
	IsTracked bool   `json:"is_tracked"`
}

func NewContainer(id int, ip string, isTracked bool) *Container {
	return &Container{
		Id:        id,
		Ip:        ip,
		IsTracked: isTracked,
	}
}
