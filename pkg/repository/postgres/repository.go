package postgres

import (
	"database/sql"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
)

type shop interface {
	Create(data models.Good) (models.Good, error)
	Update(data models.Good) (models.Good, error)
	Delete(data models.Good) (models.Good, error)
	GetAll(limit, offset int) ([]models.Good, error)
	Reprioritiize(id, projectId, newPriority int) ([]models.Good, error)
	WriteLogs(logs [][]byte) error
}

type PostgresRepository struct {
	Shop shop
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{Shop: newShopPostgres(db)}
}
