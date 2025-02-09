package pinger

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgtype"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type Pinger struct {
	Ips map[string]int
}

func (pinger *Pinger) PostIp(update []models.Container) {
	UpdateMap := make(map[string]int, len(update))
	for _, container := range update {
		UpdateMap[container.Ip] = container.Id
	}
	pinger.Ips = UpdateMap
}

func (pinger *Pinger) PingAll() []models.Ping {
	var m sync.Mutex
	var w sync.WaitGroup
	pings := make([]models.Ping, 0, len(pinger.Ips))
	for ip, id := range pinger.Ips {
		w.Add(1)
		go func(ip string, id int) {
			defer w.Done()
			ping := PingOne(&models.NewContainer(id, ip))
			m.Lock()
			pings = append(pings, *ping)
			m.Unlock()
		}(ip, id)
	}
	w.Wait()
	return pings
}

func PingOne(container *models.Container) *models.Ping {
	data, err := exec.Command("ping", "-c", "1", "-W", "1", container.Ip).Output()
	if err != nil {
		return &models.Ping{
			PingAt:      time.Now(),
			Success:     false,
			ContainerID: container.Id,
			Ip:          ip,
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
				PingAt: time.Now(),
				LastSuccess: pgtype.Timestamp{
					time: time.Now(),
				},
				Success:     true,
				Latency:     time.Duration(nanoseconds),
				ContainerID: container.Id,
				Ip:          ip,
			}
		} else {
			return &models.Ping{
				PingAt:      time.Now(),
				Success:     true,
				ContainerID: container.Id,
				Ip:          ip,
				LastSuccess: pgtype.Timestamp{
					time: time.Now(),
				},
			}
		}
	} else {
		return &models.Ping{
			PingAt:      time.Now(),
			Success:     false,
			ContainerID: container.Id,
			Ip:          ip,
		}
	}
}
