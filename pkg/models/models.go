package models

import (
	"time"
)

type ConfigServer struct {
	Host string
	Port string
}

type ConfigPostgresRepository struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ConfigClickhouseRepository struct {
	Host string
	Port string
	Name string
}

type ConfigRedis struct {
	Host string
	Port string
}

type Good struct {
	Id          int       `db:"id"`
	ProjectId   int       `db:"project_id"`
	Name        string    `json:"name" 		  db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

type ConfigNats struct {
	Host string
	Port string
}
