package postgres

import (
	"database/sql"
)


type shopPostgres struct {
	dbP *sql.DB
}

func newShopPostgres(db *sql.DB) *shopPostgres {
	return &shopPostgres{dbP: db}
}