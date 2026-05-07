package event

import (
	"context"
	"log/slog"
	"sync"
)

type Event struct {
	Type    string
	Payload any
}

type Publisher interface {
	Publish(ctx context.Context, e Event)
}

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]func(context.Context, Event) error
}

func NewEventBus() *EventBus {
	return &EventBus{handlers: make(map[string][]func(context.Context, Event) error)}
}

func (b *EventBus) Subscribe(eventType string, fn func(context.Context, Event) error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], fn)
	slog.Debug("eventbus: handler subscribed", "event_type", eventType)
}

func (b *EventBus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	handlers := b.handlers[e.Type]
	b.mu.RUnlock()

	for _, fn := range handlers {
		if err := fn(ctx, e); err != nil {
			slog.Error("eventbus: handler error", "event_type", e.Type, "err", err)
		}
		slog.Debug("eventbus: event published", "event_type", e.Type)
	}
}

var _ Publisher = (*EventBus)(nil)

// package event

// import (
// 	"context"
// 	"log/slog"
// 	"sync"
// )

// type EventType string

// type Event struct {
// 	Type    EventType
// 	Payload any
// }

// type Publisher interface {
// 	Publish(ctx context.Context, e Event)
// }

// type EventHandler func(context.Context, Event) error

// type EventBus struct {
// 	mu       sync.RWMutex
// 	handlers map[EventType][]EventHandler
// }

// func NewEventBus() *EventBus {
// 	return &EventBus{
// 		handlers: make(map[EventType][]EventHandler),
// 	}
// }

// func (b *EventBus) Subscribe(eventType EventType, fn EventHandler) {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()

// 	b.handlers[eventType] = append(b.handlers[eventType], fn)

// 	slog.Debug("eventbus: handler subscribed", "event_type", eventType)
// }

// func (b *EventBus) Publish(ctx context.Context, e Event) {
// 	b.mu.RLock()
// 	handlers := b.handlers[e.Type]
// 	b.mu.RUnlock()

// 	if len(handlers) == 0 {
// 		slog.Debug("eventbus: no handlers found", "event_type", e.Type)
// 		return
// 	}

// 	for _, fn := range handlers {
// 		func(h EventHandler) {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					slog.Error("eventbus: handler panic",
// 						"event_type", e.Type,
// 						"panic", r,
// 					)
// 				}
// 			}()

// 			if err := h(ctx, e); err != nil {
// 				slog.Error("eventbus: handler error",
// 					"event_type", e.Type,
// 					"err", err,
// 				)
// 			}
// 		}(fn)
// 	}

// 	slog.Debug("eventbus: event published", "event_type", e.Type)
// }

// var _ Publisher = (*EventBus)(nil)
