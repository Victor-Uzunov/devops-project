package main

import (
	"context"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/server"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error on setup %+v", err)
		return
	}
	var cfg log.Config
	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Printf("Error on setup %+v", err)
		return
	}
	var apiConfig client.APIConfig
	err := envconfig.Process("", &apiConfig)
	if err != nil {
		fmt.Printf("Failed to load API config: %v", err)
		return
	}
	ctx := context.Background()
	ctx, err = log.SetupLogger(ctx, cfg)
	if err != nil {
		fmt.Printf("Error on setup %+v", err)
		return
	}
	s := server.NewServer(apiConfig)
	s.Start()
}
