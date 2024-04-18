package report

import (
	"time"

	"github.com/mathiasdonoso/chrono/internal/chrono/task"
)

type Report struct {
	detail map[task.Task]time.Duration
}

type Service interface {
	DailyReport() (Report, error)
}

type ReportRepository interface {
	GetDailyReport() (Report, error)
}
