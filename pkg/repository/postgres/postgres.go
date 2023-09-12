package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/test-task-hezzl/pkg/models"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewPostgresDB(ctx context.Context, c models.ConfigPostgresRepository) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	dbP, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to connection postgres: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		dbP.Close()
	}(ctx)

	if err = dbP.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	m, err := migrate.New("file://migrations", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil {
		return dbP, nil
		// return nil, fmt.Errorf("error to make migrations postgres: %w", err)
	}

	return dbP, nil
}