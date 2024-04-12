package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
)

type service struct {
	TaskRepository TaskRepository
	ProgressRepository progress.ProgressRepository
}

func NewService(t TaskRepository, p progress.ProgressRepository) Service {
	return &service{
		TaskRepository:     t,
		ProgressRepository: p,
	}
}

func (s *service) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	tasks, err := s.TaskRepository.ListTasksByStatus(statuses...)
	if err != nil {
		return nil, fmt.Errorf("error consulting the database: %v", err)
	}

	return tasks, nil
}

func (s *service) CreateTask(name, description string) error {
	task, err := s.TaskRepository.FindPendingTaskByName(name, Filter{
		Statuses: []Status{PENDING, IN_PROGRESS, PAUSED},
	})
	if err != nil {
		if err.Error() != "not found" {
			return fmt.Errorf("error consulting the database: %v", err)
		}
	}

	if task.ID != "" {
		return fmt.Errorf("pending task \"%s\" already exists", name)
	}

	task.ID = uuid.New().String()
	task.Name = name
	task.Description = description
	task.Status = PENDING
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err = s.TaskRepository.CreateTask(task)
	if err != nil {
		return fmt.Errorf("error creating task: \"%v\"", err)
	}

	return nil
}

// RemoveTaskByPartialId implements Service.
func (s *service) RemoveTaskByPartialId(partialId string) error {
	task, err := s.TaskRepository.FindTaskByPartialId(partialId, Filter{
		Statuses: []Status{},
	})

	if err != nil {
		return fmt.Errorf("error consulting the database: %v", err)
	}

	err = s.TaskRepository.RemoveTaskById(task.ID)
	if err != nil {
		return fmt.Errorf("error removing task: %v", err)
	}

	return nil
}

func (s *service) StartTask(idOrName string) (string, error) {
	var task Task

	task, err := s.TaskRepository.FindTaskByPartialId(idOrName, Filter{
		Statuses: []Status{PENDING, IN_PROGRESS, PAUSED},
	})

	if err == nil {
		// err = s.TaskRepository.UpdateTaskStatus(task.ID, IN_PROGRESS)
		// if err != nil {
		// 	return "", fmt.Errorf("error updating task status: %v", err)
		// }
		// return "Task started", nil
	}

	if err.Error() == "not found" {
		task.ID = uuid.New().String()
		task.Name = idOrName
		task.Status = IN_PROGRESS
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()

		err = s.TaskRepository.CreateTask(&task)
		if err != nil {
			return "", fmt.Errorf("error creating task: \"%v\"", err)
		}
		task, err = s.TaskRepository.CreateTask()
	}

	if err != nil {
		if err.Error() == "not found" {
			task.ID = uuid.New().String()
			task.Name = idOrName
			task.Status = IN_PROGRESS
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()

			err = s.TaskRepository.CreateTask(&task)
			if err != nil {
				return "", fmt.Errorf("error creating task: \"%v\"", err)
			}
			task, err = s.TaskRepository.CreateTask()
		} else {
			return "", fmt.Errorf("error consulting the database: %v", err)
		}
	}

	err = s.ProgressRepository.AddProgress(progress.Progress{
		ID:			  uuid.New().String(),
		TaskID:       task.ID,
		StatusInit:   string(task.Status),
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	})

	if err != nil {
		return "", fmt.Errorf("error creating progress: %v", err)
	}

	return "Task started", nil
}
