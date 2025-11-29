package tasks

import (
	"time"
)

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"index"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"` // LOW, MEDIUM, HIGH
	Status      string     `json:"status"`   // TODO, IN_PROGRESS, DONE
	DueDate     *time.Time `json:"due_date"`
	Tags        []Tag      `json:"tags" gorm:"many2many:task_tags;"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
