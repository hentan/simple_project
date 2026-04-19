package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"simple_project/internal/models"
)

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
}
