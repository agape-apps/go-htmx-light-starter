package todonow

import (
	"context"

	"go-htmx-light-starter/todonow/models" // Import the models package
)

// TodoRepository defines the methods for interacting with Todo data storage.
type TodoRepository interface {
	GetAll(ctx context.Context) ([]models.Todo, error)           // Use models.Todo
	GetByID(ctx context.Context, id int64) (*models.Todo, error) // Use models.Todo
	Add(ctx context.Context, todo *models.Todo) error            // Use models.Todo
	Toggle(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
}
