package task

import (
	"fmt"

	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
)

type service struct {
	TaskRepository     TaskRepository
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

func (s *service) CreateTask(name string) error {
	task, err := s.TaskRepository.FindTaskByNameAndStatus(name, Filter{
		Statuses: []Status{TODO, IN_PROGRESS},
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
	task.Status = TODO

	err = s.TaskRepository.CreateTask(&task)
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

func (s *service) StartTask(idOrName string) error {
	curr, err := s.TaskRepository.ListTasksByStatus(IN_PROGRESS)
	if err != nil {
		return fmt.Errorf("error consulting the database: %v", err)
	}

	if len(curr) > 0 {
		if len(curr) > 1 {
			return fmt.Errorf("More than 1 task with status \"in_progress\". Please delete or update duplicated.")
		} else {
			curr[0].Status = TODO
			err = s.TaskRepository.UpdateTask(&curr[0])
			if err != nil {
				return fmt.Errorf("error updating task: %v", err)
			}

			p := s.ProgressRepository.GetLastProgressByTaskID(curr[0].ID)
			p.Status = string(DONE)
			err = s.ProgressRepository.UpdateProgress(&p)
			if err != nil {
				return fmt.Errorf("error updating progress: %v", err)
			}
		}
	}

	task, err := s.TaskRepository.FindByIdOrCreate(idOrName, Filter{
		Statuses: []Status{TODO, IN_PROGRESS},
	})
	if err != nil {
		return fmt.Errorf("error consulting the database: %v", err)
	}

	progress := progress.Progress{}
	progress.TaskID = task.ID
	progress.Status = string(IN_PROGRESS)

	err = s.ProgressRepository.AddProgress(&progress)
	if err != nil {
		return fmt.Errorf("error creating progress: %v", err)
	}

	if task.Status != IN_PROGRESS {
		task.Status = IN_PROGRESS
		err = s.TaskRepository.UpdateTask(&task)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}
	}

	return nil
}

func (s *service) FinishTask(id string) error {
	task, err := s.TaskRepository.FindTaskByPartialId(
		id,
		Filter{Statuses: []Status{TODO, IN_PROGRESS}},
	)

	if err != nil {
		return fmt.Errorf("error consulting the database: %v", err)
	}

	p := progress.Progress{}
	p.TaskID = task.ID
	p.Status = string(DONE)

	err = s.ProgressRepository.AddProgress(&p)
	if err != nil {
		return fmt.Errorf("error creating progress: %v", err)
	}

	task.Status = DONE
	err = s.TaskRepository.UpdateTask(&task)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return nil
}
