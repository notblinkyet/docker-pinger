package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
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
			tx, err := storage.pool.Begin(ctx)
			if err != nil {
				//TODO logs
				continue
			}
			err = Create(tx, &ping, ctx)
			if err != nil {
				//TODO logs
				continue
			}
			err = tx.Commit(ctx)
			if err != nil {
				//TODO logs
				continue
			}
		} else {
			tx, err := storage.pool.Begin(ctx)
			if err != nil {
				//TODO logs
				continue
			}
			defer tx.Rollback(ctx)

			sql := `
				SELECT last_success_at
				FROM last_pings
				WHERE container_id=$1;
			`
			row := tx.QueryRow(ctx, sql, ping.ContainerId)
			err = row.Scan(&ping.LastSuccess)
			if errors.Is(err, pgx.ErrNoRows) {
				ping.LastSuccess = pgtype.Timestamptz{
					Status: pgtype.Null,
				}
			} else if err != nil {
				//TODO logs
				continue
			}
			err = Create(tx, &ping, ctx)
			if err != nil {
				//TODO logs
				continue
			}
			err = tx.Commit(ctx)
			if err != nil {
				//TODO logs
				continue
			}
		}
	}
}

func (storage *PingStorage) Get(from, op string) ([]models.Ping, error) {
	ctx := context.Background()
	sql := fmt.Sprintf(`
		SELECT c.id, p.id, p.latency, p.last_success_at, p.ping_at
		FROM %s p JOIN containers c ON p.container_id=c.id
	`, from)
	rows, err := storage.pool.Query(ctx, sql, from)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	pings := make([]models.Ping, 0, 10)
	for rows.Next() {
		var ping models.Ping
		var latency int
		err = rows.Scan(
			&ping.ContainerId,
			&ping.Id,
			&latency,
			&ping.LastSuccess,
			&ping.PingAt,
		)
		if err != nil {
			//TODO add logs
			continue
		}
		ping.Latency = time.Duration(latency)
		pings = append(pings, ping)
	}
	return pings, nil
}

func (storage *PingStorage) GetLast() ([]models.Ping, error) {
	const op = "storage/ping/getlast"
	return storage.Get("last_pings", op)
}

func (storage *PingStorage) GetAll() ([]models.Ping, error) {
	const op = "storage/ping/getall"
	return storage.Get("pings", op)
}

func Create(tx pgx.Tx, ping *models.Ping, ctx context.Context) error {
	sql1 := `
			INSERT INTO pings(container_id, latency, last_success_at, ping_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id;
			`
	sql2 := `
			UPDATE last_pings
			SET id=$1, latency=$2, last_success_at=$3, ping_at=$4
			WHERE container_id=$5;
		`
	sql3 := `
			INSERT INTO last_pings(id, container_id, latency, last_success_at, ping_at)
			VALUES ($1, $2, $3, $4, $5)
			`

	row := tx.QueryRow(ctx, sql1, ping.ContainerId, ping.Latency.Nanoseconds(), ping.LastSuccess, ping.PingAt)
	err := row.Scan(&ping.Id)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, sql2, ping.Id, ping.Latency.Nanoseconds(), ping.LastSuccess, ping.PingAt, ping.ContainerId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 0 {
		return nil
	}
	_, err = tx.Exec(ctx, sql3, ping.Id, ping.Latency.Nanoseconds(), ping.LastSuccess, ping.PingAt, ping.ContainerId)
	if err != nil {
		return err
	}
	return nil
}
