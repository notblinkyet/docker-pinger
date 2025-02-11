package ping

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type IPingStorage interface {
	Create([]models.Ping) []error
	GetLast() ([]models.Ping, error)
	GetAll() ([]models.Ping, error)
	Close()
}

type PingStorage struct {
	pool *pgxpool.Pool
}

func NewPingStorage(pool *pgxpool.Pool) *PingStorage {
	return &PingStorage{
		pool: pool,
	}
}

func (storage *PingStorage) Close() {
	storage.pool.Close()
}
