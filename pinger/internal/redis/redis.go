package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notblinkyet/docker-pinger/pinger/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client redis.Client
	exp    time.Duration
}

func NewRedisClient(cfg *config.Redis) (*RedisClient, error) {
	password := os.Getenv("REDIS_PASS")
	addres := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addres,
		Password: password,
		DB:       cfg.Db,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{
		client: *client,
		exp:    cfg.Exp,
	}, nil
}

func (redis *RedisClient) Set(ip string, SuccessTime time.Time) error {
	ctx := context.Background()
	err := redis.client.Set(ctx, ip, SuccessTime.UTC().String(), redis.exp).Err()
	return err
}

func (redis *RedisClient) Get(ip string) (time.Time, error) {
	ctx := context.Background()
	res, err := redis.client.Get(ctx, ip).Result()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.UTC.String(), res)

}
