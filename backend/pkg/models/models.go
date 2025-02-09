package models

import "time"

type Ping struct {
	Id          int
	ContainerId int
	Ip          string
	Latency     time.Duration
	LastSuccess time.Time
	PingAt      time.Time
	Success     bool
}

type Container struct {
	Id int
	Ip string
}
