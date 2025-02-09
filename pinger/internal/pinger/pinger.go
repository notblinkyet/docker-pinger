package pinger

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	var m sync.Mutex
	var w sync.WaitGroup
	pings := make([]models.Ping, 0, len(pinger.Ips))
	for ip := range pinger.Ips {
		w.Add(1)
		go func(ip string) {
			defer w.Done()
			ping := PingOne(ip)
			m.Lock()
			pings = append(pings, *ping)
			m.Unlock()
		}(ip)
	}
	w.Wait()
	return pings
}

func PingOne(ip string) *models.Ping {
	data, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
	if err != nil {
		return &models.Ping{
			PingAt:  time.Now(),
			Success: false,
		}
	}
	output := string(data)
	if strings.Contains(output, "0% packet loss") {
		re := regexp.MustCompile(`time=([0-9.]+) ms`)
		matches := re.FindStringSubmatch(output)
		if len(matches) > 1 {
			latency, _ := strconv.ParseFloat(matches[1], 64)
			nanoseconds := int64(latency * 1000000)
			return &models.Ping{
				PingAt:  time.Now(),
				Success: true,
				Latency: time.Duration(nanoseconds),
			}
		} else {
			return &models.Ping{
				PingAt:  time.Now(),
				Success: true,
			}
		}
	} else {
		return &models.Ping{
			PingAt:  time.Now(),
			Success: false,
		}
	}
}
