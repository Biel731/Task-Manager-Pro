package tasks

import (
	"time"
)

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"required"`
	Status      string     `json:"status" binding:"required"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority"`
	Status      *string    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	Tags        *[]string  `json:"tags"`
}

type TaskFilter struct {
	Status   string
	Priority string
	Tags     string
	Query    string
}
