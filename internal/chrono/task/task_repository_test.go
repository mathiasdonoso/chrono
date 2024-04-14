package task

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTaskShouldInsertNewRowInDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlmock.NewResult(1, 1))

	r := NewRepository(db)

	task := Task{
		ID:          "1",
		Name:        "Task 1",
		Status:      PENDING,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = r.CreateTask(&task)
	if err != nil {
		t.Errorf("error was not expected while creating task: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindPendingTaskByNameShouldFindTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "status", "created_at", "updated_at"}).
		AddRow("1", "Task 1", PENDING, time.Now(), time.Now())

	query := regexp.QuoteMeta("SELECT id, name, status, created_at, updated_at FROM tasks WHERE name = $1 AND status IN ('pending');")

	mock.ExpectQuery(query).WithArgs("Task 1").WillReturnRows(rows)

	r := NewRepository(db)

	task, err := r.FindPendingTaskByName("Task 1", Filter{Statuses: []Status{PENDING}})
	if err != nil {
		t.Errorf("error was not expected while finding task: %s", err)
	}

	if task.ID != "1" {
		t.Errorf("expected task ID to be 1, got %s", task.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindPendingTaskByNameShouldNotFindTaskWhenStatusIsNotPending(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlmock.NewRows([]string{"id", "name", "status", "created_at", "updated_at"}).
		AddRow("1", "Task 1", DONE, time.Now(), time.Now())

	query := regexp.QuoteMeta("SELECT id, name, status, created_at, updated_at FROM tasks WHERE name = $1 AND status IN ('pending');")

	mock.ExpectQuery(query).WithArgs("Task 1").WillReturnError(sql.ErrNoRows)

	r := NewRepository(db)

	_, err = r.FindPendingTaskByName("Task 1", Filter{Statuses: []Status{PENDING}})
	if err == nil {
		t.Errorf("expected error while finding task, got nil")
	}

	if err.Error() != "not found" {
		t.Errorf("expected error to be 'not found', got %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindPendingTaskByNameShouldNotFindTaskWhenNameDoesNotMatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlmock.NewRows([]string{"id", "name", "status", "created_at", "updated_at"}).
		AddRow("1", "Task 2", PENDING, time.Now(), time.Now())

	query := regexp.QuoteMeta("SELECT id, name, status, created_at, updated_at FROM tasks WHERE name = $1 AND status IN ('pending');")

	mock.ExpectQuery(query).WithArgs("Task 1").WillReturnError(sql.ErrNoRows)

	r := NewRepository(db)

	_, err = r.FindPendingTaskByName("Task 1", Filter{Statuses: []Status{PENDING}})
	if err == nil {
		t.Errorf("expected error while finding task, got nil")
	}

	if err.Error() != "not found" {
		t.Errorf("expected error to be 'not found', got %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
