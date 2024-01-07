package main

import (
	"chatapp/internal/infra/http/handler"
	"chatapp/internal/infra/repository/chatmem"
	"chatapp/internal/infra/repository/contactmem"
	"chatapp/internal/infra/repository/messagemem"
	"chatapp/internal/infra/repository/usermem"
	"chatapp/internal/infra/websocket"
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

	userRepo := usermem.New()
	chatRepo := chatmem.New()
	contactRepo := contactmem.New()
	messageRepo := messagemem.New()

	userHand := handler.NewUser(userRepo)
	userHand.Register(app.Group("/api"))

	chatHand := handler.NewChat(chatRepo, userRepo)
	chatHand.Register(app.Group("/api"))

	contactHand := handler.NewContact(contactRepo, userRepo)
	contactHand.Register(app.Group("/api"))

	messageHand := handler.NewMessage(messageRepo, userRepo, chatRepo)

	WSHand := websocket.NewWebSocketConnection(messageHand)
	WSHand.Register(app.Group("/api"))

	// app.Use(middleware.Logger())
	// app.Use(middleware.Recover())

	if err := app.Start(":8000"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
