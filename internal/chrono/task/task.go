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

type TaskRepository interface {
	// CreateTask creates a new task in the database with status "pending".
	CreateTask(task *Task) error
	// FindTaskById returns a task with the given id.
	FindTaskById(id string) (*Task, error)
	// FindTaskByPartialId returns a task if there is only one match with the given partial id.
	FindTaskByPartialId(partialId string, filter Filter) (*Task, error)
	// FindPendingTaskByName returns a task with the given name and status "pending".
	FindPendingTaskByName(name string, filter Filter) (*Task, error)
	// FindTasksByStatus returns all tasks filtering by statuses in the database.
	ListTasksByStatus(statuses ...Status) ([]Task, error)
	// RemoveTaskById removes a task by id.
	RemoveTaskById(id string) error
}

type Service interface {
	// CreateTask creates a new task if one with the same name and status "pending" does not exists already.
	CreateTask(name, description string) error
	// ListTasks returns tasks filtering by statuses.
	ListTasksByStatus(statuses ...Status) ([]Task, error)
	// RemoveTaskByPartialId removes a task by partial id.
	RemoveTaskByPartialId(partialId string) error
	// StartTask starts a task by id or name.
	StartTask(idOrName string) (string, error)
}
