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

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func NewPool(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pc, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres dsn: %w", err)
	}

	pc.MaxConns = 20
	pc.MinConns = 2
	pc.MaxConnLifetime = 30 * time.Minute
	pc.MaxConnIdleTime = 5 * time.Minute
	pc.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return pool, nil
}

func (repo *PostgresRepository) Close() {
	repo.pool.Close()
}

func (repo *PostgresRepository) GetExpenses(ctx context.Context) ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
			&expense.ID,
			&expense.Date,
			&expense.Purpose,
			&expense.PupilID,
			&expense.Surname,
			&expense.Amount,
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

func (repo *PostgresRepository) GetExpensesArchive(ctx context.Context) ([]models.ExpenseArchive, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			ea.archived_at,
			ea.op,
			ea.id,
			ea.date,
			ea.gift_for,
			ea.pupil_id,
			COALESCE(p.surname, '') AS surname,
			ea.summ
		FROM expenses_archive ea
		LEFT JOIN pupils p ON p.id = ea.pupil_id
		ORDER BY ea.archived_at DESC, ea.id DESC
	`

	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses archive: %w", err)
	}
	defer rows.Close()

	expenses := make([]models.ExpenseArchive, 0)
	for rows.Next() {
		var expense models.ExpenseArchive
		if err := rows.Scan(
			&expense.ArchivedAt,
			&expense.Operation,
			&expense.ID,
			&expense.Date,
			&expense.Purpose,
			&expense.PupilID,
			&expense.Surname,
			&expense.Amount,
		); err != nil {
			return nil, fmt.Errorf("could not scan expense archive row: %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return expenses, nil
}

func (repo *PostgresRepository) AddExpense(ctx context.Context, expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO expenses(date, gift_for, pupil_id, summ)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repo.pool.Exec(ctx, query,
		expense.Date,
		expense.Purpose,
		expense.PupilID,
		expense.Amount,
	)
	if err != nil {
		return fmt.Errorf("cannot insert expense: %w", err)
	}
	return nil
}

func (repo *PostgresRepository) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE expenses
		SET date = $1, gift_for = $2, pupil_id = $3, summ = $4
		WHERE id = $5
	`

	res, err := repo.pool.Exec(ctx, query,
		expense.Date,
		expense.Purpose,
		expense.PupilID,
		expense.Amount,
		expense.ID,
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

func (repo *PostgresRepository) DeleteExpense(ctx context.Context, expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM expenses WHERE id = $1`
	res, err := repo.pool.Exec(ctx, query, expense.ID)
	if err != nil {
		return fmt.Errorf("cannot delete expense: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresRepository) AddPayment(ctx context.Context, payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO payments(pupil_id, summ, date, purpose) VALUES ($1, $2, $3, $4)`
	_, err := repo.pool.Exec(ctx, query, payment.PupilID, payment.Amount, payment.Date, payment.Purpose)
	if err != nil {
		return fmt.Errorf("cannot insert payment: %w", err)
	}
	return nil
}

func (repo *PostgresRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE payments SET pupil_id = $1, summ = $2, date = $3, purpose = $4 WHERE id = $5`
	res, err := repo.pool.Exec(ctx, query, payment.PupilID, payment.Amount, payment.Date, payment.Purpose, payment.ID)
	if err != nil {
		return fmt.Errorf("cannot update payment: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresRepository) DeletePayment(ctx context.Context, payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM payments WHERE id = $1`
	res, err := repo.pool.Exec(ctx, query, payment.ID)
	if err != nil {
		return fmt.Errorf("cannot delete payment: %w", err)
	}
	aff := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *PostgresRepository) GetPayments(ctx context.Context) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
			&pmt.ID,
			&pmt.Date,
			&pmt.Purpose,
			&pmt.PupilID,
			&pmt.Surname,
			&pmt.Amount,
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

func (repo *PostgresRepository) GetPaymentsTotal(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT COALESCE(SUM(summ), 0) FROM payments`

	var total int
	err := repo.pool.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("could not get sum payments: %w", err)
	}
	return total, nil
}

func (repo *PostgresRepository) GetExpensesTotal(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT COALESCE(SUM(summ), 0) FROM expenses`

	var total int
	if err := repo.pool.QueryRow(ctx, query).Scan(&total); err != nil {
		return 0, fmt.Errorf("could not get sum expenses: %w", err)
	}
	return total, nil
}

func (repo *PostgresRepository) GetPupils(ctx context.Context) ([]models.Pupil, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func (repo *PostgresRepository) AddPupil(ctx context.Context, pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func (repo *PostgresRepository) UpdatePupil(ctx context.Context, pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func (repo *PostgresRepository) DeletePupil(ctx context.Context, pupil *models.Pupil) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
