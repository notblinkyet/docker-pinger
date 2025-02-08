package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/internal/models"
)

var (
	ErrIpAlreadyTracked error = errors.New("container with this ip is already tracke")
)

type IContainerStorage interface {
	Create(ip string) error
	Delete(ip string) error
	GetAll() ([]models.Container, error)
}

type ContainerStorage struct {
	pool *pgxpool.Pool
}

func NewContainerStorage(pool *pgxpool.Pool) *ContainerStorage {
	return &ContainerStorage{
		pool: pool,
	}
}

func (storage *ContainerStorage) GetByIpInTx(tx pgx.Tx, ip string) (*models.Container, error) {

	sql := `SELECT id, ip FROM container WHERE ip = $1;`
	var container models.Container
	err := tx.QueryRow(context.Background(), sql, ip).Scan(&container.Id, &container.Ip)
	if err != nil {
		//TODO logs
		return nil, err
	}
	return &container, nil
}

func (storage *ContainerStorage) Create(ip string) error {
	const op = "storage/container/create"

	tx, err := storage.pool.Begin(context.Background())
	if err != nil {
		//TODO logs
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			//TODO logs
			tx.Rollback(context.Background())
		}
	}()

	_, err = storage.GetByIpInTx(tx, ip)

	if err == nil {
		//TODO logs
		return ErrIpAlreadyTracked
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		//TODO logs
		return fmt.Errorf("%s:%w", op, err)
	}

	sql := `INSERT INTO container(ip)
			VALUES ($1);
	`
	_, err = tx.Exec(context.Background(), sql, ip)
	if err != nil {
		//TODO logs
		return fmt.Errorf("%s:%w", op, err)
	}
	if err = tx.Commit(context.Background()); err != nil {
		//TODO logs
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (storage *ContainerStorage) Delete(ip string) error {
	const op = "storage/container/delete"

	tx, err := storage.pool.Begin(context.Background())
	if err != nil {
		//TODO logs
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			//TODO logs
			tx.Rollback(context.Background())
		}
	}()

	_, err = storage.GetByIpInTx(tx, ip)

	if err == nil {
		//TODO logs
		return ErrIpAlreadyTracked
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		//TODO logs
		return fmt.Errorf("%s:%w", op, err)
	}

	sql := `DELETE FROM container
		WHERE id=$1;
	`
	_, err = tx.Exec(context.Background(), sql, ip)
	if err != nil {
		//TODO logs
		return fmt.Errorf("%s:%w", op, err)
	}
	if err = tx.Commit(context.Background()); err != nil {
		//TODO logs
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (storage *ContainerStorage) GetAll() ([]models.Container, error) {
	const op = "storage/container/getall"
	containers := make([]models.Container, 0, 10)
	sql := `SELECT id, ip FROM container;`
	rows, err := storage.pool.Query(context.Background(), sql)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var container models.Container
		err = rows.Scan(
			&container.Id,
			&container.Ip,
		)
		if err != nil {
			//TODO logs
			continue
		}
		containers = append(containers, container)
	}
	return containers, nil
}
