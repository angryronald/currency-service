package currency

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/lib/test/docker/rabbitmq"
)

var rabbitmqClientConsumer *amqp.Connection
var rabbitmqClientPublisher *amqp.Connection

// Uncomment when using dockertest
// var resource *dockertest.Resource

// Using dockertest but failing on starting rabbitmq in container
// func TestMain(m *testing.M) {
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		logrus.Fatalf("Could not connect to Docker: %s\n", err)
// 	}

// 	rabbitmqClientPublisher, resource = rabbitmq.GetConnection(pool)
// 	if rabbitmqClientPublisher != nil {
// 		logrus.Debugf("Success create rabbitmq connection (publisher)\n")
// 	}

// 	rabbitmqClientConsumer, _ = rabbitmq.GetConnection(pool)
// 	if rabbitmqClientPublisher != nil {
// 		logrus.Debugf("Success create rabbitmq connection (consumer)\n")
// 	}

// 	chPublisher, err := rabbitmqClientPublisher.Channel()
// 	if err != nil {
// 		logrus.Fatalf("Could not open a channel: %s", err)
// 	}
// 	defer chPublisher.Close()

// 	chConsumer, err := rabbitmqClientConsumer.Channel()
// 	if err != nil {
// 		logrus.Fatalf("Could not open a channel: %s", err)
// 	}
// 	defer chConsumer.Close()

// 	err = chConsumer.ExchangeDeclare(
// 		string(constants.OTP_GENERATED_EVENT),
// 		"fanout",
// 		true,  // Durable (exchange survives server restarts)
// 		false, // Auto-deleted (exchange is deleted when no longer in use)
// 		false, // Internal (used by other exchanges but not clients)
// 		false, // No-wait (do not wait for a server response)
// 		nil,   // Arguments (optional)
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	err = chConsumer.ExchangeDeclare(
// 		string(constants.OTP_VALIDATED_EVENT),
// 		"fanout",
// 		true,  // Durable (exchange survives server restarts)
// 		false, // Auto-deleted (exchange is deleted when no longer in use)
// 		false, // Internal (used by other exchanges but not clients)
// 		false, // No-wait (do not wait for a server response)
// 		nil,   // Arguments (optional)
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	_, err = chConsumer.QueueDeclare(
// 		string(constants.OTP_GENERATED_EVENT),
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	_, err = chConsumer.QueueDeclare(
// 		string(constants.OTP_VALIDATED_EVENT),
// 		true,  // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	err = chPublisher.ExchangeDeclare(
// 		string(constants.OTP_GENERATED_EVENT),
// 		"fanout",
// 		true,  // Durable (exchange survives server restarts)
// 		false, // Auto-deleted (exchange is deleted when no longer in use)
// 		false, // Internal (used by other exchanges but not clients)
// 		false, // No-wait (do not wait for a server response)
// 		nil,   // Arguments (optional)
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	err = chPublisher.ExchangeDeclare(
// 		string(constants.OTP_VALIDATED_EVENT),
// 		"fanout",
// 		true,  // Durable (exchange survives server restarts)
// 		false, // Auto-deleted (exchange is deleted when no longer in use)
// 		false, // Internal (used by other exchanges but not clients)
// 		false, // No-wait (do not wait for a server response)
// 		nil,   // Arguments (optional)
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	err = chConsumer.QueueBind(
// 		string(constants.OTP_GENERATED_EVENT),
// 		"event",
// 		string(constants.OTP_GENERATED_EVENT),
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	err = chConsumer.QueueBind(
// 		string(constants.OTP_VALIDATED_EVENT),
// 		"event",
// 		string(constants.OTP_VALIDATED_EVENT),
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		logrus.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	code := m.Run() // execute all the tests

// 	// Delete the Docker container
// 	if err := pool.Purge(resource); err != nil {
// 		logrus.Fatalf("Could not purge rabbitmq resource: %s\n", err)
// 	}

// 	defer os.Exit(code)
// }

// Using docker command
func TestMain(m *testing.M) {
	var err error

	port := rabbitmq.StartRabbitMQ()

	// port := "5672"
	time.Sleep(10 * time.Second)

	if rabbitmqClientPublisher, err = amqp.Dial(fmt.Sprintf(
		"amqp://user:password@localhost:%s/",
		port,
	)); err != nil {
		logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	if rabbitmqClientConsumer, err = amqp.Dial(fmt.Sprintf(
		"amqp://user:password@localhost:%s/",
		port,
	)); err != nil {
		logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	chPublisher, err := rabbitmqClientPublisher.Channel()
	if err != nil {
		logrus.Fatalf("Could not open a channel: %s", err)
	}

	chConsumer, err := rabbitmqClientConsumer.Channel()
	if err != nil {
		logrus.Fatalf("Could not open a channel: %s", err)
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
		logrus.Fatalf("Could not declare the queue: %s", err)
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
		logrus.Fatalf("Could not declare the queue: %s", err)
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
		logrus.Fatalf("Could not declare the queue: %s", err)
	}

	err = chConsumer.QueueBind(
		string(constant.CURRENCY_ADDED_EVENT),
		"event",
		string(constant.CURRENCY_ADDED_EVENT),
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalf("Could not declare the queue: %s", err)
	}

	code := m.Run() // execute all the tests

	defer os.Exit(code)
	defer rabbitmq.StopRabbitMQ()
	defer chPublisher.Close()
	defer chConsumer.Close()
}
