package task

import (
	"database/sql"
	"fmt"
	"log"
)

// type DBTX interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// }

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindTaskByName(name string) (*Task, error) {
	task := Task{}

	query := "SELECT id, name, description, status, created_at, updated_at FROM tasks WHERE name = $1;"
	err := r.db.QueryRow(query, name).
		Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Task{}, fmt.Errorf("task with name %s not found", name)
		}

		log.Println("find task by name error:", err)
		return &Task{}, err
	}

	return &task, nil
}

func (r *repository) CreateTask(task *Task) (*Task, error) {
	query := "INSERT INTO tasks (id, name, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := r.db.Exec(query, task.ID, task.Name, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		log.Println("create task error:", err)
		return &Task{}, err
	}

	return task, nil
}

