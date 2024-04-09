package task

import (
	"testing"
	"time"
)

type MockRepository interface {
	CreateTask(task *Task) error
	FindTaskById(id string) (*Task, error)
	FindPendingTaskByName(name string) (*Task, error)
	ListTasksByStatus(statuses ...Status) ([]Task, error)
}

type mockTaskRepository struct{
	taskCreated *Task
	MockRepository
}

func (r mockTaskRepository) CreateTask(task *Task) error {
	*r.taskCreated = *task
	return nil
}

func (r mockTaskRepository) FindPendingTaskByName(name string) (*Task, error) {
	return &Task{}, nil
}

func (r mockTaskRepository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	return []Task{
		{
			ID:          "1",
			Name:        "Task 1",
			Description: "Description 1",
			Status:      PENDING,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Description: "Description 2",
			Status:      PENDING,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}, nil
}

type mockTaskFindDataRepository struct {
	MockRepository
}

func (r mockTaskFindDataRepository) FindPendingTaskByName(name string) (*Task, error) {
	return &Task{
		ID:          "1",
		Name:        name,
		Description: "Description 1",
		Status:      PENDING,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}, nil
}

func TestCreateTaskShouldCreateNewTask(t *testing.T) {
	r := mockTaskRepository{}
	r.taskCreated = &Task{}
	s := NewService(r)

	name := "Task 1"
	description := "Description 1"

	err := s.CreateTask(name, description)

	if err != nil {
		t.Errorf("error was not expected while creating task: %s", err)
	}

	taskCreated := r.taskCreated

	if taskCreated.Status != PENDING {
		t.Errorf("expected task status to be pending, got %s", taskCreated.Status)
	}
}

func TestCreateTaskShouldNotCreateTaskWhenTaskExists(t *testing.T) {
	r := mockTaskFindDataRepository{}
	s := NewService(r)

	name := "Task 1"
	description := "Description 1"

	err := s.CreateTask(name, description)

	if err == nil {
		t.Errorf("expected error while creating task")
	}
}


func TestListTasksByStatusShouldReturnTasks(t *testing.T) {
	r := mockTaskRepository{}
	s := NewService(r)

	statuses := []Status{PENDING}

	tasks, err := s.ListTasksByStatus(statuses...)

	if err != nil {
		t.Errorf("error was not expected while listing tasks: %s", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}
