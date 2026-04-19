package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"simple_project/internal/models"
)

func (app *Application) GetExpenses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expenses, err := app.DB.GetExpenses(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := app.GetBalance(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	expenseWithBalance := models.ExpenseWithBalance{
		Expenses: expenses,
		Balance:  balance,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenseWithBalance); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func (app *Application) AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}

	err := app.DB.AddExpense(r.Context(), &expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("add expense to db error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}

	err := app.DB.UpdateExpense(r.Context(), &expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("update expense to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}

	err := app.DB.DeleteExpense(r.Context(), &expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete expense to db error: %v", err)
		return
	}
}
