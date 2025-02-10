package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/internal/config"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage/container"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage/ping"
)

type Storage struct {
	Ping      ping.IPingStorage
	Container container.IContainerStorage
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		Ping:      ping.NewPingStorage(pool),
		Container: container.NewContainerStorage(pool),
	}
}

func Open(cfg *config.Config) (*Storage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Storage.Username,
		os.Getenv("POSTGRES_PASS"), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return NewStorage(pool), nil
}

func (storage *Storage) Close() {
	storage.Ping.Close()
}
