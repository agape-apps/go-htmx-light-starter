package todonow

import (
	"log"
	"net/http"
	"strconv"

	"go-htmx-light-starter/templates/layouts"
	"go-htmx-light-starter/todonow/models" // Import models
	todonowtemplates "go-htmx-light-starter/todonow/templates"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
)

// TodoHandler holds dependencies for todo-related handlers.
type TodoHandler struct {
	Repo TodoRepository // Type is in this package
}

// NewTodoHandler creates a new TodoHandler instance.
func NewTodoHandler(repo TodoRepository) *TodoHandler { // Type is in this package
	return &TodoHandler{Repo: repo}
}

// RegisterTodoRoutes sets up the routes for the TODO feature.
func (h *TodoHandler) RegisterTodoRoutes(r chi.Router) {
	r.Get("/", h.ShowPage)
	r.Post("/add", h.Add)
	r.Put("/toggle/{id}", h.Toggle)    // Use PUT for idempotent update
	r.Delete("/delete/{id}", h.Delete) // Use DELETE for deletion
}

// ShowPage displays the main TODO list page.
func (h *TodoHandler) ShowPage(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting all todos for page: %v", err)
		http.Error(w, "Failed to load TODO items", http.StatusInternalServerError)
		return
	}

	// Create the todo page component, passing the todos
	todoPage := todonowtemplates.TodoNowPage(todos)
	// Wrap it in the base layout
	baseComponent := layouts.Base("TodoNow App", todoPage)

	// Render the full page
	// Use htmx response helper in case this handler is called via HTMX request later
	htmx.NewResponse().RenderTempl(r.Context(), w, baseComponent)
}

// Add handles adding a new TODO item.
func (h *TodoHandler) Add(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")
	if task == "" {
		// Optionally return an error snippet via HTMX
		log.Println("Add Todo: Task cannot be empty")
		http.Error(w, "Task cannot be empty", http.StatusBadRequest) // Simple error for now
		return
	}

	newTodo := &models.Todo{ // Use models.Todo
		Task:      task,
		Completed: false,
	}

	err := h.Repo.Add(r.Context(), newTodo) // Add returns the ID in newTodo
	if err != nil {
		log.Printf("Error adding todo: %v", err)
		http.Error(w, "Failed to add TODO item", http.StatusInternalServerError)
		return
	}

	// Render just the new item component for HTMX
	todoItemComponent := todonowtemplates.TodoItem(*newTodo)
	htmx.NewResponse().RenderTempl(r.Context(), w, todoItemComponent)
}

// Toggle handles marking a TODO item as complete or incomplete.
func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID for toggle: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.Toggle(r.Context(), id)
	if err != nil {
		// Handle potential "not found" error from repo if necessary
		log.Printf("Error toggling todo ID %d: %v", id, err)
		http.Error(w, "Failed to toggle TODO item", http.StatusInternalServerError)
		return
	}

	// Fetch the updated item to render it
	updatedTodo, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching updated todo ID %d after toggle: %v", id, err)
		// If toggle succeeded but fetch failed, what to return?
		// For simplicity, return error, but could return old state or success message
		http.Error(w, "Failed to fetch updated TODO status", http.StatusInternalServerError)
		return
	}

	// Render the updated item component for HTMX swap
	todoItemComponent := todonowtemplates.TodoItem(*updatedTodo)
	htmx.NewResponse().RenderTempl(r.Context(), w, todoItemComponent)
}

// Delete handles deleting a TODO item.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID for delete: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting todo ID %d: %v", id, err)
		http.Error(w, "Failed to delete TODO item", http.StatusInternalServerError)
		return
	}

	// Return empty 200 OK for HTMX to remove the element
	w.WriteHeader(http.StatusOK)
}
