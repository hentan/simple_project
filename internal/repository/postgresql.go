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

	query := `SELECT id, date, gift_for, surname, summ 
			  FROM expenses
			  ORDER BY surname
		`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses: %w", err)
	}
	defer rows.Close()
	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.Id,
			&expense.Date,
			&expense.GiftFor,
			&expense.Surname,
			&expense.Sum)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (repo *PostgresqlRepository) AddExpense(expense models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO expenses(date, gift_for, surname, summ)
			  VALUES ($1, $2, $3, $4)`

	err := repo.db.QueryRowContext(ctx, query,
		expense.Date,
		expense.GiftFor,
		expense.Surname,
		expense.Sum)
	if err != nil {
		return fmt.Errorf("cannot insert expense: %w", err)
	}
	return nil

}
