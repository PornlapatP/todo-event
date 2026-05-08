package application

import (
	"context"

	"todoe/domain/task/domain"
	"todoe/domain/task/port"
	"todoe/internal/event"
)

func NewProjectionHandler(viewRepo port.ViewRepository) func(context.Context, event.Event) error {
	return func(ctx context.Context, e event.Event) error {
		switch e.Type {
		case domain.EventCreated, domain.EventStatusChanged:
			task, ok := e.Payload.(domain.Task)
			if !ok {
				return nil
			}
			result := viewRepo.Upsert(ctx, task)
			if result.IsError() {
				return result.Error()
			}
		}
		return nil
	}
}
