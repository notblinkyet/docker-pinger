package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

var (
	ErrIpAlreadyTracked error = errors.New("container with this ip is already tracked")
	ErrNotExist         error = errors.New("this ip don't tracked")
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

	sql := `SELECT id, ip, is_tracked FROM containers WHERE ip = $1;`
	var container models.Container
	err := tx.QueryRow(context.Background(), sql, ip).Scan(&container.Id, &container.Ip, &container.IsTracked)
	if err != nil {
		//TODO logs
		return nil, err
	}
	return &container, nil
}

func (storage *ContainerStorage) Create(ip string) error {
	const op = "storage/container/create"
	ctx := context.Background()

	tx, err := storage.pool.Begin(context.Background())
	if err != nil {
		//TODO logs
		return fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback(ctx)

	container, err := storage.GetByIpInTx(tx, ip)
	fmt.Println(container)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		//TODO logs
		return err
	} else if container != nil && !container.IsTracked {
		sql :=
			`
          UPDATE containers
          SET is_tracked = TRUE
          WHERE id = $1;
        `
		fmt.Printf("Updating container ID: %d, is_tracked: %t\n", container.Id, true)

		result, err := tx.Exec(ctx, sql, container.Id)
		if err != nil {
			//TODO logs
			return fmt.Errorf("%s:%w", op, err)
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			//TODO logs
			return fmt.Errorf("%s: no rows updated", op)
		}
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	} else if container != nil {
		//TODO logs
		return fmt.Errorf("%s:%w", op, ErrIpAlreadyTracked)
	}

	sql := `INSERT INTO containers(ip, is_tracked)
			VALUES ($1, $2);
	`
	_, err = tx.Exec(context.Background(), sql, ip, true)
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

	container, err := storage.GetByIpInTx(tx, ip)

	if err != nil {
		//TODO logs
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s:%w", op, ErrNotExist)
		}
		return fmt.Errorf("%s:%w", op, err)
	}

	sql := `
		UPDATE containers
		SET is_tracked=FALSE
		WHERE id=$1;
	`
	_, err = tx.Exec(context.Background(), sql, container.Id)
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
	sql := `
		SELECT id, ip
		FROM containers
		WHERE is_tracked=TRUE;
		`
	rows, err := storage.pool.Query(context.Background(), sql)
	if err != nil {
		//TODO logs
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var container models.Container
		container.IsTracked = true
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
