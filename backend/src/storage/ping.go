package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/internal/models"
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
	for _, ping := range pings {
		var (
			last interface{}
			id   int
		)
		sql1 := `
		SELECT last_success_at 
		FROM last_ping 
		WHERE container_id=$1;`
		sql2 := `
		INSERT INTO ping(container_id, latency, last_success_at, ping_at)
		VALUES($1, $2, $3, $4)
		RETURNING id;
		`
		sql3 := `
		UPDATE last_ping
		SET id=$1, latency=$2,last_success_at=$3,ping_at=$4
		WHERE container_id=$5;
		`
		tx, err := storage.pool.Begin(context.Background())
		if err != nil {
			//TODO logs
		}
		row := tx.QueryRow(context.Background(), sql1, ping.ContainerID)
		err = row.Scan(&last)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				//TODO logs
				tx.Rollback(context.Background())
				continue
			}
			//TODO logs
			last = nil
		}
		if ping.Success {
			last = ping.PingAt
		}
		row = tx.QueryRow(context.Background(), sql2, ping.ContainerID, ping.Latency, last, ping.PingAt)
		err = row.Scan(&id)
		if err != nil {
			//TODO logs
			tx.Rollback(context.Background())
			continue
		}
		_, err = tx.Exec(context.Background(), sql3, id, ping.Latency, last, ping.PingAt, ping.ContainerID)
		if err != nil {
			//TODO logs
			tx.Rollback(context.Background())
			continue
		}
		err = tx.Commit(context.Background())
		if err != nil {
			//TODO logs
		}
	}
}

func (storage *PingStorage) GetLast() ([]models.Ping, error) {
	const op = "storage/ping/getlast"
	pings := make([]models.Ping, 0, 10)
	sql := `SELECT id, container_id, latency, last_success_at, ping_at FROM last_ping;`
	rows, err := storage.pool.Query(context.Background(), sql)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var ping models.Ping
		err = rows.Scan(
			&ping.Id,
			&ping.ContainerID,
			&ping.Latency,
			&ping.LastSuccess,
			&ping.PingAt)
		if err != nil {
			//TODO logs
			continue
		}
		pings = append(pings, ping)
	}

	return pings, nil

}

func (storage *PingStorage) GetAll() ([]models.Ping, error) {
	const op = "storage/ping/getall"
	pings := make([]models.Ping, 0, 10)
	sql := `SELECT id, container_id, latency, last_success_at, ping_at FROM ping;`
	rows, err := storage.pool.Query(context.Background(), sql)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var ping models.Ping
		err = rows.Scan(
			&ping.Id,
			&ping.ContainerID,
			&ping.Latency,
			&ping.LastSuccess,
			&ping.PingAt)
		if err != nil {
			//TODO logs
			continue
		}
		pings = append(pings, ping)
	}

	return pings, nil
}
