package dbrepo

import (
	"database/sql"

	repository "github.com/dsundar/bookings/internal/Repository"
	"github.com/dsundar/bookings/internal/config"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}
