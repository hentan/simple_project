package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simple_project/internal/config"
	"simple_project/internal/repository"
)

type Application struct {
	DB     repository.Database
	config config.Config
}

type Handler interface {
	Start(h http.Handler) error
	GetExpenses(w http.ResponseWriter, r *http.Request)
}

func (app *Application) Start(h http.Handler) error {
	addr := app.config.AppPort
	log.Printf("Starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, h); err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
