package counterapp

import (
	"log"
	"net/http"
	"strconv"

	counterapptemplates "go-htmx-light-starter/counterapp/templates"
	"go-htmx-light-starter/templates/layouts"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

// Handlers holds the application state and dependencies for handlers.
type Handlers struct {
	DB         *bun.DB // Database connection pool
	counter    int
	maxCounter int
}

// NewHandlers creates a new Handlers instance with dependencies.
func NewHandlers(db *bun.DB) *Handlers {
	return &Handlers{
		DB:         db, // Store the database connection
		counter:    0,
		maxCounter: 100, // Default max counter value
	}
}

// Home handles the request for the home page.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	// Create the index page with counter and max value
	indexPage := counterapptemplates.Index(h.counter, h.maxCounter)
	baseComponent := layouts.Base("My Web App", indexPage)
	templ.Handler(baseComponent).ServeHTTP(w, r)
}

// Increment handles the request to increment the counter.
func (h *Handlers) Increment(w http.ResponseWriter, r *http.Request) {
	// Only increment if counter is less than max
	if h.counter < h.maxCounter {
		h.counter++
	}
	log.Printf("Increment counter: %d (max: %d)", h.counter, h.maxCounter)
	// Return the updated counter component using htmx-go and Templ
	counterComponent := counterapptemplates.Counter(h.counter, h.maxCounter)
	htmx.NewResponse().RenderTempl(r.Context(), w, counterComponent)
}

// Decrement handles the request to decrement the counter.
func (h *Handlers) Decrement(w http.ResponseWriter, r *http.Request) {
	// Only decrement if counter is greater than negative max
	if h.counter > -h.maxCounter {
		h.counter--
	}
	log.Printf("Decrement counter: %d (max: %d)", h.counter, h.maxCounter)
	// Return the updated counter component using htmx-go and Templ
	counterComponent := counterapptemplates.Counter(h.counter, h.maxCounter)
	htmx.NewResponse().RenderTempl(r.Context(), w, counterComponent)
}

// SetMax handles the request to set the maximum counter value.
func (h *Handlers) SetMax(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get max value from form
	maxStr := r.FormValue("max")
	newMax, err := strconv.Atoi(maxStr)
	// Basic validation: ensure max is between 1 and 100 (adjust as needed)
	if err != nil || newMax < 1 || newMax > 1000 { // Increased upper limit for flexibility
		http.Error(w, "Invalid max value (must be between 1 and 1000)", http.StatusBadRequest)
		return
	}

	// Update max counter
	h.maxCounter = newMax
	log.Printf("Set max counter: %d", h.maxCounter)

	// Return the entire counter component using htmx-go and Templ
	counterComponent := counterapptemplates.Counter(h.counter, h.maxCounter)
	htmx.NewResponse().RenderTempl(r.Context(), w, counterComponent)
}

// Routes sets up the Chi router with all application routes.
func (h *Handlers) Routes() chi.Router {
	r := chi.NewRouter()

	// Page Routes
	r.Get("/", h.Home)

	// API/HTMX Routes
	r.Post("/increment", h.Increment)
	r.Post("/decrement", h.Decrement)
	r.Post("/set-max", h.SetMax)

	return r
}
