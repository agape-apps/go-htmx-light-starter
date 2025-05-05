# Directory Structure for Go/Chi, Templ, HTMX, Alpine.js Stack

```
project-root/
├── cmd/                    # Application entry points
│   └── server/             # Main server application
│       └── main.go         # Server initialization and configuration
├── docs/                   # Project documentation, API documentation for LLMs
├── frontend/               # Frontend source files
│   ├── js/                 # Alpine.js and custom JavaScript
│   └── styles/             # Bulma CSS source files
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   ├── models/             # Data domain models
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── repositories/       # Data access logic
│   └── services/           # Business logic
├── migrations/             # Database migrations
├── static/                 # Static assets
│   ├── css/                # Processed CSS output
│   ├── js/                 # Bundled JS output
│   └── images/             # Image assets
├── templates/              # Templ templates
│   ├── components/         # Reusable UI components
│   ├── layouts/            # Page layouts
│   │   └── base.templ      # Base layout with common structure (head, body, etc.)
│   └── pages/              # Page templates
│       └── index.templ     # Index/homepage template
├── .env                    # Environment variables
├── .gitignore              # Git ignore file
├── .air.toml               # Air config (with Templ step)
├── bunfig.toml             # Bun configuration
├── go.mod                  # Go modules
├── go.sum                  # Go modules checksum
├── Makefile                # Combined dev + build commands
├── package.json            # Bun-managed manifest
└── README.md               # Project documentation
```

## `internal` directory structure

Let's break down the `internal` directory structure for this Go project. This directory is standard practice in Go projects to hold code that isn't meant to be imported by other external projects, keeping your application's core logic private.

Here's an explanation of each subdirectory within `internal`, along with simple examples:

1.  **`internal/config/`**
    *   **Purpose:** Manages application configuration. This could involve loading settings from environment variables, configuration files (like `.env` or YAML files), or command-line flags.
    *   **Example:** You might have a `config.go` file here with functions to load database connection strings, API keys, or server port numbers.

    ```go
    // internal/config/config.go (Simplified Example)
    package config

    import "os"

    type AppConfig struct {
        DatabaseURL string
        ServerPort  string
    }

    func LoadConfig() (*AppConfig, error) {
        dbURL := os.Getenv("DATABASE_URL")
        port := os.Getenv("SERVER_PORT")
        if port == "" {
            port = "8080" // Default port
        }
        // Add error handling for missing required configs
        return &AppConfig{
            DatabaseURL: dbURL,
            ServerPort:  port,
        }, nil
    }
    ```

2.  **`internal/models/`**
    *   **Purpose:** Defines the core data structures (models or entities) of your application. These represent the fundamental concepts your application works with.
    *   **Example:** If you're building a blog, you might define `User`, `Post`, and `Comment` structs here.

    ```go
    // internal/models/post.go (Simplified Example)
    package models

    import "time"

    type Post struct {
        ID        int
        Title     string
        Content   string
        CreatedAt time.Time
        AuthorID  int
    }

    // You might also define interfaces here that repositories must implement
    type PostRepository interface {
        GetByID(id int) (*Post, error)
        Create(post *Post) error
        // ... other methods
    }
    ```

3.  **`internal/handlers/`**
    *   **Purpose:** Contains the HTTP request handlers (often called controllers in other frameworks). These functions receive incoming HTTP requests, parse data (like form values or URL parameters), call appropriate services to perform actions, and generate HTTP responses, often using Templ templates to render HTML.
    *   **Example:** A handler to display a specific blog post. It would get the post ID from the URL, call a service to fetch the post data, and then render a Templ template with that data.

    ```go
    // internal/handlers/post_handlers.go (Simplified Example)
    package handlers

    import (
        "net/http"
        "strconv"
        // Assume 'templates' package exists for rendering
        // Assume 'services' package exists
        "your_project/internal/services"
        "your_project/templates" // Adjust import path

        "github.com/go-chi/chi/v5"
    )

    type PostHandler struct {
        PostService services.PostService // Inject the service
    }

    func NewPostHandler(ps services.PostService) *PostHandler {
        return &PostHandler{PostService: ps}
    }

    func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "postID")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid Post ID", http.StatusBadRequest)
            return
        }

        post, err := h.PostService.GetPostByID(id) // Call the service
        if err != nil {
            // Handle errors appropriately (e.g., not found)
            http.Error(w, "Post not found", http.StatusNotFound)
            return
        }

        // Render the Templ template for the post page
        templates.PostPage(post).Render(r.Context(), w)
    }
    ```

