package clickhouse

import (
	"context"
	"fmt"
	"database/sql"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	_ "github.com/mailru/go-clickhouse/v2"
)


func NewClickhouseDB(ctx context.Context, c models.ConfigClickhouseRepository) (*sql.DB, error) {
	dbC, err := sql.Open("chhttp", fmt.Sprintf("http://%s:%s/%s", c.Host, c.Port, c.Name))
	if err != nil {
		return nil, fmt.Errorf("error to connection clickhouse: %v", err)
	}
	if err := dbC.Ping(); err != nil {
		return nil, fmt.Errorf("error to test connection clickhouse: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		dbC.Close()
	}(ctx)

	return dbC, nil
}