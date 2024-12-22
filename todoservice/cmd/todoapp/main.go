package main

import (
	"context"
	"fmt"
	database "github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
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
		fmt.Printf("Error on setup logger config %+v", err)
		return
	}
	var dbConfig database.Config
	if err := envconfig.Process("", &dbConfig); err != nil {
		fmt.Printf("Error on setup databse config %+v", err)
		return
	}
	ctx := context.Background()
	ctx, err := log.SetupLogger(ctx, cfg)
	if err != nil {
		fmt.Printf("Error on setup %+v", err)
		return
	}
	db, err := database.Create(ctx, dbConfig)
	if err != nil {
		log.C(ctx).Fatal(err)
		return
	}
	var oauth2Config jwt.ConfigOAuth2
	if err = envconfig.Process("", &oauth2Config); err != nil {
		fmt.Printf("Error on setup oauth2 config %+v", err)
	}
	restServer := http.NewServer(db, oauth2Config)
	restServer.Start()
}
