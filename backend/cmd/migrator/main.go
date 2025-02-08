package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/notblinkyet/docker-pinger/backend/internal/config"
)

func main() {
	var d bool

	cfg := config.MustLoad()
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Storage.Username,
		os.Getenv("POSTGRES_PASS"), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)
	migrationsPath := cfg.MigrationPath

	flag.BoolVar(&d, "d", false, "use down migrations")
	flag.Parse()

	if dbURL == "" {
		panic("db-url is required")
	}
	if migrationsPath == "" {
		panic("migrations path is required")
	}
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dbURL)

	if err != nil {
		panic(fmt.Errorf("failed to  create migration engine: %v", err))
	}

	if d {
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}
			panic(err)
		}
		fmt.Println("migrations down applied successfully")
		return
	}
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations up applied successfully")

}