4.  **`internal/middleware/`**
    *   **Purpose:** Holds HTTP middleware functions. Middleware intercepts requests before they reach the handler or after the handler generates a response. It's used for cross-cutting concerns like logging, authentication, authorization, CORS, request ID generation, etc.
    *   **Example:** A simple logging middleware that prints information about each incoming request.

    ```go
    // internal/middleware/logger.go (Simplified Example)
    package middleware

    import (
        "log"
        "net/http"
        "time"
    )

    func Logger(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            log.Printf("Started %s %s", r.Method, r.URL.Path)
            next.ServeHTTP(w, r) // Pass request to the next middleware or handler
            log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
        })
    }

    // Another example: Authentication check
    func Authenticate(next http.Handler) http.Handler {
         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check for session cookie or token
            // If not authenticated: http.Error(w, "Unauthorized", http.StatusUnauthorized); return
            // If authenticated:
            //   Maybe add user info to request context: ctx := context.WithValue(r.Context(), userKey, user)
            //   next.ServeHTTP(w, r.WithContext(ctx))
            next.ServeHTTP(w, r) // Placeholder
         })
    }
    ```

5.  **`internal/repositories/`**
    *   **Purpose:** Contains the logic for data persistence and retrieval. This layer abstracts the database interactions (or interactions with other data sources like external APIs). It typically implements interfaces defined in the `domain` layer.
    *   **Example:** A repository for managing `Post` data in a SQL database.

    ```go
    // internal/repositories/post_sql_repo.go (Simplified Example)
    package repositories

    import (
        "database/sql"
        "your_project/internal/domain" // Adjust import path
    )

    type PostSQLRepository struct {
        DB *sql.DB // Database connection pool
    }

    func NewPostSQLRepository(db *sql.DB) *PostSQLRepository {
        return &PostSQLRepository{DB: db}
    }

    // Implement the domain.PostRepository interface
    func (repo *PostSQLRepository) GetByID(id int) (*domain.Post, error) {
        post := &domain.Post{}
        query := "SELECT id, title, content, created_at, author_id FROM posts WHERE id = ?"
        err := repo.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.AuthorID)
        if err != nil {
            // Handle sql.ErrNoRows specifically
            return nil, err
        }
        return post, nil
    }

    func (repo *PostSQLRepository) Create(post *domain.Post) error {
        query := "INSERT INTO posts (title, content, author_id, created_at) VALUES (?, ?, ?, ?)"
        _, err := repo.DB.Exec(query, post.Title, post.Content, post.AuthorID, post.CreatedAt)
        // Potentially retrieve and set the generated ID
        return err
    }
    ```

6.  **`internal/services/`**
    *   **Purpose:** Contains the core business logic of the application. Services orchestrate operations, often coordinating calls to one or more repositories and applying business rules. Handlers call services, keeping the handlers thin and focused on HTTP concerns.
    *   **Example:** A `PostService` that handles creating a post, perhaps performing validation or notifying other systems, before saving it via the repository.

    ```go
    // internal/services/post_service.go (Simplified Example)
    package services

    import (
        "errors"
        "time"
        "your_project/internal/domain" // Adjust import path
    )

    // Define an interface for the service (good for testing and dependency injection)
    type PostService interface {
        GetPostByID(id int) (*domain.Post, error)
        CreatePost(title, content string, authorID int) (*domain.Post, error)
    }

    // Implementation of the service
    type postServiceImpl struct {
        repo domain.PostRepository // Inject the repository interface
    }

    func NewPostService(repo domain.PostRepository) PostService {
        return &postServiceImpl{repo: repo}
    }

    func (s *postServiceImpl) GetPostByID(id int) (*domain.Post, error) {
        // Maybe add caching logic here later
        return s.repo.GetByID(id)
    }

    func (s *postServiceImpl) CreatePost(title, content string, authorID int) (*domain.Post, error) {
        if title == "" || content == "" {
            return nil, errors.New("title and content cannot be empty")
        }
        // Add more complex business rules/validation if needed

        post := &domain.Post{
            Title:     title,
            Content:   content,
            AuthorID:  authorID,
            CreatedAt: time.Now(),
        }

        err := s.repo.Create(post)
        if err != nil {
            return nil, err // Handle repository errors
        }

        // Maybe trigger notifications or other side effects here

        return post, nil // Return the created post (potentially with ID set by repo)
    }
    ```

This structure promotes separation of concerns, making the application easier to understand, test, maintain, and scale. Let me know if you'd like a deeper dive into any specific part, Sir Christian!