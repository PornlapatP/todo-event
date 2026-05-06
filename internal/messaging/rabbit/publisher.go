package rabbit

import (
	"context"
	"encoding/json"

	"todoe/internal/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewRabbitPublisher(ch *amqp.Channel) *RabbitPublisher {
	return &RabbitPublisher{
		ch:       ch,
		exchange: "captcha.exchange",
	}
}

func (r *RabbitPublisher) Publish(ctx context.Context, evt event.Event) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	return r.ch.PublishWithContext(
		context.Background(),
		r.exchange,
		evt.Type,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
