package models

type ExpenseWithBalance struct {
	Expenses []Expense `json:"expenses"`
	Balance  int       `json:"balance"`
}
