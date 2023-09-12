package clickhouse

import (
	"database/sql"
)


type logsClickhouse struct {
	dbC *sql.DB
}

func newLogsClickhouse(db *sql.DB) *logsClickhouse {
	return &logsClickhouse{dbC: db}
}