package todonow

import (
	"context"
	"database/sql"
	"log"

	"go-htmx-light-starter/todonow/models" // Import models

	"github.com/uptrace/bun"
)

// todoRepo implements the TodoRepository interface using Bun.
type todoRepo struct {
	db *bun.DB
}

// NewTodoRepository creates a new repository instance.
func NewTodoRepository(db *bun.DB) TodoRepository { // TodoRepository is in this package
	return &todoRepo{db: db}
}

// GetAll retrieves all Todo items from the database.
func (r *todoRepo) GetAll(ctx context.Context) ([]models.Todo, error) { // Use models.Todo
	var todos []models.Todo // Use models.Todo
	err := r.db.NewSelect().Model(&todos).Order("id ASC").Scan(ctx)
	if err != nil {
		log.Printf("Error getting all todos: %v", err)
		return nil, err
	}
	return todos, nil
}

// GetByID retrieves a single Todo item by its ID.
func (r *todoRepo) GetByID(ctx context.Context, id int64) (*models.Todo, error) { // Use models.Todo
	todo := new(models.Todo) // Use models.Todo
	err := r.db.NewSelect().Model(todo).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Todo with ID %d not found", id)
			return nil, err // Or return a specific "not found" error
		}
		log.Printf("Error getting todo by ID %d: %v", id, err)
		return nil, err
	}
	return todo, nil
}

// Add inserts a new Todo item into the database.
// It updates the passed Todo struct with the newly generated ID.
func (r *todoRepo) Add(ctx context.Context, todo *models.Todo) error { // Use models.Todo
	_, err := r.db.NewInsert().Model(todo).Exec(ctx)
	if err != nil {
		log.Printf("Error adding todo: %v", err)
		return err
	}
	log.Printf("Added todo ID: %d, Task: %s", todo.ID, todo.Task)
	return nil
}

// Toggle updates the completion status of a Todo item.
func (r *todoRepo) Toggle(ctx context.Context, id int64) error {
	// Fetch the current status first to toggle it
	currentTodo, err := r.GetByID(ctx, id)
	if err != nil {
		return err // Error already logged in GetByID
	}

	newStatus := !currentTodo.Completed

	_, err = r.db.NewUpdate().
		Model((*models.Todo)(nil)). // Use models.Todo
		Set("completed = ?", newStatus).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		log.Printf("Error toggling todo ID %d: %v", id, err)
		return err
	}
	log.Printf("Toggled todo ID %d to completed=%t", id, newStatus)
	return nil
}

// Delete removes a Todo item from the database by its ID.
func (r *todoRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.NewDelete().Model((*models.Todo)(nil)).Where("id = ?", id).Exec(ctx) // Use models.Todo
	if err != nil {
		log.Printf("Error deleting todo ID %d: %v", id, err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Attempted to delete non-existent todo ID %d", id)
		// Depending on requirements, you might return sql.ErrNoRows or a custom error
	} else {
		log.Printf("Deleted todo ID %d", id)
	}

	return nil
}
