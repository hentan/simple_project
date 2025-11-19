package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"simple_project/internal/config"
	"simple_project/internal/handlers"
	"simple_project/internal/repository"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	fmt.Println(ctx)
	envFilePath := ".env"
	cfg := config.NewConfig(envFilePath)
	repo := repository.New(cfg)
	handler := handlers.New(repo)
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(":8080", router)
}
