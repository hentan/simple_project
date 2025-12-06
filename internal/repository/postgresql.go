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

func (repo *PostgresqlRepository) Connect() *sql.DB {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) UpdateExpense(expense *models.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) DeleteExpense(expense *models.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) AddPayment(payment *models.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) UpdatePayment(payment *models.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) DeletePayment(payment *models.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (repo *PostgresqlRepository) GetAllPayments() ([]models.Payment, error) {
	//TODO implement me
	panic("implement me")
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

	query := `SELECT e.id, e.date, e.gift_for, p.surname, e.summ 
			  FROM expenses e
			  INNER JOIN pupils p ON p.id = e.pupil_id
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

func (repo *PostgresqlRepository) AddExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
			  pupil_id := select pupil_id from pupils where id = $1;
			  INSERT INTO expenses(date, gift_for, surname, summ)
			  VALUES ($2, $3, pupil_id, $5)`

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

func UpdateExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE expenses set summ=$1 where`
}
