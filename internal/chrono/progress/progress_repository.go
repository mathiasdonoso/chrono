package progress

import "database/sql"

type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) ProgressRepository {
    return &repository{db: db}
}

func (r *repository) AddProgress(progress Progress) error {
    query := "INSERT INTO works (id, task_id, status_init, status_end, created_at, updated_at, finished_at) VALUES ($1, $2, $3, $4, $5, $6, $7);"

    _, err := r.db.Exec(
        query,
        progress.ID,
        progress.TaskID,
        progress.StatusInit,
        progress.StatusFinish,
        progress.CreatedAt,
        progress.UpdatedAt,
        progress.FinishedAt,
    )
    if err != nil {
        return err
    }

    return nil
}
