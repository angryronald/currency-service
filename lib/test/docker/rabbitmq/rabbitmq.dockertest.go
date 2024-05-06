package rabbitmq

import (
	"fmt"
	"time"

	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
	// "github.com/streadway/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
)

var resource *dockertest.Resource

func GenerateInstance(pool *dockertest.Pool) *dockertest.Resource {
	var err error

	// Set up options for RabbitMQ container
	options := &dockertest.RunOptions{
		Repository: "rabbitmq",
		Tag:        "3-management",
		Env: []string{
			"RABBITMQ_DEFAULT_USER=user",
			"RABBITMQ_DEFAULT_PASS=password",
		},
	}

	// Run RabbitMQ container
	resource, err = pool.RunWithOptions(options)
	if err != nil {
		logrus.Fatalf("Could not start RabbitMQ container: %s", err)
	}

	time.Sleep(10 * time.Second)

	// Wait for RabbitMQ to be ready
	// if err := pool.Retry(func() error {
	// 	conn, err := amqp.Dial(fmt.Sprintf(
	// 		"amqp://user:password@%s:%s/",
	// 		resource.GetBoundIP("5672/tcp"),
	// 		resource.GetPort("5672/tcp"),
	// 	))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer conn.Close()
	// 	return nil
	// }); err != nil {
	// 	logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	// }
	if err := pool.Retry(func() error {
		conn, err := amqp.Dial(fmt.Sprintf(
			"amqp://user:password@%s:%s/",
			host,
			port,
		))
		if err != nil {
			return err
		}
		defer conn.Close()
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	// At this point, RabbitMQ should be running and ready to use.

	// You can use the streadway/amqp library to interact with RabbitMQ.
	// For example, you can declare a queue and publish a message:
	// conn, err := amqp.Dial(fmt.Sprintf(
	// 	"amqp://user:password@%s:%s/",
	// 	host,
	// 	port,
	// ))
	// if err != nil {
	// 	logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	// }

	// ch, err := conn.Channel()
	// if err != nil {
	// 	logrus.Fatalf("Could not open a channel: %s", err)
	// }

	return resource
}

func GetConnection(pool *dockertest.Pool) (*amqp.Connection, *dockertest.Resource) {
	if resource == nil {
		resource = GenerateInstance(pool)
	}

	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://user:password@%s:%s/",
		host,
		port,
	))
	if err != nil {
		logrus.Fatalf("Could not connect to RabbitMQ: %s", err)
	}

	return conn, resource
}
