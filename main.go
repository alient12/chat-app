package main

import (
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/handler"
	"chatapp/internal/infra/repository/memory"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := echo.New()

	var repo userrepo.Repository = memory.New()

	h := handler.NewUser(repo)
	h.Register(app.Group("/api"))

	if err := app.Start(":8000"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
