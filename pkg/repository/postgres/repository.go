package postgres

import (
	"database/sql"
)

type shop interface{}

type PostgresRepository struct {
	sh shop
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{sh: newShopPostgres(db)}
}