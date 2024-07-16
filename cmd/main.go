package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"nedorez-test/database"
	"nedorez-test/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	internal.Route(app, db)

	go func() {
		log.Fatal(app.Listen(":8888"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")
	app.ShutdownWithContext(context.Background())
}
