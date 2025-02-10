package ping

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

func (storage *PingStorage) Create(pings []models.Ping) {
	//const op = "storage/ping/update"
	ctx := context.Background()
	for _, ping := range pings {
		tx, err := storage.pool.Begin(ctx)
		if err != nil {
			//TODO logs
			continue
		}
		err = Create(tx, &ping, ctx)
		if err != nil {
			//TODO: logs
			tx.Rollback(ctx)
			continue
		}
		err = tx.Commit(ctx)
		if err != nil {
			//TODO: logs
		}
	}
}

func Create(tx pgx.Tx, ping *models.Ping, ctx context.Context) error {
	var lastSuccessAt pgtype.Timestamp
	if ping.WasSuccessBefore == false {
		lastSuccessAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		lastSuccessAt = pgtype.Timestamp{
			Time: ping.LastSuccess,
		}
	}
	sql := `
			SELECT id FROM containers
			WHERE ip=$1;
	`
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
	row := tx.QueryRow(ctx, sql, ping.Ip)
	err := row.Scan(&ping.ContainerId)
	if err != nil {
		return err
	}
	row = tx.QueryRow(ctx, sql1, ping.ContainerId, ping.Latency.Nanoseconds(), lastSuccessAt, ping.PingAt)
	err = row.Scan(&ping.Id)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, sql2, ping.Id, ping.Latency.Nanoseconds(), lastSuccessAt, ping.PingAt, ping.ContainerId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 0 {
		return nil
	}
	_, err = tx.Exec(ctx, sql3, ping.Id, ping.Latency.Nanoseconds(), lastSuccessAt, ping.PingAt, ping.ContainerId)
	if err != nil {
		return err
	}
	return nil
}
