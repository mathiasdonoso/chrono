package task

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

type Filter struct {
	Statuses []Status
}

func NewRepository(db *sql.DB) TaskRepository {
	return &repository{db: db}
}

func (r *repository) ListTasksByStatus(statuses ...Status) ([]Task, error) {
	s := make([]string, len(statuses))
	for i, v := range statuses {
		s[i] = fmt.Sprintf("'%s'", v)
	}

	query := "SELECT id, name, status, created_at, updated_at FROM tasks WHERE status IN (" + strings.Join(s, ",") + ");"

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

func (r *repository) FindTaskByPartialId(partialId string, filter Filter) (Task, error) {
	query := "SELECT id, name, status, created_at, updated_at FROM tasks WHERE id LIKE $1"

	if len(filter.Statuses) > 0 {
		s := make([]string, len(filter.Statuses))
		for i, v := range filter.Statuses {
			s[i] = fmt.Sprintf("'%s'", v)
		}
		query += " AND status IN (" + strings.Join(s, ",") + ")"
	}

	query += " LIMIT 2;"

	rows, err := r.db.Query(query, partialId+"%")
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("not found")
		}

		return Task{}, err
	}

	result := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return Task{}, err
		}
		result = append(result, t)
	}

	if len(result) > 1 {
		return Task{}, fmt.Errorf("multiple tasks found")
	}

	return result[0], nil
}

func (r *repository) UpdateTask(task *Task) error {
	task.UpdatedAt = time.Now()

	query := "UPDATE tasks SET name = ?, status = ?, updated_at = ? WHERE id = ?;"
	_, err := r.db.Exec(
		query,
		task.Name,
		task.Status,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByIdOrCreate(idOrName string, filter Filter) (Task, error) {
	task := Task{}

	fq := ""

	if len(filter.Statuses) > 0 {
		s := make([]string, len(filter.Statuses))
		for i, v := range filter.Statuses {
			s[i] = fmt.Sprintf("'%s'", v)
		}
		fq += " AND status IN (" + strings.Join(s, ",") + ")"
	}

	query := fmt.Sprintf(
		"SELECT id, name, status, created_at, updated_at FROM tasks WHERE id LIKE $1 %s;",
		fq,
	)

	err := r.db.QueryRow(query, idOrName+"%").Scan(
		&task.ID,
		&task.Name,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			newTask := Task{}
			newTask.Name = idOrName

			err = r.CreateTask(&newTask)

			if err != nil {
				return Task{}, fmt.Errorf("not found")
			}

			return newTask, nil
		} else {
			return Task{}, fmt.Errorf("not found")
		}
	}

	return task, nil
}

func (r *repository) FindTaskById(id string) (Task, error) {
	task := Task{}

	query := "SELECT id, name, status, created_at, updated_at FROM tasks WHERE id = $1;"
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("not found")
		}

		return Task{}, err
	}

	return task, nil
}

func (r *repository) FindPendingTaskByName(name string, filter Filter) (Task, error) {
	task := Task{}

	query := "SELECT id, name, status, created_at, updated_at FROM tasks WHERE name = $1"

	if len(filter.Statuses) > 0 {
		s := make([]string, len(filter.Statuses))
		for i, v := range filter.Statuses {
			s[i] = fmt.Sprintf("'%s'", v)
		}
		query += " AND status IN (" + strings.Join(s, ",") + ")"
	}

	query += ";"

	err := r.db.QueryRow(query, name).Scan(
		&task.ID,
		&task.Name,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("not found")
		}

		return Task{}, err
	}

	return task, nil
}

func (r *repository) CreateTask(task *Task) error {
	task.ID = uuid.New().String()
	if task.Status == "" {
		task.Status = PENDING
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	query := "INSERT INTO tasks (id, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);"
	_, err := r.db.Exec(
		query,
		task.ID,
		task.Name,
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

