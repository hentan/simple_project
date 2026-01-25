package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"simple_project/internal/config"
	"simple_project/internal/models"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func (repo *PostgresqlRepository) Connect() *sql.DB {
	return repo.db
}

func New(cfg config.Postgresql) *PostgresqlRepository {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	fmt.Println(connString)
	db, err := connectToDB(connString)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	fmt.Println("Successfully connected to database")

	return &PostgresqlRepository{db: db}
}

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("успешное подключение к БД!")
	return db, nil
}

func (repo *PostgresqlRepository) Connection() *sql.DB { return repo.db }

func (repo *PostgresqlRepository) GetExpenses() ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.date, e.gift_for, e.pupil_id, p.surname, e.summ
		FROM expenses e
		INNER JOIN pupils p ON p.id = e.pupil_id
		ORDER BY p.surname
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses: %w", err)
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(
			&expense.Id,
			&expense.Date,
			&expense.GiftFor,
			&expense.PupilId,
			&expense.Surname,
			&expense.Summ,
		); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil { // FIX: корректная проверка ошибки итерации
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return expenses, nil
}

func (repo *PostgresqlRepository) AddExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO expenses(date, gift_for, pupil_id, summ)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repo.db.ExecContext(ctx, query,
		expense.Date,
		expense.GiftFor,
		expense.PupilId,
		expense.Summ,
	)
	if err != nil {
		return fmt.Errorf("cannot insert expense: %w", err)
	}
	return nil
}

func (repo *PostgresqlRepository) UpdateExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE expenses
		SET date = $1, gift_for = $2, pupil_id = $3, summ = $4
		WHERE id = $5
	`

	res, err := repo.db.ExecContext(ctx, query,
		expense.Date,
		expense.GiftFor,
		expense.PupilId,
		expense.Summ,
		expense.Id,
	)
	if err != nil {
		return fmt.Errorf("cannot update expense: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 { // FIX: полезно явно сигнализировать “ничего не обновилось”
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) DeleteExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// FIX: метод был TODO
	query := `DELETE FROM expenses WHERE id = $1`
	res, err := repo.db.ExecContext(ctx, query, expense.Id)
	if err != nil {
		return fmt.Errorf("cannot delete expense: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) AddPayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO payments(pupil_id, summ) VALUES ($1, $2)`
	_, err := repo.db.ExecContext(ctx, query, payment.PupilID, payment.Summ)
	if err != nil {
		return fmt.Errorf("cannot insert payment: %w", err)
	}
	return nil
}

func (repo *PostgresqlRepository) UpdatePayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE payments SET pupil_id = $1, summ = $2 WHERE id = $3`
	res, err := repo.db.ExecContext(ctx, query, payment.PupilID, payment.Summ, payment.Id)
	if err != nil {
		return fmt.Errorf("cannot update payment: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) DeletePayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM payments WHERE id = $1`
	res, err := repo.db.ExecContext(ctx, query, payment.Id)
	if err != nil {
		return fmt.Errorf("cannot delete payment: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) GetAllPayments() ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, pupil_id, summ FROM payments ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.Id, &p.PupilID, &p.Summ); err != nil {
			return nil, fmt.Errorf("could not scan payment row: %w", err)
		}
		payments = append(payments, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return payments, nil
}
