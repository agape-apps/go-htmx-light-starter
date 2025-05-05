# PRD: Minimal TODO Application (TodoNow)

**Version:** 1.1
**Date:** 2025-05-04

## Files Created:
data/todonow.db SQLite               database file
internal/database/db.go              initializes the database and table
internal/handlers/todo_handlers.go 
internal/models/todo.go              table definition
internal/repositories/interfaces.go
internal/repositories/todo_repo.go

### UI
templates/components/todo_list.templ
templates/components/todo_item.templ
templates/pages/todonow.templ

## 1. Introduction

This document outlines the requirements for a minimal TODO application feature, named "TodoNow", to be integrated into the Go/HTMX Light Starter project. The primary goal is to demonstrate the usage of the Bun ORM with a SQLite database within the existing project architecture.

## 2. Goals

*   Implement a basic Create, Read, Update (Toggle Completion), Delete (CRUD) functionality for TODO items.
*   Demonstrate integration of Bun ORM for database interactions.
*   Utilize SQLite as the database backend.
*   Adhere to the existing project structure and conventions.
*   Keep the Go code specific to the TODO feature concise (approx. 100-200 LOC). Actuallly it became about 430 LOC
*   Provide a seamless user experience using HTMX for interactions.

## 3. Functional Requirements

### 3.1. User Interface (UI)

*   **Route:** The TODO application will be accessible at the `/todonow` route.
*   **Display:**
    *   The page will display a list of existing TODO items.
    *   Each item will show its task description and completion status.
    *   A form will be present to add new TODO items.
*   **Interactions (HTMX-driven):**
    *   **Add Todo:** Users can type a task into an input field and click an "Add" button. The new task should appear in the list without a full page reload.
    *   **Toggle Completion:** Users can click a checkbox or button associated with a TODO item to mark it as complete or incomplete. The item's visual state should update instantly.
    *   **Delete Todo:** Users can click a "Delete" button associated with a TODO item. The item should be removed from the list without a full page reload.

### 3.2. Backend & Data

*   **Data Model:** A `Todo` model will be defined with the following fields:
    *   `ID` (integer, primary key, auto-incrementing)
    *   `Task` (string, non-empty)
    *   `Completed` (boolean, default: false)
*   **Database:**
    *   SQLite will be used as the database.
    *   A database file (e.g., `data/todonow.db`) will store the TODO items. Need to create `data` directory at project root.
    *   Bun ORM will be used for all database operations (setup, migrations, CRUD).
*   **API Endpoints (Internal/HTMX):** Specific endpoints will be created within the Go application to handle the HTMX requests for adding, toggling, and deleting TODO items. These will reside within a dedicated handler: internal/handlers/todohandlers.go

## 4. Non-Functional Requirements

*   **Technology Stack:** Go, Chi, Templ, Bun ORM, SQLite, HTMX, Alpine.js (available, but should not be needed for this simple app), TailwindCSS/DaisyUI.
*   **Code Structure:** Follow the established `internal` directory structure (models, repositories, handlers, potentially services if logic becomes complex, though aiming for simplicity). Templates will be placed in `templates/pages` and potentially `templates/components`.
*   **Performance:** Interactions should feel responsive due to HTMX partial updates.
*   **Maintainability:** Code should be clean, well-commented where necessary, and follow Go best practices.

## 5. Out of Scope (DO NOT IMPLEMENT)

*   User authentication or multiple user support.
*   Due dates, priorities, or categories for TODO items.
*   Advanced features like searching or filtering.
*   Offline support.