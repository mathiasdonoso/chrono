package task

import (
	"testing"
	"time"
)

type MockRepository interface {
	CreateTask(task *Task) error
	FindTaskById(id string) (*Task, error)
	FindPendingTaskByName(name string) (*Task, error)
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
