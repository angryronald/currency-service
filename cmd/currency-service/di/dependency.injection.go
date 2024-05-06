package di

import (
	"sync"

	"gocloud.dev/pubsub"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	currencyCommand "github.com/angryronald/currency-service/internal/currency/application/command"
	currencyQuery "github.com/angryronald/currency-service/internal/currency/application/query"
	currencyService "github.com/angryronald/currency-service/internal/currency/domain/currency"
	currencyEndpoint "github.com/angryronald/currency-service/internal/currency/endpoint"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/publisher/rabbitmq"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/subscriber"
	currencySubscriber "github.com/angryronald/currency-service/internal/currency/infrastructure/event/subscriber/rabbitmq/currency"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	currencyMemcached "github.com/angryronald/currency-service/internal/currency/infrastructure/repository/memcached/redis/currency"
	currencySQL "github.com/angryronald/currency-service/internal/currency/infrastructure/repository/sql/postgre/currency"
)

type Dependencies struct {
	CurrencyEndpoint            currencyEndpoint.CurrencyEndpointInterface
	CurrencySubscriber          subscriber.Subscriber
	CurrencyMemcachedRepository repository.CurrencyRepository
	CurrencySQLRepository       repository.CurrencyRepository
	SQLDB                       *gorm.DB
}

var syncOnce sync.Once
var AllDependencies Dependencies
var redisClient *redis.Client
var redisDefaultExpiration int
var postgresDB *gorm.DB
var pubsubTopics map[string]*pubsub.Topic
var subscriptions map[string]*pubsub.Subscription

func CollectDependencies(log *logrus.Logger) {
	syncOnce.Do(func() {
		instantiateRedis(log)
		instantiatePostgres(log)
		instantiatePubSub(log)

		currencyMemcachedRepository := currencyMemcached.NewCurrencyRepository(
			redisClient, redisDefaultExpiration,
		)

		currencySQLRepository := currencySQL.NewCurrencyRepository(
			postgresDB,
		)

		AllDependencies = Dependencies{
			CurrencyEndpoint: currencyEndpoint.NewCurrencyEndpoint(
				currencyCommand.NewCurrencyCommand(
					currencyService.NewCurrencyService(
						rabbitmq.NewPublisher(
							pubsubTopics,
						),
						currencySQLRepository,
						log,
					),
				),
				currencyQuery.NewCurrencyQuery(
					currencyService.NewCurrencyService(
						rabbitmq.NewPublisher(
							pubsubTopics,
						),
						currencyMemcachedRepository,
						log,
					),
					currencyService.NewCurrencyService(
						rabbitmq.NewPublisher(
							pubsubTopics,
						),
						currencySQLRepository,
						log,
					),
					log,
				),
				log,
			),
			CurrencySubscriber: currencySubscriber.NewCurrencySubscriber(
				subscriptions,
				currencyMemcachedRepository,
				log,
			),
			CurrencyMemcachedRepository: currencyMemcachedRepository,
			CurrencySQLRepository:       currencySQLRepository,
			SQLDB:                       postgresDB,
		}
	})
}
