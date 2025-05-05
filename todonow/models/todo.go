package models

import "github.com/uptrace/bun"

// Todo represents a single task item in the TODO list.
type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"` // Define table name and alias

	ID        int64  `bun:"id,pk,autoincrement"`     // Primary key, auto-increment
	Task      string `bun:"task,notnull"`            // Task description, cannot be null
	Completed bool   `bun:"completed,default:false"` // Completion status, defaults to false
}
