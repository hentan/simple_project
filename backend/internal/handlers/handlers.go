package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simple_project/internal/config"
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
