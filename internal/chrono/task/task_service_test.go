package task

import (
	"testing"
	"time"

	"github.com/mathiasdonoso/chrono/internal/chrono/progress"
)

var createTaskMock func(task *Task) error
var findTaskByIdMock func(id string) (Task, error)
var findTaskByPartialIdMock func(partialId string, filter Filter) (Task, error)
var findTaskByNameAndStatusMock func(name string, filter Filter) (Task, error)
var listTasksByStatusMock func(statuses ...Status) ([]Task, error)
var removeTaskByIdMock func(id string) error
var findByIdOrCreateMock func(idOrName string, filter Filter) (Task, error)
var updateTaskMock func(task *Task) error

type MockTaskRepository interface {
	CreateTask(task *Task) error
	FindTaskById(id string) (Task, error)
	FindTaskByPartialId(partialId string, filter Filter) (Task, error)
	FindTaskByNameAndStatus(name string, filter Filter) (Task, error)
	ListTasksByStatus(statuses ...Status) ([]Task, error)
	RemoveTaskById(id string) error
	FindByIdOrCreate(idOrName string, filter Filter) (Task, error)
	UpdateTask(task *Task) error
}

func (r mockTaskRepository) CreateTask(task *Task) error {
	return createTaskMock(task)
}
func (r mockTaskRepository) FindTaskById(id string) (Task, error) {
	return findTaskByIdMock(id)
}
func (r mockTaskRepository) FindTaskByPartialId(partialId string, filter Filter) (Task, error) {
	return findTaskByPartialIdMock(partialId, filter)
}
func (r mockTaskRepository) FindTaskByNameAndStatus(name string, filter Filter) (Task, error) {
	return findTaskByNameAndStatusMock(name, filter)
}
func (r mockTaskRepository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	return listTasksByStatusMock(statuses...)
}
func (r mockTaskRepository) RemoveTaskById(id string) error {
	return removeTaskByIdMock(id)
}
func (r mockTaskRepository) FindByIdOrCreate(idOrName string, filter Filter) (Task, error) {
	return findByIdOrCreateMock(idOrName, filter)
}
func (r mockTaskRepository) UpdateTask(task *Task) error {
	return updateTaskMock(task)
}

var addProgressMock func(progress *progress.Progress) error
var getLastProgressByTaskIDMock func(taskID string) progress.Progress
var updateProgressMock func(progress *progress.Progress) error

type MockProgressRepository interface {
	AddProgress(progress *progress.Progress) error
	GetLastProgressByTaskID(taskID string) progress.Progress
	UpdateProgress(progress *progress.Progress) error
}

func (r mockProgressRepository) AddProgress(progress *progress.Progress) error {
	return addProgressMock(progress)
}
func (r mockProgressRepository) GetLastProgressByTaskID(taskID string) progress.Progress {
	return getLastProgressByTaskIDMock(taskID)
}
func (r mockProgressRepository) UpdateProgress(progress *progress.Progress) error {
	return updateProgressMock(progress)
}

type mockTaskRepository struct {
	MockTaskRepository
}

type mockProgressRepository struct {
	MockProgressRepository
}

func TestCreateTaskShouldCreateNewTaskWithStatusTODO(t *testing.T) {
	r := mockTaskRepository{}
	p := mockProgressRepository{}

	name := "Task 1"

	findTaskByNameAndStatusMock = func(_ string, _ Filter) (Task, error) {
		return Task{}, nil
	}

	createTaskMock = func(task *Task) error {
		if task.Status != TODO {
			t.Errorf("expected task status to be %s, got %s", TODO, task.Status)
		}
		if task.Name != name {
			t.Errorf("expected task name to be %s, got %s", name, task.Name)
		}
		return nil
	}

	s := NewService(r, p)
	err := s.CreateTask(name)

	if err != nil {
		t.Errorf("error was not expected while creating task: %s", err)
	}
}

func TestCreateTaskShouldNotCreateTaskWhenTaskExists(t *testing.T) {
	r := mockTaskRepository{}
	p := mockProgressRepository{}

	findTaskByNameAndStatusMock = func(name string, _ Filter) (Task, error) {
		return Task{
			ID:        "1",
			Name:      name,
			Status:    TODO,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}, nil
	}

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

	statuses := []Status{TODO}

	listTasksByStatusMock = func(_ ...Status) ([]Task, error) {
		return []Task{{
			ID:        "1",
			Name:      "Task",
			Status:    TODO,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}, {
			ID:        "2",
			Name:      "Another Task",
			Status:    TODO,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}}, nil
	}

	tasks, err := s.ListTasksByStatus(statuses...)

	if err != nil {
		t.Errorf("error was not expected while listing tasks: %s", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}
