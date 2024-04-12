package task

type Handler struct {
    Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service}
}

func (h *Handler) CreateTask(name, description string) (string, error) {
    err := h.Service.CreateTask(name, description)
    if err != nil {
        return "", err
    }

    return "Task created", nil
}

func (h *Handler) ListTasksByStatus(statuses ...Status) ([]Task, error) {
    tasks, err := h.Service.ListTasksByStatus(statuses...)
    if err != nil {
        return []Task{}, err
    }

    return tasks, nil
}

func (h *Handler) RemoveTaskByPartialId(partialId string) (string, error) {
    err := h.Service.RemoveTaskByPartialId(partialId)
    if err != nil {
        return "", err
    }

    return "Task removed", nil
}

func (h *Handler) StartTask(idOrName string) (string, error) {
    err := h.Service.StartTask(idOrName)
    if err != nil {
        return "", err
    }

    return "Task started", nil
}

