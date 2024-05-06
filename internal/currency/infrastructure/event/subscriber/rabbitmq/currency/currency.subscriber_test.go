package currency

import (
	"context"
	"testing"
	"time"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
)

var pubsubTopics map[string]*pubsub.Topic
var subscription map[string]*pubsub.Subscription

func init() {
	pubsubTopics = map[string]*pubsub.Topic{}
	subscription = map[string]*pubsub.Subscription{}
}

func TestCurrencyRabbitmqSubscriber_onCurrencyAdded(t *testing.T) {
	pubsubTopics[string(constant.CURRENCY_ADDED_EVENT)] = rabbitpubsub.OpenTopic(rabbitmqClientPublisher, string(constant.CURRENCY_ADDED_EVENT), nil)
	subscription[string(constant.CURRENCY_ADDED_EVENT)] = rabbitpubsub.OpenSubscription(rabbitmqClientConsumer, string(constant.CURRENCY_ADDED_EVENT), nil)

	// Create a testing context.
	ctx := context.TODO()

	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OTPRepository.
	mockRepository := repository.NewMockCurrencyRepository(ctrl)

	// Create a pubsub.Subscription for testing (use mempubsub for testing).
	log := logrus.New()
	subscriber := &CurrencyRabbitmqSubscriber{
		subscription: subscription,
		repository:   mockRepository,
		log:          log,
	}

	// Define a test message.
	testMessage := &pubsub.Message{
		Metadata: map[string]string{
			constant.EVENT_IDEMPOTENT_ID: "test-id",
		},
		Body: []byte(`{"Code":"IDR", "Name":"Indonesian Rupiah"}`), // Add your JSON message here.
	}

	// simulate publishing message
	if err := pubsubTopics[string(constant.CURRENCY_ADDED_EVENT)].Send(ctx, testMessage); err != nil {
		log.Fatalf("error when publishing event: %v", err)
	}

	ctxWithDone, cancel := context.WithDeadline(ctx, time.Now().UTC().Add(5*time.Second))
	defer cancel()

	// Define expectations on the mock OTPRepository.
	mockRepository.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(1)

	// Call the function to test.
	go subscriber.onCurrencyAdded(ctxWithDone)

	// Wait for the operation to complete or the deadline to expire.
	select {
	case <-ctx.Done():
		// The context has been canceled due to the deadline.
		log.Println("Operation canceled due to deadline.")
	case <-time.After(5 * time.Second):
		// This simulates the operation taking longer than the deadline.
		log.Println("Operation completed after deadline.")
	}

	// Assert that the expectations were met.
	assert.Error(t, ctxWithDone.Err(), "context deadline exceeded")
}
