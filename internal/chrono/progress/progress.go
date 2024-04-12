package progress

import "time"

type Progress struct {
    ID string `json:"id"`
    TaskID string `json:"task_id"`
    StatusInit string `json:"status_init"`
    StatusFinish string `json:"status_finish"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    FinishedAt time.Time `json:"finished_at"`
}

type ProgressRepository interface {
    // AddProgress adds a new progress to the database.
    AddProgress(progress *Progress) error
}
