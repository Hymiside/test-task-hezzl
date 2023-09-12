package postgres

import (
	"database/sql"
)


type shopPostgres struct {
	db *sql.DB
}

func newShopPostgres(db *sql.DB) *shopPostgres {
	return &shopPostgres{db: db}
}