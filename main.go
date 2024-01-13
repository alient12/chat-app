package main

import (
	"chatapp/internal/infra/http/handler"
	"chatapp/internal/infra/repository/chatdb"
	"chatapp/internal/infra/repository/contactdb"
	"chatapp/internal/infra/repository/filedb"
	"chatapp/internal/infra/repository/messagedb"
	"chatapp/internal/infra/repository/userdb"
	"chatapp/internal/infra/websocket"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := echo.New()
	// app.Use(middleware.CORS())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	userRepo, err := userdb.New()
	if err != nil {
		log.Fatalf("cannot load users datavase")
	}
	chatRepo, err := chatdb.New()
	if err != nil {
		log.Fatalf("cannot load chats datavase")
	}
	contactRepo, err := contactdb.New()
	if err != nil {
		log.Fatalf("cannot load contacts datavase")
	}
	messageRepo, err := messagedb.New()
	if err != nil {
		log.Fatalf("cannot load messages datavase")
	}
	fileRepo, err := filedb.New()
	if err != nil {
		log.Fatalf("cannot load messages datavase")
	}

	userHand := handler.NewUser(userRepo)
	userHand.Register(app.Group("/api"))

	chatHand := handler.NewChat(chatRepo, userRepo)
	chatHand.Register(app.Group("/api"))

	contactHand := handler.NewContact(contactRepo, userRepo)
	contactHand.Register(app.Group("/api"))

	messageHand := handler.NewMessage(messageRepo, userRepo, chatRepo)
	messageHand.Register(app.Group("/api"))

	fileHand := handler.NewFile(fileRepo, userRepo, chatRepo)
	fileHand.Register(app.Group("/api"))

	WSHand := websocket.NewWebSocketConnection(messageHand)
	WSHand.Register(app.Group("/api"))

	// app.Use(middleware.Logger())
	// app.Use(middleware.Recover())

	if err := app.Start(":8000"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
	// if err := app.StartTLS(":8000", "localhost.crt", "localhost.key"); err != nil {
	// 	log.Fatalf("server failed to start %v", err)
	// }
}
