package repository

import (
	"database/sql"
	"simple_project/internal/models"
)

type Database interface {
	Connect() *sql.DB
	GetExpenses() ([]models.Expense, error)
	AddExpense(expense *models.Expense) error
	UpdateExpense(expense *models.Expense) error
	DeleteExpense(expense *models.Expense) error
	AddPayment(payment *models.Payment) error
	UpdatePayment(payment *models.Payment) error
	DeletePayment(payment *models.Payment) error
	GetAllPayments() ([]models.Payment, error)
}
