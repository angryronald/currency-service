package rabbitmq

import (
	"context"

	"gocloud.dev/pubsub"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/publisher"
)

type RabbitmqPublisher struct {
	pubsubTopics map[string]*pubsub.Topic
}

func (p *RabbitmqPublisher) Publish(ctx context.Context, topic string, idempotentId string, event []byte) error {
	return p.pubsubTopics[topic].Send(ctx, &pubsub.Message{
		Body: event,
		Metadata: map[string]string{
			constant.EVENT_IDEMPOTENT_ID: idempotentId,
		},
	})
}

func NewPublisher(pubsubTopics map[string]*pubsub.Topic) publisher.Publisher {
	return &RabbitmqPublisher{
		pubsubTopics: pubsubTopics,
	}
}
