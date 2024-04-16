package progress

import "time"

type Progress struct {
    ID string `json:"id"`
    TaskID string `json:"task_id"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    FinishedAt time.Time `json:"finished_at"`
}

type ProgressRepository interface {
    // AddProgress adds a new progress to the database.
    AddProgress(progress *Progress) error
    // GetLastProgressByTaskID returns the last progress of a task.
    GetLastProgressByTaskID(taskID string) Progress
    // UpdateProgress updates a progress in the database.
    UpdateProgress(progress *Progress) error
}
