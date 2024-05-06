package currency

import (
	"context"

	"gocloud.dev/pubsub"

	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/model"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/subscriber"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/lib/cast"
)

type CurrencyRabbitmqSubscriber struct {
	subscription map[string]*pubsub.Subscription
	repository   repository.CurrencyRepository
	log          *logrus.Logger
}

func (s *CurrencyRabbitmqSubscriber) Run(ctx context.Context) {
	go s.onCurrencyAdded(context.WithoutCancel(ctx))
}

func (s *CurrencyRabbitmqSubscriber) onCurrencyAdded(ctx context.Context) {
	for {
		msg, err := s.subscription[string(constant.CURRENCY_ADDED_EVENT)].Receive(ctx)
		if err != nil {
			s.log.Debugf("Receiving message: %v", err)
			continue
		}

		s.log.Debugf("Got message: %q\n", msg.Body)

		idempotentID := msg.Metadata[constant.EVENT_IDEMPOTENT_ID]
		currency := &model.CurrencyEventModel{}
		if err = cast.FromBytes(msg.Body, currency); err != nil {
			s.log.Warnf("Error casting object (%s): %q\n", idempotentID, msg.Body)
			continue
		}

		if _, err = s.repository.Insert(ctx, currency.ToRepositoryModel()); err != nil {
			s.log.Warnf("Error inserting object (%s): %v\n", idempotentID, currency)
			continue
		}

		msg.Ack()
	}
}

func NewCurrencySubscriber(
	subscription map[string]*pubsub.Subscription,
	repository repository.CurrencyRepository,
	log *logrus.Logger,
) subscriber.Subscriber {
	return &CurrencyRabbitmqSubscriber{
		subscription: subscription,
		repository:   repository,
		log:          log,
	}
}
