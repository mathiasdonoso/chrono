package progress

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type uuidCreator interface {
	NewUUID() string
}

type uuidGenerator struct{}

func (u *uuidGenerator) NewUUID() string {
	return uuid.New().String()
}

var idGenerator uuidCreator

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ProgressRepository {
	return &repository{db: db}
}

func (r *repository) GetLastProgressByTaskID(taskID string) Progress {
	query := "SELECT id, task_id, status, created_at, updated_at FROM progress WHERE task_id = $1 ORDER BY created_at DESC LIMIT 1;"

	row := r.db.QueryRow(query, taskID)

	var progress Progress
	if err := row.Scan(
		&progress.ID,
		&progress.TaskID,
		&progress.Status,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	); err != nil {
		return Progress{}
	}

	return progress
}

func (r *repository) AddProgress(progress *Progress) error {
	progress.ID = idGenerator.NewUUID()
	progress.CreatedAt = time.Now()
	progress.UpdatedAt = time.Now()

	query := "INSERT INTO progress(id, task_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);"

	_, err := r.db.Exec(
		query,
		progress.ID,
		progress.TaskID,
		progress.Status,
		progress.CreatedAt,
		progress.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateProgress(progress *Progress) error {
	progress.UpdatedAt = time.Now()

	query := "UPDATE progress SET status = $1, updated_at = $2 WHERE id = $3;"

	_, err := r.db.Exec(
		query,
		progress.Status,
		progress.UpdatedAt,
		progress.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
