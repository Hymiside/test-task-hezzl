package postgres

import (
	"database/sql"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
)

type shop interface{
	Create(data models.Good) (models.Good, error)
	Update(data models.Good) (models.Good, error)
	Delete(data models.Good) (map[string]string, error)
}

type PostgresRepository struct {
	Shop shop
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{Shop: newShopPostgres(db)}
}