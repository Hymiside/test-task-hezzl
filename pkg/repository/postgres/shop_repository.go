package postgres

import (
	"context"
	"database/sql"
	"strconv"
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

	var res models.Good
	if err := tx.QueryRowContext(ctx, "INSERT INTO goods (name, project_id) VALUES ($1, $2) RETURNING *", data.Name, data.ProjectId).Scan(
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
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
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
	
	var res models.Good
	if err = tx.QueryRowContext(ctx, "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4 RETURNING *", data.Name, data.Description, data.Id, data.ProjectId).Scan(
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
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}

	if err = tx.Commit(); err != nil {
        return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
    }
	return res, nil
}

func (s *shopPostgres) Delete(data models.Good) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := s.dbP.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error to begin transaction: %v", err)
	}

	var (
		id, projectId int
		removed bool
		res map[string]string
	)

	if err = tx.QueryRowContext(ctx, "UPDATE goods SET removed = $1 WHERE id = $2 AND project_id = $3 RETURNING id, project_id, removed", true, data.Id, data.ProjectId).Scan(&id, &projectId, &removed); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("good not found: %v", err)
		}
		return nil, fmt.Errorf("error to delete good: %v", err)
	}

	res = map[string]string{
		"id": strconv.Itoa(id),
		"projectId": strconv.Itoa(projectId),
		"removed": strconv.FormatBool(removed),
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error to commit transaction: %v", err)
	}

	return res, nil
}