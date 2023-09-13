package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"strconv"
	"strings"
    "encoding/json"

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
	if err := tx.QueryRowContext(ctx, "INSERT INTO goods (name, project_id) VALUES ($1, $2) RETURNING id, project_id, name, description, priority, removed, created_at", 
	data.Name, data.ProjectId).Scan(
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
			return models.Good{}, fmt.Errorf("error to get description good: %v", err)
		}
	}

	if err := tx.QueryRowContext(ctx, "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4",  data.Name, data.Description, data.Id, data.ProjectId); err.Err() != nil {
		return models.Good{}, fmt.Errorf("error to update good: %v", err.Err())
	}

	if err = tx.Commit(); err != nil {
		return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
	}
	return data, nil
}

func (s *shopPostgres) Delete(data models.Good) (models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := s.dbP.BeginTx(ctx, nil)
	if err != nil {
		return models.Good{}, fmt.Errorf("error to begin transaction: %v", err)
	}

	var res models.Good
	if err := tx.QueryRowContext(ctx, "UPDATE goods SET removed = $1 WHERE id = $2 AND project_id = $3 RETURNING id, project_id, name, description, priority, removed, created_at", 
	true, data.Id, data.ProjectId).Scan(
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
		return models.Good{}, fmt.Errorf("error to delete good: %v", err)
	}
	
	if err = tx.Commit(); err != nil {
		return models.Good{}, fmt.Errorf("error to commit transaction: %v", err)
	}
	
	return res, nil
}

func (s *shopPostgres) GetAll(limit, offset int) ([]models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var req string
	if limit == 0 && offset == 0 {
		req = "SELECT * FROM goods;"
	} else {
		req = fmt.Sprintf("SELECT * FROM goods LIMIT %d OFFSET %d;", limit, offset)
	}

	rows, err := s.dbP.QueryContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error to get all goods: %v", err)
	}
	defer rows.Close()

	var goods []models.Good
	for rows.Next() {
		var good models.Good
		if err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt); err != nil {
			return nil, fmt.Errorf("error to scan good: %v", err)
		}
		goods = append(goods, good)
	}
	return goods, nil
}

func (s *shopPostgres) WriteLogs(logs [][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var (
		good models.Good
		data []models.Good

	)

	for _, byteGood := range logs {
		if err := json.Unmarshal(byteGood, &good); err != nil {fmt.Printf("error to unmarshal byteGood: %v", err)}
		data = append(data, good)
	}
	vals := []interface{}{}
	for _, row := range data {
		vals = append(vals, row.Id, row.ProjectId, row.Name, row.Description, row.Priority, row.Removed, row.CreatedAt)
	}

	sqlStr := `INSERT INTO logs(good_id, project_id, name, description, priority, removed, created_at) VALUES %s`
	sqlStr = s.replaceSQL(sqlStr, "(?, ?, ?, ?, ?, ?, ?)", len(data))

	tx, err := s.dbP.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to begin transaction: %v", err)
	}
	defer tx.Rollback()

	if err := tx.QueryRowContext(ctx, sqlStr, vals...); err.Err() != nil {
		return fmt.Errorf("error to bulk insert logs: %v", err.Err())
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error to commit transaction: %v", err)
	}
	
	return nil
}

func (s *shopPostgres) replaceSQL(stmt, pattern string, len int) string {
	pattern += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(stmt, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		stmt = strings.Replace(stmt, "?", param, 1)
	}
	return strings.TrimSuffix(stmt, ",")
}