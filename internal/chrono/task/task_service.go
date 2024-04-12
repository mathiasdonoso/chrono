package task

import (
	"fmt"

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

	task.Name = name

	err = s.TaskRepository.CreateTask(&task); if err != nil {
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

	err = s.TaskRepository.RemoveTaskById(task.ID); if err != nil {
		return fmt.Errorf("error removing task: %v", err)
	}

	return nil
}

func (s *service) StartTask(idOrName string) error {
	task, err := s.TaskRepository.FindByIdOrCreate(idOrName, Filter{
		Statuses: []Status{PENDING, IN_PROGRESS, PAUSED},
	})
	if err != nil {
		fmt.Errorf("error consulting the database: %v", err)
	}

	// progress := progress.Progress{}
	// progress.TaskID = task.ID
	// progress.StatusInit = string(task.Status)

	// err = s.ProgressRepository.AddProgress(&progress)
	// fmt.Println("2", progress)
	// if err != nil {
	// 	fmt.Errorf("error creating progress: %v", err)
	// }

	fmt.Println("status is", task.Status)

	// FIX: This is not updating the task status
	if task.Status != IN_PROGRESS {
		task.Status = IN_PROGRESS
		err = s.TaskRepository.UpdateTask(&task); if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}
	}

	return nil
}
