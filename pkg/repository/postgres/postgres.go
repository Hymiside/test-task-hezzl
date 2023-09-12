package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	_ "github.com/lib/pq"
)

func NewPostgresDB(ctx context.Context, c models.ConfigPostgresRepository) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to connection postgres: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}
	return db, nil
}