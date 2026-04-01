package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simple_project/internal/config"
	"simple_project/internal/models"
	"simple_project/internal/repository"
)

type Application struct {
	DB     repository.Database
	config config.Config
}

type Handler interface {
	Start(h http.Handler) error
	GetExpenses(w http.ResponseWriter, r *http.Request)
	AddExpense(w http.ResponseWriter, r *http.Request)
	UpdateExpense(w http.ResponseWriter, r *http.Request)
	DeleteExpense(w http.ResponseWriter, r *http.Request)
	GetPayments(w http.ResponseWriter, r *http.Request)
	AddPayment(w http.ResponseWriter, r *http.Request)
	UpdatePayment(w http.ResponseWriter, r *http.Request)
	DeletePayment(w http.ResponseWriter, r *http.Request)
	GetPupils(w http.ResponseWriter, r *http.Request)
	AddPupil(w http.ResponseWriter, r *http.Request)
	UpdatePupil(w http.ResponseWriter, r *http.Request)
	DeletePupil(w http.ResponseWriter, r *http.Request)
}

func (app *Application) Start(h http.Handler) error {
	addr := app.config.AppPort
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(":"+addr, h); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func New(db repository.Database, cfg config.Config) Handler {
	return &Application{
		DB:     db,
		config: cfg,
	}
}

func (app *Application) GetExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := app.DB.GetExpenses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := app.GetBalance()
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

	err := app.DB.AddExpense(&expense)
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

	err := app.DB.UpdateExpense(&expense)
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

	err := app.DB.DeleteExpense(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete expense to db error: %v", err)
		return
	}
}

func (app *Application) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := app.DB.GetAllPayments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func (app *Application) AddPayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.AddPayment(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("add payment to db error: %v", err)
		return
	}
}

func (app *Application) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.UpdatePayment(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("update payment to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) DeletePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.DeletePayment(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete payment to db error: %v", err)
		return
	}
}

func (app *Application) GetBalance() (int, error) {
	sumExpenses, err := app.DB.GetSumExpenses()
	if err != nil {
		return 0, fmt.Errorf("get sum expenses: %w", err)
	}

	sumPayments, err := app.DB.GetSumPayments()
	if err != nil {
		return 0, fmt.Errorf("get sum payments: %w", err)
	}

	summBalance := sumExpenses - sumPayments
	return summBalance, nil
}

func (app *Application) GetPupils(w http.ResponseWriter, r *http.Request) {
	pupils, err := app.DB.GetPupils()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pupils); err != nil {
		log.Printf("json encode error: %v", err)
	}

}

func (app *Application) AddPupil(w http.ResponseWriter, r *http.Request) {
	var pupil models.Pupil
	if err := json.NewDecoder(r.Body).Decode(&pupil); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.AddPupil(&pupil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("add pupil to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func (app *Application) UpdatePupil(w http.ResponseWriter, r *http.Request) {
	var pupil models.Pupil
	if err := json.NewDecoder(r.Body).Decode(&pupil); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.UpdatePupil(&pupil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("update pupil to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}

func (app *Application) DeletePupil(w http.ResponseWriter, r *http.Request) {
	var pupil models.Pupil
	if err := json.NewDecoder(r.Body).Decode(&pupil); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("json decode error: %v", err)
		return
	}
	err := app.DB.DeletePupil(&pupil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete pupil to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}
