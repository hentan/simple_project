package repository

import (
	"database/sql"
	"simple_project/internal/config"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func New(cfg config.Postgresql) *PostgresqlRepository {
	return &PostgresqlRepository{}
}
