package main

import (
	"context"
	"log"
	"os/signal"
	"simple_project/internal/config"
	"simple_project/internal/handlers"
	"simple_project/internal/repository"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	envFilePath := ".env"
	cfg := config.NewConfig(envFilePath)

	pool, err := repository.NewPool(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("create postgresql pool: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	handler := handlers.New(repo, *cfg)

	if err := handler.Start(ctx, handlers.Routes(handler)); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
