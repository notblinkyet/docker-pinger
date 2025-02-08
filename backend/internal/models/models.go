package models

import "time"

type Ping struct {
	Id          int
	ContainerID int
	Latency     int
	LastSuccess time.Time
	PingAt      time.Time
	Success     bool
}

type Container struct {
	Id int
	Ip string
}
