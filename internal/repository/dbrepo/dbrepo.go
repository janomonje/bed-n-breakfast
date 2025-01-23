package dbrepo

import (
	"database/sql"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	"github.com/janomonje/bed-n-breakfast/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPosgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
