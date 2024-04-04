package task

import "time"

type Status string

const (
	PENDING Status = "pending"
	IN_PROGRESS Status = "in_progress"
	PAUSED Status = "paused"
	DONE Status = "done"
	CANCELED Status = "canceled"
)

type Task struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	// CreateTask creates a new task in the database with status "pending".
	CreateTask(task *Task) error
	// FindTaskById returns a task with the given id.
	FindTaskById(id string) (*Task, error)
	// FindPendingTaskByName returns a task with the given name and status "pending".
	FindPendingTaskByName(name string) (*Task, error)
	// FindTasksByStatus returns all tasks filtering by statuses in the database.
	ListTasksByStatus(statuses ...Status) ([]Task, error)
}

type Service interface {
	// CreateTask creates a new task if one with the same name and status "pending" does not exists already.
	CreateTask(name, description string) error
	// ListTasks returns tasks filtering by statuses.
	ListTasksByStatus(statuses ...Status) ([]Task, error)
}