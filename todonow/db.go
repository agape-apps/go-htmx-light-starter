package todonow

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"go-htmx-light-starter/todonow/models" // Import models

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// InitDB initializes the database connection and ensures the necessary table exists.
func InitDB(dataSourceName string) (*bun.DB, error) {
	// Ensure the data directory exists
	dataDir := filepath.Dir(dataSourceName)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		log.Printf("Data directory '%s' does not exist, creating it...", dataDir)
		if err := os.MkdirAll(dataDir, 0755); err != nil { // Use 0755 for permissions
			log.Printf("Error creating data directory '%s': %v", dataDir, err)
			return nil, err
		}
		log.Printf("Data directory '%s' created successfully.", dataDir)
	} else if err != nil {
		log.Printf("Error checking data directory '%s': %v", dataDir, err)
		return nil, err // Return error if stat fails for other reasons
	}

	// Register the SQLite driver (important for Bun)
	sqldb, err := sql.Open(sqliteshim.ShimName, dataSourceName)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Create a Bun DB instance
	db := bun.NewDB(sqldb, sqlitedialect.New())

	// Add query hook for debugging (optional, remove in production)
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Ping the database to verify connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")

	// Create the todos table if it doesn't exist
	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*models.Todo)(nil)).IfNotExists().Exec(ctx) // Use models.Todo
	if err != nil {
		log.Printf("Error creating 'todos' table: %v", err)
		return nil, err
	}

	log.Println("'todos' table checked/created successfully.")

	return db, nil
}
