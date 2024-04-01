package task

import (
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindTaskById(id string) (*Task, error) {
	task := Task{}

	query := "SELECT id, name, description, status, created_at, updated_at FROM tasks WHERE id = $1;"
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &Task{}, fmt.Errorf("not found")
		}

		return &Task{}, err
	}

	return &task, nil
}

func (r *repository) FindPendingTaskByName(name string) (*Task, error) {
	task := Task{}

	query := "SELECT id, name, description, status, created_at, updated_at FROM tasks WHERE name = $1 AND status = 'pending';"
	err := r.db.QueryRow(query, name).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &Task{}, fmt.Errorf("not found")
		}

		return &Task{}, err
	}

	return &task, nil
}

func (r *repository) CreateTask(task *Task) error {
	query := "INSERT INTO tasks (id, name, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := r.db.Exec(
		query,
		task.ID,
		task.Name,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

