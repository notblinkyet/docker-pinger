package models

import "time"

type Ping struct {
	Id          int
	ContainerID int
	Latency     time.Duration
	LastSuccess time.Time
	PingAt      time.Time
	Success     bool
}

type Container struct {
	Id int
	Ip string
}
