package report

type Handler struct {
    Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service}
}

func (h *Handler) CreateReport(reportType string) (Report, error) {
    switch reportType {
    case "daily":
        res, err := h.Service.DailyReport()
        if err != nil {
            return Report{}, err
        }
        return res, nil
    default:
        return Report{}, nil
    }
}
