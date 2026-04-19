package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"simple_project/internal/models"
)

func (app *Application) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := app.DB.GetPayments(r.Context())
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
