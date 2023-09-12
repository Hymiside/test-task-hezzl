package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
)


type shopPostgres struct {
	dbP *sql.DB
}

func newShopPostgres(db *sql.DB) *shopPostgres {
	return &shopPostgres{dbP: db}
}

func (s *shopPostgres) Create(data models.Good) (models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := s.dbP.BeginTx(ctx, nil)
    if err != nil {
        return models.Good{}, fmt.Errorf("error to begin transaction: %v", err)
    }
	defer tx.Rollback()

	var goodId int
	if err := tx.QueryRowContext(ctx, "INSERT INTO goods (name, project_id) VALUES ($1, $2) RETURNING id", data.Name, data.ProjectId).Scan(&goodId); err != nil {
		if err == sql.ErrNoRows {
            return models.Good{}, fmt.Errorf("good not found: %v", err)
        }
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}

	var res models.Good
	if err = tx.QueryRowContext(ctx, "SELECT * FROM goods WHERE id = $1", goodId).Scan(
		&res.Id, 
		&res.ProjectId, 
		&res.Name, 
		&res.Description, 
		&res.Priority, 
		&res.Removed, 
		&res.CreatedAt); err != nil {
        if err == sql.ErrNoRows {
            return models.Good{}, fmt.Errorf("good not found: %v", err)
        }
        return models.Good{}, fmt.Errorf("error to get good: %v", err)
	}

	if err = tx.Commit(); err != nil {
        return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
    }
	return res, nil
}

func (s *shopPostgres) Update(data models.Good) (models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := s.dbP.BeginTx(ctx, nil)
	if err != nil {
		return models.Good{}, fmt.Errorf("error to begin transaction: %v", err)
	}
	defer tx.Rollback()

	if data.Description == "" {
		if err = tx.QueryRowContext(ctx, "SELECT description FROM goods WHERE id = $1", data.Id).Scan(&data.Description); err != nil {
			if err == sql.ErrNoRows {
				return models.Good{}, fmt.Errorf("good not found: %v", err)
			}
			return models.Good{}, fmt.Errorf("error to get description good: %v", err)
		}
	} 
	
	var goodId int
	if err = tx.QueryRowContext(ctx, "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4 RETURNING id", data.Name, data.Description, data.Id, data.ProjectId).Scan(&goodId); err != nil {
		if err == sql.ErrNoRows {
            return models.Good{}, fmt.Errorf("good not found: %v", err)
        }
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}

	var res models.Good
	if err = tx.QueryRowContext(ctx, "SELECT * FROM goods WHERE id = $1", goodId).Scan(
		&res.Id, 
		&res.ProjectId, 
		&res.Name, 
		&res.Description, 
		&res.Priority, 
		&res.Removed, 
		&res.CreatedAt); err != nil {
        if err == sql.ErrNoRows {
            return models.Good{}, fmt.Errorf("good not found: %v", err)
        }
        return models.Good{}, fmt.Errorf("error to get good: %v", err)
	}

	if err = tx.Commit(); err != nil {
        return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
    }
	return res, nil
}