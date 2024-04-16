package task

import (
	"testing"
	"time"

	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
)

type MockRepository interface {
	CreateTask(task *Task) error
	FindTaskById(id string) (Task, error)
	FindTaskByPartialId(partialId string, filter Filter) (Task, error)
	FindPendingTaskByName(name string, filter Filter) (Task, error)
	ListTasksByStatus(statuses ...Status) ([]Task, error)
	RemoveTaskById(id string) error
	FindByIdOrCreate(idOrName string, filter Filter) (Task, error)
	UpdateTask(task *Task) error
}

type mockProgressRepository struct {
	MockRepository
}

func (r mockProgressRepository) AddProgress(progress *progress.Progress) error {
	return nil
}

type mockTaskRepository struct {
	taskCreated *Task
	MockRepository
}

func (r mockTaskRepository) FindTaskById(id string) (Task, error) {
	return Task{}, nil
}

func (r mockTaskRepository) FindByIdOrCreate(idOrName string, filter Filter) (Task, error) {
	return Task{}, nil
}

func (r mockTaskRepository) UpdateTask(task *Task) error {
	return nil
}

func (r mockTaskRepository) FindTaskByPartialId(partialId string, filter Filter) (Task, error) {
	return Task{}, nil
}

func (r mockTaskRepository) RemoveTaskById(id string) error {
	return nil
}

func (r mockTaskRepository) CreateTask(task *Task) error {
	*r.taskCreated = *task
	return nil
}

func (r mockTaskRepository) FindPendingTaskByName(name string, filter Filter) (Task, error) {
	return Task{}, nil
}

func (r mockTaskRepository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	return []Task{
		{
			ID:        "1",
			Name:      "Task 1",
			Status:    PENDING,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		{
			ID:        "2",
			Name:      "Task 2",
			Status:    PENDING,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}, nil
}

type mockTaskFindDataRepository struct {
	MockRepository
}

func (r mockTaskFindDataRepository) CreateTask(task *Task) error {
	return nil
}

func (r mockTaskFindDataRepository) FindTaskById(id string) (Task, error) {
	return Task{}, nil
}

func (r mockTaskFindDataRepository) FindTaskByPartialId(partialId string, filter Filter) (Task, error) {
	return Task{}, nil
}

func (r mockTaskFindDataRepository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	return []Task{}, nil
}

func (r mockTaskFindDataRepository) RemoveTaskById(id string) error {
	return nil
}

func (r mockTaskFindDataRepository) FindPendingTaskByName(name string, filter Filter) (Task, error) {
	return Task{
		ID:        "1",
		Name:      name,
		Status:    PENDING,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, nil
}

func TestCreateTaskShouldCreateNewTaskWithStatusPending(t *testing.T) {
	r := mockTaskRepository{}
	p := mockProgressRepository{}
	r.taskCreated = &Task{}
	s := NewService(r, p)

	name := "Task 1"

	err := s.CreateTask(name)

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
	p := mockProgressRepository{}
	s := NewService(r, p)

	name := "Task 1"

	err := s.CreateTask(name)

	if err == nil {
		t.Errorf("expected error while creating task")
	}
}

func TestListTasksByStatusShouldReturnTasks(t *testing.T) {
	r := mockTaskRepository{}
	p := mockProgressRepository{}
	s := NewService(r, p)

	statuses := []Status{PENDING}

	tasks, err := s.ListTasksByStatus(statuses...)

	if err != nil {
		t.Errorf("error was not expected while listing tasks: %s", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}
