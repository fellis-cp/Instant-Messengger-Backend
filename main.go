package main

import (
	"context"
	"instant-messenger-backend/database"
	"instant-messenger-backend/rabbitmq"
	"instant-messenger-backend/routes"
	"instant-messenger-backend/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.New()

	// Connect to RabbitMQ
	rabbitmq.ConnectToRabbitMQ()

	// Connect to the database
	database.ConnectDB()

	// Setup routes
	routes.SetupRouter(r)

	// Start WebSocket server
	go websocket.StartWS()

	// Start the HTTP server
	srv := &http.Server{
		Addr:    "localhost:5000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run HTTP server: %s", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down the server...")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
