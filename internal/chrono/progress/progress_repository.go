package progress

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) ProgressRepository {
    return &repository{db: db}
}

func (r *repository) AddProgress(progress *Progress) error {
    progress.ID = uuid.New().String()
    progress.CreatedAt = time.Now()
    progress.UpdatedAt = time.Now()

    query := "INSERT INTO progress(id, task_id, status_init, status_end, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);"

    _, err := r.db.Exec(
        query,
        progress.ID,
        progress.TaskID,
        progress.StatusInit,
        progress.StatusFinish,
        progress.CreatedAt,
        progress.UpdatedAt,
    )
    if err != nil {
        return err
    }

    return nil
}
