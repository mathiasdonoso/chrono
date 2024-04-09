package task

import (
	"database/sql"
	"fmt"
	"strings"
)

type repository struct {
	db *sql.DB
}

type Filter struct {
	Statuses []Status
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	s := make([]string, len(statuses))
	for i, v := range statuses {
		s[i] = fmt.Sprintf("'%s'", v)
	}

	query := "SELECT id, name, description, status, created_at, updated_at FROM tasks WHERE status IN (" + strings.Join(s, ",") + ");"

	rows, err := r.db.Query(query)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []Task{}, err
	}

	return tasks, nil
}

func (r *repository) FindTaskByPartialId(partialId string, filter Filter) (*Task, error) {
	task := Task{}

	query := "SELECT id, name, description, status, created_at, updated_at FROM tasks WHERE id LIKE $1"

	if len(filter.Statuses) > 0 {
		s := make([]string, len(filter.Statuses))
		for i, v := range filter.Statuses {
			s[i] = fmt.Sprintf("'%s'", v)
		}
		query += " AND status IN (" + strings.Join(s, ",") + ")"
	}

	query += ";"

	err := r.db.QueryRow(query, partialId+"%").Scan(
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

func (r *repository) RemoveTaskById(id string) error {
	query := "DELETE FROM tasks WHERE id = $1;"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

