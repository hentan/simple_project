package handlers

import (
	"fmt"
	"log"
	"net/http"
	"simple_project/internal/config"
)

type Application struct {
	config config.Config
}

type Handler interface {
	Start(h http.Handler)
	GetExpenses(w http.ResponseWriter, r *http.Request)
}

func (app *Application) Start(h http.Handler) error {
	err := http.ListenAndServe(app.config.AppPort, h)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("Listening on port %s", app.config.AppPort)
	log.Println(msg)
	return nil
}
