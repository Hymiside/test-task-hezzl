package clickhouse

import (
	"database/sql"
)

type logs interface{}

type ClickhouseRepository struct {
	l logs
}

func NewClickhouseRepository(db *sql.DB) *ClickhouseRepository {
	return &ClickhouseRepository{l: newLogsClickhouse(db)}
}