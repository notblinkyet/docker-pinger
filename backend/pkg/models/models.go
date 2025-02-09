package models

import (
	"time"

	"github.com/jackc/pgtype"
)

type Ping struct {
	Id          int                `json:"id"`
	ContainerId int                `json:"container_id"`
	Ip          string             `json:"ip"`
	Latency     time.Duration      `json:"latency"`
	LastSuccess pgtype.Timestamptz `json:"last_success"`
	PingAt      time.Time          `json:"ping_at"`
	Success     bool               `json:"success"`
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
