package ping

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

func (storage *PingStorage) Get(from, op string) ([]models.Ping, error) {
	ctx := context.Background()
	sql := fmt.Sprintf(`
		SELECT c.ip, p.id, p.latency, p.last_success_at, p.ping_at
		FROM %s p JOIN containers c ON p.container_id=c.id
	`, from)
	rows, err := storage.pool.Query(ctx, sql)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	pings := make([]models.Ping, 0, 10)
	for rows.Next() {
		var ping models.Ping
		var lastSuccessAt pgtype.Timestamp
		var latency int
		err = rows.Scan(
			&ping.Ip,
			&ping.Id,
			&latency,
			&lastSuccessAt,
			&ping.PingAt,
		)
		if err != nil {
			//TODO add logs
			continue
		}
		if lastSuccessAt.Status == pgtype.Null {
			ping.WasSuccessBefore = false
		} else {
			ping.WasSuccessBefore = true
			ping.LastSuccess = lastSuccessAt.Time
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
