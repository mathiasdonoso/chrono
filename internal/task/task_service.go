package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type service struct {
	Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateTask(name, description string) (string, error) {
	task, err := s.Repository.FindTaskByName(name)
	if err != nil {
		if err.Error() != "task with name "+name+" not found" {
			return "", err
		}
	}

	if task.ID != "" {
		return "", fmt.Errorf("task with name %s already exists", name)
	}

	// create task
	task.ID = uuid.New().String()
	task.Name = name
	task.Description = description
	task.Status = PENDING
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err = s.Repository.CreateTask(task)
	if err != nil {
		return "", fmt.Errorf("error creating task: %v", err)
	}

	return "Task created!", nil
}
