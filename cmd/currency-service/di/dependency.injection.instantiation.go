package di

import (
	"fmt"
	"strconv"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/config"
	"github.com/angryronald/currency-service/internal/currency/constant"
)

func instantiateRedis(log *logrus.Logger) {
	db, err := strconv.Atoi(config.GetValue(config.REDIS_DB))
	if err != nil {
		log.Fatalf("cannot connecting to redis: %v", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetValue(config.REDIS_ADDRESS),
		Password: config.GetValue(config.REDIS_PASSWORD), // no password set
		DB:       db,                                     // use default DB
	})

	redisDefaultExpirationRaw, err := strconv.Atoi(config.GetValue(config.REDIS_DEFAULT_EXPIRATION))
	if err != nil {
		// default 60 minutes
		redisDefaultExpirationRaw = 3600
	}
	redisDefaultExpiration = redisDefaultExpirationRaw
}

func instantiatePostgres(log *logrus.Logger) {
	var err error
	dsn := fmt.Sprintf(
		config.GetValue(config.DATABASE_CONNECTION_STRING),
		config.GetValue(config.DATABASE_HOST),
		config.GetValue(config.DATABASE_USER),
		config.GetValue(config.DATABASE_PASS),
		config.GetValue(config.DATABASE_NAME),
		config.GetValue(config.DATABASE_PORT),
		config.GetValue(config.DATABASE_SSL),
	)
	postgresDB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("cannot connecting to postgres: %v", err)
	}
}

func instantiatePubSub(log *logrus.Logger) {
	var rabbitmqClientConsumer *amqp.Connection
	var rabbitmqClientPublisher *amqp.Connection
	var err error

	pubsubTopics = map[string]*pubsub.Topic{}
	subscriptions = map[string]*pubsub.Subscription{}

	if rabbitmqClientPublisher, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.GetValue(config.RABBITMQ_USERNAME),
		config.GetValue(config.RABBITMQ_PASSWORD),
		config.GetValue(config.RABBITMQ_HOST),
		config.GetValue(config.RABBITMQ_PORT),
	)); err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	if rabbitmqClientConsumer, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.GetValue(config.RABBITMQ_USERNAME),
		config.GetValue(config.RABBITMQ_PASSWORD),
		config.GetValue(config.RABBITMQ_HOST),
		config.GetValue(config.RABBITMQ_PORT),
	)); err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	chPublisher, err := rabbitmqClientPublisher.Channel()
	if err != nil {
		log.Fatalf("Could not open a channel: %s", err)
	}

	chConsumer, err := rabbitmqClientConsumer.Channel()
	if err != nil {
		log.Fatalf("Could not open a channel: %s", err)
	}

	err = chConsumer.ExchangeDeclare(
		string(constant.CURRENCY_ADDED_EVENT),
		"fanout",
		true,  // Durable (exchange survives server restarts)
		false, // Auto-deleted (exchange is deleted when no longer in use)
		false, // Internal (used by other exchanges but not clients)
		false, // No-wait (do not wait for a server response)
		nil,   // Arguments (optional)
	)
	if err != nil {
		log.Fatalf("Could not declare the queue: %s", err)
	}

	_, err = chConsumer.QueueDeclare(
		string(constant.CURRENCY_ADDED_EVENT),
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Could not declare the queue: %s", err)
	}

	err = chPublisher.ExchangeDeclare(
		string(constant.CURRENCY_ADDED_EVENT),
		"fanout",
		true,  // Durable (exchange survives server restarts)
		false, // Auto-deleted (exchange is deleted when no longer in use)
		false, // Internal (used by other exchanges but not clients)
		false, // No-wait (do not wait for a server response)
		nil,   // Arguments (optional)
	)
	if err != nil {
		log.Fatalf("Could not declare the queue: %s", err)
	}

	err = chConsumer.QueueBind(
		string(constant.CURRENCY_ADDED_EVENT),
		"event",
		string(constant.CURRENCY_ADDED_EVENT),
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Could not declare the queue: %s", err)
	}

	pubsubTopics[string(constant.CURRENCY_ADDED_EVENT)] = rabbitpubsub.OpenTopic(rabbitmqClientPublisher, string(constant.CURRENCY_ADDED_EVENT), nil)
	subscriptions[string(constant.CURRENCY_ADDED_EVENT)] = rabbitpubsub.OpenSubscription(rabbitmqClientConsumer, string(constant.CURRENCY_ADDED_EVENT), nil)
}
