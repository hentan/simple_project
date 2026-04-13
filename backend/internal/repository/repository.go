package repository

import (
	"context"
	"database/sql"
	"simple_project/internal/models"
)

type Database interface {
	Connect() *sql.DB
	GetExpenses(ctx context.Context) ([]models.Expense, error)
	AddExpense(ctx context.Context, expense *models.Expense) error
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	DeleteExpense(ctx context.Context, expense *models.Expense) error
	AddPayment(ctx context.Context, payment *models.Payment) error
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	DeletePayment(ctx context.Context, payment *models.Payment) error
	GetAllPayments(ctx context.Context) ([]models.Payment, error)
	GetSumPayments(ctx context.Context) (int, error)
	GetSumExpenses(ctx context.Context) (int, error)
	GetPupils(ctx context.Context) ([]models.Pupil, error)
	AddPupil(ctx context.Context, pupil *models.Pupil) error
	UpdatePupil(ctx context.Context, pupil *models.Pupil) error
	DeletePupil(ctx context.Context, pupil *models.Pupil) error
}
