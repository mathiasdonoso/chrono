package report

type service struct {
	ReportRepository ReportRepository
}

func NewService(r ReportRepository) Service {
	return &service{
		ReportRepository: r,
	}
}

func (s *service) DailyReport() (Report, error) {
	// Get all task that I been working the previous working day (with the progress)
	// (previous working day is the last day that I worked)
	report, err := s.ReportRepository.GetDailyReport()
	if err != nil {
		return Report{}, err
	}

	// build the report
	// report := Report{
	// 	detail: map[task.Task]time.Duration{},
	// }

	return report, nil
}
