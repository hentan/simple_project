package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simple_project/internal/config"
	"simple_project/internal/models"
	"simple_project/internal/repository"
	"time"
)

type Application struct {
	DB     repository.Database
	config config.Config
}

type Handler interface {
	Start(ctx context.Context, h http.Handler) error
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

func (app *Application) Start(ctx context.Context, h http.Handler) error {
	addr := app.config.AppPort
	log.Printf("Starting server on %s\n", addr)

	server := &http.Server{
		Addr:              ":" + addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	log.Println("server stopped")
	return nil
}

func New(db repository.Database, cfg config.Config) Handler {
	return &Application{
		DB:     db,
		config: cfg,
	}
}

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

func (app *Application) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := app.DB.GetAllPayments(r.Context())
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
	err := app.DB.AddPayment(r.Context(), &payment)
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
	err := app.DB.UpdatePayment(r.Context(), &payment)
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
	err := app.DB.DeletePayment(r.Context(), &payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete payment to db error: %v", err)
		return
	}
}

func (app *Application) GetBalance(ctx context.Context) (int, error) {
	sumExpenses, err := app.DB.GetSumExpenses(ctx)
	if err != nil {
		return 0, fmt.Errorf("get sum expenses: %w", err)
	}

	sumPayments, err := app.DB.GetSumPayments(ctx)
	if err != nil {
		return 0, fmt.Errorf("get sum payments: %w", err)
	}

	summBalance := sumPayments - sumExpenses
	return summBalance, nil
}

func (app *Application) GetPupils(w http.ResponseWriter, r *http.Request) {
	pupils, err := app.DB.GetPupils(r.Context())
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
	err := app.DB.AddPupil(r.Context(), &pupil)
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
	err := app.DB.UpdatePupil(r.Context(), &pupil)
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
	err := app.DB.DeletePupil(r.Context(), &pupil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("delete pupil to db error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}
