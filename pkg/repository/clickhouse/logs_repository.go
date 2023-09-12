package clickhouse

import (
	"database/sql"
)


type logsClickhouse struct {
	db *sql.DB
}

func newLogsClickhouse(db *sql.DB) *logsClickhouse {
	return &logsClickhouse{db: db}
}