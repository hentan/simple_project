package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"simple_project/internal/config"
	"simple_project/internal/models"
	"time"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func New(cfg config.Postgresql) *PostgresqlRepository {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := connectToDB(connString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database")

	return &PostgresqlRepository{
		db: db,
	}
}

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("успешное подключение к БД!")
	return db, nil
}

func (repo *PostgresqlRepository) Connection() *sql.DB { return repo.db }

func (repo *PostgresqlRepository) GetExpenses() ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

}

func (repo *PostgresqlRepository) AddExpense(expense models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

}
