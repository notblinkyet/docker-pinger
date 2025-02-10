package pinger

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgtype"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/pinger/internal/client"
)

type Pinger struct {
	Ips    *sync.Map
	count  *atomic.Int64
	client *client.Client
	logger *log.Logger
}

func NewPinger(newIps []string, client *client.Client, logger *log.Logger) *Pinger {
	Ips := &sync.Map{}
	count := &atomic.Int64{}
	count.Store(int64(len(newIps)))
	for _, ip := range newIps {
		Ips.Store(ip, struct{}{})
	}
	return &Pinger{
		Ips:    Ips,
		count:  count,
		client: client,
		logger: logger,
	}
}

func (pinger *Pinger) PingAll() []models.Ping {
	var m sync.Mutex
	var w sync.WaitGroup
	pings := make([]models.Ping, 0, pinger.count.Load())
	pinger.Ips.Range(func(key, value any) bool {
		w.Add(1)
		ip := key.(string)
		go func(ip string) {
			defer w.Done()
			ping := PingOne(ip)
			m.Lock()
			pings = append(pings, *ping)
			m.Unlock()
		}(ip)
		return true
	})
	w.Wait()
	return pings
}

func PingOne(ip string) *models.Ping {
	data, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
	if err != nil {
		return &models.Ping{
			PingAt:  time.Now(),
			Success: false,
			Ip:      ip,
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
				LastSuccess: pgtype.Timestamptz{
					Time: time.Now(),
				},
				Success: true,
				Latency: time.Duration(nanoseconds),
				Ip:      ip,
			}
		} else {
			return &models.Ping{
				PingAt:  time.Now(),
				Success: true,
				Ip:      ip,
				LastSuccess: pgtype.Timestamptz{
					Time: time.Now(),
				},
			}
		}
	} else {
		return &models.Ping{
			PingAt:  time.Now(),
			Success: false,
			Ip:      ip,
		}
	}
}

func (pinger *Pinger) PingOnceAfterDelay(delay time.Duration) {
	time.AfterFunc(delay, func() {
		pings := pinger.PingAll()
		err := pinger.client.Post(pings)
		pinger.logger.Println(err)
	})
}
