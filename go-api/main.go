package main

import (
	"log"
	"os"

	delivery "matrix-orchestrator/delivery/http"
	"matrix-orchestrator/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("FATAL ERROR: JWT_SECRET environment variable is not set")
	}
	if os.Getenv("NODE_API_URL") == "" {
		log.Fatal("FATAL ERROR: NODE_API_URL environment variable is not set")
	}
	if os.Getenv("AUTH_USERNAME") == "" || os.Getenv("AUTH_PASSWORD") == "" {
		log.Fatal("FATAL ERROR: AUTH_USERNAME or AUTH_PASSWORD environment variable is not set")
	}

	app := fiber.New(fiber.Config{
		AppName: "Matrix Orchestrator API v1.0",
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	matrixUsecase := usecase.NewMatrixUsecase()
	delivery.NewMatrixHandler(app, matrixUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Go Fiber server on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
