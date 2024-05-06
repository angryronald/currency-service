package rabbitmq

// uncomment this to try run rabbitmq
// package main

// import (
// 	"fmt"
// 	"log"
// 	"time"
// 

// 	"github.com/ory/dockertest"
// 	// "github.com/streadway/amqp"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func main() {
// 	// Create a new pool of Docker containers
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	// Set up options for RabbitMQ container
// 	options := &dockertest.RunOptions{
// 		Repository: "rabbitmq",
// 		Tag:        "3-management",
// 		Env: []string{
// 			"RABBITMQ_DEFAULT_USER=user",
// 			"RABBITMQ_DEFAULT_PASS=password",
// 		},
// 	}

// 	// Run RabbitMQ container
// 	resource, err := pool.RunWithOptions(options)
// 	if err != nil {
// 		log.Fatalf("Could not start RabbitMQ container: %s", err)
// 	}
// 	defer func() {
// 		// Clean up and remove the RabbitMQ container when done
// 		if err := pool.Purge(resource); err != nil {
// 			log.Fatalf("Could not purge RabbitMQ container: %s", err)
// 		}
// 	}()
// 	fmt.Printf("RabbitMQ container started: ID=%s\n", resource.Container.ID)

// 	// Sleep for some time to allow RabbitMQ to run
// 	time.Sleep(10 * time.Second) // Adjust the sleep duration as needed

// 	// Wait for RabbitMQ to be ready
// 	if err := pool.Retry(func() error {
// 		conn, err := amqp.Dial(fmt.Sprintf(
// 			"amqp://user:password@%s:%s/",
// 			resource.GetBoundIP("5672/tcp"),
// 			resource.GetPort("5672/tcp"),
// 		))
// 		if err != nil {
// 			return err
// 		}
// 		defer conn.Close()
// 		return nil
// 	}); err != nil {
// 		log.Fatalf("Could not connect to RabbitMQ: %s", err)
// 	}

// 	// At this point, RabbitMQ should be running and ready to use.

// 	// You can use the streadway/amqp library to interact with RabbitMQ.
// 	// For example, you can declare a queue and publish a message:
// 	conn, err := amqp.Dial(fmt.Sprintf(
// 		"amqp://user:password@%s:%s/",
// 		resource.GetBoundIP("5672/tcp"),
// 		resource.GetPort("5672/tcp"),
// 	))
// 	if err != nil {
// 		log.Fatalf("Could not connect to RabbitMQ: %s", err)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("Could not open a channel: %s", err)
// 	}
// 	defer ch.Close()

// 	queueName := "test-queue"
// 	_, err = ch.QueueDeclare(
// 		queueName,
// 		false, // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	if err != nil {
// 		log.Fatalf("Could not declare the queue: %s", err)
// 	}

// 	messageBody := "Hello, RabbitMQ!"
// 	err = ch.Publish(
// 		"",        // exchange
// 		queueName, // routing key
// 		false,     // mandatory
// 		false,     // immediate
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(messageBody),
// 		})
// 	if err != nil {
// 		log.Fatalf("Could not publish a message: %s", err)
// 	}

// 	fmt.Printf("Published message: %s\n", messageBody)

// 	// You can now consume messages, perform integration tests, and clean up resources as needed.
// 	// Remember to handle errors appropriately in a real-world scenario.
// }
