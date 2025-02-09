package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type IPingStorage interface {
	Create([]models.Ping)
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

func (storage *PingStorage) Create(pings []models.Ping) {
	//const op = "storage/ping/update"
	ctx := context.Background()
	for _, ping := range pings {
		if ping.Success {

		} else {
			tx, err := storage.pool.Begin(ctx)
			if err != nil {
				//TODO logs
				continue
			}
			defer func() {
				if err != nil {
					//TODO logs
					tx.Rollback(ctx)
				}
			}()
			sql := `
			SELECT last_success_at 
			FROM ping
			WHERE container.id=$1
			ORDER BY ping_at DESC
			LIMIT 1;
			`
			tx.QueryRow(ctx, sql, ping.Ip)
		}
	}
}

func (storage *PingStorage) GetLast() ([]models.Ping, error) {
	const op = "storage/ping/getlast"
	return nil, nil
}

func (storage *PingStorage) GetAll() ([]models.Ping, error) {
	const op = "storage/ping/getall"
	return nil, nil
}
