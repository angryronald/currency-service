package rabbitmq

import (
	"context"
	"testing"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var pubsubTopics map[string]*pubsub.Topic
var subscription map[string]*pubsub.Subscription

func init() {
	pubsubTopics = map[string]*pubsub.Topic{}
	subscription = map[string]*pubsub.Subscription{}
}

func TestRabbitmqPublisher_Publish(t *testing.T) {
	pubsubTopics[topicName] = rabbitpubsub.OpenTopic(rabbitmqClientPublisher, topicName, nil)
	subscription[topicName] = rabbitpubsub.OpenSubscription(rabbitmqClientConsumer, topicName, nil)

	// Create a testing context.
	ctx := context.TODO()

	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publisher := NewPublisher(pubsubTopics)

	if err := publisher.Publish(ctx, topicName, uuid.NewString(), []byte(`{"ID":"123456"}`)); err != nil {
		t.Errorf("error when publishing event: %v", err)
	}

	// Call the function to test.
	for {
		msg, err := subscription[topicName].Receive(ctx)
		if err != nil {
			t.Errorf("error when consuming event: %v", err)
		}

		msg.Ack()
		ctx.Done()
		break
	}

	// Assert that the expectations were met.
	assert.Nil(t, ctx.Err())
}
