package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple_project/internal/config"
	"simple_project/internal/models"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresqlRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresqlRepository {
	return &PostgresqlRepository{pool: pool}
}

func NewPool(ctx context.Context, cfg config.Postgresql) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pc, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgresql dsn: %w", err)
	}

	pc.MaxConns = 20
	pc.MinConns = 2
	pc.MaxConnLifetime = 30 * time.Minute
	pc.MaxConnIdleTime = 5 * time.Minute
	pc.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgresql database: %w", err)
	}

	return pool, nil
}

func (repo *PostgresqlRepository) Close() {
	repo.pool.Close()
}

func (repo *PostgresqlRepository) GetExpenses() ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.date, e.gift_for, e.pupil_id, p.surname, e.summ
		FROM expenses e
		INNER JOIN pupils p ON p.id = e.pupil_id
		ORDER BY e.date DESC
	`

	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses: %w", err)
	}
	defer rows.Close()

	expenses := make([]models.Expense, 0)
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

	_, err := repo.pool.Exec(ctx, query,
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

	res, err := repo.pool.Exec(ctx, query,
		expense.Date,
		expense.GiftFor,
		expense.PupilId,
		expense.Summ,
		expense.Id,
	)
	if err != nil {
		return fmt.Errorf("cannot update expense: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 { // FIX: полезно явно сигнализировать “ничего не обновилось”
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) DeleteExpense(expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM expenses WHERE id = $1`
	res, err := repo.pool.Exec(ctx, query, expense.Id)
	if err != nil {
		return fmt.Errorf("cannot delete expense: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) AddPayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO payments(pupil_id, summ, date, purpose) VALUES ($1, $2, $3, $4)`
	_, err := repo.pool.Exec(ctx, query, payment.PupilId, payment.Summ, payment.Date, payment.Purpose)
	if err != nil {
		return fmt.Errorf("cannot insert payment: %w", err)
	}
	return nil
}

func (repo *PostgresqlRepository) UpdatePayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE payments SET pupil_id = $1, summ = $2, date = $3, purpose = $4 WHERE id = $5`
	res, err := repo.pool.Exec(ctx, query, payment.PupilId, payment.Summ, payment.Date, payment.Purpose, payment.Id)
	if err != nil {
		return fmt.Errorf("cannot update payment: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) DeletePayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM payments WHERE id = $1`
	res, err := repo.pool.Exec(ctx, query, payment.Id)
	if err != nil {
		return fmt.Errorf("cannot delete payment: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) GetAllPayments() ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			pay.id,
			pay.date,
			pay.purpose,
			pay.pupil_id,
			COALESCE(p.surname, '') AS surname,
			pay.summ
		FROM payments pay
		LEFT JOIN pupils p ON p.id = pay.pupil_id
		ORDER BY p.surname
	`

	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var pmt models.Payment
		if err := rows.Scan(
			&pmt.Id,
			&pmt.Date,
			&pmt.Purpose,
			&pmt.PupilId,
			&pmt.Surname,
			&pmt.Summ,
		); err != nil {
			return nil, fmt.Errorf("could not scan payment row: %w", err)
		}
		payments = append(payments, pmt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return payments, nil
}

func (repo *PostgresqlRepository) GetSumPayments() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT COALESCE(SUM(summ), 0) FROM payments`

	var total int
	err := repo.pool.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("could not get sum payments: %w", err)
	}
	return total, nil
}

func (repo *PostgresqlRepository) GetSumExpenses() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT COALESCE(SUM(summ), 0) FROM expenses`

	var total int
	if err := repo.pool.QueryRow(ctx, query).Scan(&total); err != nil {
		return 0, fmt.Errorf("could not get sum expenses: %w", err)
	}
	return total, nil
}

func (repo *PostgresqlRepository) GetPupils() ([]models.Pupil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, name, surname, parent_name, parent_phone
			  FROM pupils
			  ORDER BY surname`
	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get pupils: %w", err)
	}
	defer rows.Close()
	var pupils []models.Pupil
	for rows.Next() {
		var pupil models.Pupil
		if err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.Surname,
			&pupil.ParentName,
			&pupil.ParentPhone); err != nil {
			return nil, fmt.Errorf("could not scan pupil row: %w", err)
		}
		pupils = append(pupils, pupil)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return pupils, nil
}

func (repo *PostgresqlRepository) AddPupil(pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO pupils( name, surname, parent_name, parent_phone)
			  VALUES($1, $2, $3, $4)`

	_, err := repo.pool.Exec(ctx, query,
		pupil.Name,
		pupil.Surname,
		pupil.ParentName,
		pupil.ParentPhone)
	if err != nil {
		return fmt.Errorf("cannot insert pupil: %w", err)
	}
	return nil
}

func (repo *PostgresqlRepository) UpdatePupil(pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE pupils
			  SET name = $1, surname = $2, parent_name = $3, parent_phone = $4
			  WHERE id = $5`
	res, err := repo.pool.Exec(ctx, query, pupil.Name, pupil.Surname, pupil.ParentName, pupil.ParentPhone, pupil.ID)
	if err != nil {
		return fmt.Errorf("cannot update pupil: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresqlRepository) DeletePupil(pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `DELETE FROM pupils WHERE id = $1`
	res, err := repo.pool.Exec(ctx, query, pupil.ID)
	if err != nil {
		return fmt.Errorf("cannot delete pupil: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}
