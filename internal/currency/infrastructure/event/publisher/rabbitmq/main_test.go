package rabbitmq

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/lib/test/docker/rabbitmq"
)

var rabbitmqClientConsumer *amqp.Connection
var rabbitmqClientPublisher *amqp.Connection
var topicName string = "CONTOH.TOPIC"

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
		topicName,
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
		topicName,
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
		topicName,
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
		topicName,
		"event",
		topicName,
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
