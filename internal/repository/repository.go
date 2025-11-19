package repository

import "database/sql"

type Database interface {
	Connect() *sql.DB
	GetExpenses()
}
