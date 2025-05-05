package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-htmx-light-starter/counterapp"
	"go-htmx-light-starter/internal/config" // Import the config package
	"go-htmx-light-starter/todonow"

	"github.com/go-chi/chi/v5" // Import chi for subrouter
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize database connection
	db, err := todonow.InitDB("todonow/data/todonow.db") // Use relative path for SQLite file
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	// Ensure the database connection is closed when main exits
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}()

	// Create a handlers instance, passing the database connection
	// Note: We will need to update handlers.NewHandlers to accept *bun.DB
	counterAppHandlers := counterapp.NewHandlers(db) // Counter app handlers

	// --- TodoNow Feature Setup ---
	// Create repository instance
	todoRepo := todonow.NewTodoRepository(db)
	// Create todo handler instance
	todoHandler := todonow.NewTodoHandler(todoRepo)
	// --- End TodoNow Feature Setup ---

	// Create the main router
	router := chi.NewRouter()

	// --- Serve Static Files ---
	// It's crucial to serve static files before mounting application routes
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	// --- End Static Files ---

	// --- Mount CounterApp Routes ---
	// Get the routes from the counterapp package and mount them
	router.Mount("/", counterAppHandlers.Routes()) // Mount counter app at root
	// --- End CounterApp Routes ---

	// --- Register TodoNow Routes ---
	router.Route("/todonow", func(r chi.Router) {
		todoHandler.RegisterTodoRoutes(r)
	})
	// --- End TodoNow Routes ---

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
