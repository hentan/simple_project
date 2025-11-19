package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"simple_project/internal/config"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	fmt.Println(ctx)
	envFilePath := ".env"
	cfg := config.NewConfig(envFilePath)
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(":8080", router)
}
