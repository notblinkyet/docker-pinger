package pinger

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type Pinger struct {
	Ips map[string]struct{}
}

func (pinger *Pinger) PostIp(update []models.Container) {
	UpdateMap := make(map[string]struct{}, len(update))
	for _, container := range update {
		UpdateMap[container.Ip] = struct{}{}
	}
	pinger.Ips = UpdateMap
}

func (pinger *Pinger) PingAll() []models.Ping {
	chanel := make(chan models.Ping)
	pings := make([]models.Ping, 0, len(pinger.Ips))

	for ip := range pinger.Ips {
		go PingOne(chanel, ip)
		ping := <-chanel
		pings = append(pings, ping)
	}
	return pings
}

func PingOne(out chan<- models.Ping, ip string) {
	data, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
	if err != nil {
		out <- models.Ping{
			PingAt:  time.Now(),
			Success: false,
		}
		return
	}
	output := string(data)
	if strings.Contains(output, "0% packet loss") {
		re := regexp.MustCompile(`time=([0-9.]+) ms`)
		matches := re.FindStringSubmatch(output)
		if len(matches) > 1 {
			latency, _ := strconv.ParseFloat(matches[1], 64)
			out <- models.Ping{
				PingAt:  time.Now(),
				Success: true,
				Latency: latency * (time.Millisecond),
			}
		} else {
			out <- models.Ping{
				PingAt:  time.Now(),
				Success: true,
			}
		}
	} else {
		out <- models.Ping{
			PingAt:  time.Now(),
			Success: false,
		}
	}
}
