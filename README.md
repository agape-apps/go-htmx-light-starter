# Go HTMX Light Starter

A light Go/Chi/Templ web application starter project using HTMX/AlpineJS and Bulma CSS.

This is my personal base project for experimenting with Go and HTMX. If you are looking for something similar, but regularly maintained and updated, check out https://github.com/gowebly

## Tech Stack

Refer to the Tech Stack documentation links of each technology for best practices and detailed examples of usage patterns.

- GoLang in the backend to handle server requests. Docs: https://go.dev/doc/
- Chi as Go web framework/router. Docs: https://go-chi.io/#/pages/routing
- Templ templating to generate HTMX. Docs: https://templ.guide/
- HTMX. Docs: https://htmx.org/reference/
- AlpineJS for reactivity with HTMX. Docs: https://alpinejs.dev/
- Bulma CSS for UI styling. Docs: https://bulma.io/documentation/customize/with-sass/


## Hot reloading
- Air for Go and Templ files
- esbuild for JavaScript files
- sass for SCSS files

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.24 or later)
- [Bun](https://bun.sh) (for JavaScript module installation)
- [Air](https://github.com/air-verse/air) (for hot reloading)

```bash
curl -fsSL https://bun.sh/install | bash
go install github.com/air-verse/air@latest
```

## Installation

- Install Go dependencies

```bash
go mod download
go mod tidy
```

- Install JavaScript dependencies
```bash
bun install
```

## Development

To start the development server with hot reloading:

```bash
bun run watch
```

This will:
- Watch for changes in Go files and **restart the server** (using Air)
- Watch for changes in Templ files and regenerate them
- Watch for changes in JavaScript files and rebuild
- Watch for changes in CSS files and rebuild

## Building

To build the application:

```bash
bun run build
```

This will:
- Generate Templ templates
- Build JavaScript and CSS files

## Running

To run the application without watching for changes:

```bash
# Just run the server
bun run goserve

# Build first, then run
bun run gostart
```

The application will be available at http://localhost:8080

## Main Project Structure

- `cmd/server/main.go` - Main Go server
- `templates/` - Templ templates
- `frontend/src/` - Frontend JavaScript and CSS
- `static/` - Compiled static assets
- `.env`- Environment variables

See also: docs/DirectoryStructure.md

## Production Build & Deployment

To create an optimized build suitable for production deployment:

```bash
bun run build:production
```

This command performs the following steps:
1.  Generates Go code from `.templ` files (`build:templ`).
2.  Builds and minifies frontend assets (JS/CSS) into the `static/` directory (`build:frontend`).
3.  Compiles the Go application into a single, optimized binary located at `./build/server` (`build:server`).

**Deployment:**

To deploy the application, copy the following to your production server:
- The compiled binary: `./build/server`
- The entire static assets directory: `./static/`
- Any necessary configuration files (e.g., `.env` if used).

**Running in Production:**

On the production server, navigate to the directory where you copied the files and simply execute the binary:

```bash
./build/server
```

Ensure the binary has execute permissions (`chmod +x ./build/server`). The server will start, typically listening on the configured port (e.g., 8080). It's recommended to run the binary using a process manager (like `systemd`, `supervisor`, or `docker`) for resilience and logging in a real production environment.
