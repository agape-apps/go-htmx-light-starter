package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/christian/mywebapp/internal/config"   // Import the config package
	"github.com/christian/mywebapp/internal/handlers" // Import the handlers package
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create a handlers instance which holds the application state
	appHandlers := handlers.NewHandlers()

	// Get the router with routes defined in the handlers package
	router := appHandlers.Routes()

	// Create a server with a timeout using the configured port
	serverAddr := ":" + cfg.Port
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router, // Use the router from the handlers package
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on http://localhost%s", serverAddr) // Use the configured address
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up channel to listen for signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigCh
	log.Printf("Got signal: %v, shutting down server...", sig)

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	log.Println("Waiting for active connections to complete...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed gracefully")
	}
}
