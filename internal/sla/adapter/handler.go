package adapter

import (
	"context"
	"fmt"

	"todoe/domain/task/domain"
	"todoe/internal/event"
)

type CSVHandler struct {
	repo *CSVRepository
}

func NewCSVHandler(repo *CSVRepository) func(context.Context, event.Event) error {
	handler := &CSVHandler{repo: repo}
	return handler.Handle
}

// func (h *CSVHandler) Handle(ctx context.Context, e event.Event) error {
// 	switch e.Type {
// 	case domain.EventStatusCompleted:
// 		task, ok := e.Payload.(domain.Task)
// 		if !ok {
// 			return fmt.Errorf("csv handler: invalid payload for %s", e.Type)
// 		}

// return h.repo.AppendTaskCreated(task)

// 	default:
// 		return nil
// 	}
// }

func (h *CSVHandler) Handle(ctx context.Context, e event.Event) error {
	switch e.Type {
	case domain.EventStatusCompleted:
		task, ok := e.Payload.(domain.Task)
		if !ok {
			return fmt.Errorf("csv handler: invalid payload for %s", e.Type)
		}

		if task.Status != domain.StatusDone {
			return nil
		}

		return h.repo.AppendCompletedTask(task)

	default:
		return nil
	}
}
