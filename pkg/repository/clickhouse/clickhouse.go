package clickhouse

import (
	"context"
	"fmt"
	"database/sql"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
)


func NewClickhouseDB(ctx context.Context, c models.ConfigClickhouseRepository) (*sql.DB, error) {
	db, err := sql.Open("chhttp", fmt.Sprintf("http://%s:%s/%s", c.Host, c.Port, c.Name))
	if err != nil {
		return nil, fmt.Errorf("error to connection clickhouse: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error to test connection clickhouse: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	return db, nil
}