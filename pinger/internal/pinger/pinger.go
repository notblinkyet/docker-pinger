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

	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/pinger/internal/client"
	"github.com/notblinkyet/docker-pinger/pinger/internal/redis"
)

type Pinger struct {
	Ips    *sync.Map
	count  *atomic.Int64
	client *client.Client
	logger *log.Logger
	redis  *redis.RedisClient
}

func NewPinger(newIps []string, client *client.Client, logger *log.Logger, redis *redis.RedisClient) *Pinger {
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
		redis:  redis,
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
			ping := pinger.PingOne(ip)
			m.Lock()
			pings = append(pings, *ping)
			m.Unlock()
		}(ip)
		return true
	})
	w.Wait()
	return pings
}

func (pinger *Pinger) PingOne(ip string) *models.Ping {
	data, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
	pinger.logger.Println(string(data))
	t := time.Now()
	if err != nil {
		LastSuccess, err := pinger.redis.Get(ip)
		if err != nil {
			pinger.logger.Println(err)
			return &models.Ping{
				PingAt:           t,
				Success:          false,
				Ip:               ip,
				WasSuccessBefore: false,
			}
		}
		return &models.Ping{
			PingAt:           t,
			Success:          false,
			Ip:               ip,
			WasSuccessBefore: true,
			LastSuccess:      LastSuccess,
		}
	}
	output := string(data)
	if strings.Contains(output, "0% packet loss") {
		re := regexp.MustCompile(`time=([0-9.]+) ms`)
		matches := re.FindStringSubmatch(output)
		err := pinger.redis.Set(ip, t)
		if err != nil {
			pinger.logger.Println(err)
		}
		if len(matches) > 1 {
			latency, _ := strconv.ParseFloat(matches[1], 64)
			nanoseconds := int64(latency * 1000000)
			return &models.Ping{
				PingAt:           t,
				LastSuccess:      t,
				Success:          true,
				Latency:          time.Duration(nanoseconds),
				Ip:               ip,
				WasSuccessBefore: true,
			}
		} else {
			return &models.Ping{
				PingAt:           t,
				Success:          true,
				Ip:               ip,
				LastSuccess:      t,
				WasSuccessBefore: true,
			}
		}
	}
	LastSuccess, err := pinger.redis.Get(ip)
	if err != nil {
		pinger.logger.Println(err)
		return &models.Ping{
			PingAt:           t,
			Success:          false,
			Ip:               ip,
			WasSuccessBefore: false,
		}
	}
	return &models.Ping{
		PingAt:           t,
		Success:          false,
		Ip:               ip,
		WasSuccessBefore: true,
		LastSuccess:      LastSuccess,
	}
}

func (pinger *Pinger) StartPinging(delta time.Duration) {
	ticker := time.NewTicker(delta)
	defer ticker.Stop()

	for range ticker.C {
		pings := pinger.PingAll()
		pinger.logger.Println(pings)
		err := pinger.client.Post(pings)
		if err != nil {
			pinger.logger.Println(err)
		}
	}
}
