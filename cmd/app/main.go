package main

import (
	"context"
	"fmt"
	"log"
	"simple_project/internal/config"
	"simple_project/internal/handlers"
	"simple_project/internal/repository"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	fmt.Println(ctx)
	envFilePath := ".env"
	cfg := config.NewConfig(envFilePath)
	repo := repository.New(cfg.Postgresql)
	handler := handlers.New(repo, *cfg)
	err := handler.Start(handlers.Routes(handler))
	if err != nil {
		log.Fatal(err)
	}
}
