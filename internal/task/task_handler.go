package task

type Handler struct {
    Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service}
}

// Handles the creation of a new task
func (h *Handler) CreateTask(name, description string) (string, error) {
    res, err := h.Service.CreateTask(name, description)
    if err != nil {
        return "", err
    }

    return res, nil
}

