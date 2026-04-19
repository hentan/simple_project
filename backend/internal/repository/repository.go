package repository

import (
	"context"
	"simple_project/internal/models"
)

type Store interface {
	GetExpenses(ctx context.Context) ([]models.Expense, error)
	GetExpensesArchive(ctx context.Context) ([]models.ExpenseArchive, error)
	AddExpense(ctx context.Context, expense *models.Expense) error
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	DeleteExpense(ctx context.Context, expense *models.Expense) error
	AddPayment(ctx context.Context, payment *models.Payment) error
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	DeletePayment(ctx context.Context, payment *models.Payment) error
	GetPayments(ctx context.Context) ([]models.Payment, error)
	GetPaymentsTotal(ctx context.Context) (int, error)
	GetExpensesTotal(ctx context.Context) (int, error)
	GetPupils(ctx context.Context) ([]models.Pupil, error)
	AddPupil(ctx context.Context, pupil *models.Pupil) error
	UpdatePupil(ctx context.Context, pupil *models.Pupil) error
	DeletePupil(ctx context.Context, pupil *models.Pupil) error
}
