package models

type ConfigServer struct {
	Host string
	Port string
}

type ConfigPostgresRepository struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

type ConfigClickhouseRepository struct {
	Host string
	Port string
	Name string
}