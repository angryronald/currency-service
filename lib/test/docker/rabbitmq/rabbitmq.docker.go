package rabbitmq

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/lib/test/docker"
)

var mu sync.Mutex
var upFlag bool
var port string
var host string

func init() {
	port = docker.GetAvailablePort(5672)
	host = "localhost"
	upFlag = false
}

// StartRabbitMQ starting rabbitmq in docker and returning active port being used by rabbitmq container
func StartRabbitMQ() string {
	mu.Lock()
	if !upFlag {
		// availablePort := docker.GetAvailablePort(1080)
		// Define the command you want to run
		cmd := exec.Command("docker", "run", "-d", "--rm", "-p", fmt.Sprintf("%s:5672", port), "--name", "rabbitmq", "-e", "RABBITMQ_DEFAULT_USER=user", "-e", "RABBITMQ_DEFAULT_PASS=password", "rabbitmq:3-management")

		// Set the working directory to the location of your docker-compose.yml file
		cmd.Dir = "."

		// Capture and print the command's output
		output, err := cmd.CombinedOutput()
		if err != nil {
			StopRabbitMQWithoutLock()
			logrus.Debug("Error:", err)
			os.Exit(1)
		}

		logrus.Debug("Command Output:", string(output))

		time.Sleep(2 * time.Second)

		cmd = exec.Command("docker", "ps", "-a")

		// Capture and print the command's output
		output, err = cmd.CombinedOutput()
		if err != nil {
			logrus.Debug("Error:", err)
			os.Exit(1)
		}

		logrus.Debug("Command Output:", string(output))

		upFlag = true

		time.Sleep(3 * time.Second)
	}
	mu.Unlock()
	return port
}

func StopRabbitMQWithoutLock() {
	// Define the command you want to run
	cmd := exec.Command("docker", "stop", "rabbitmq")

	// Set the working directory to the location of your docker-compose.yml file
	cmd.Dir = "."

	// Capture and print the command's output
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Debug("Error:", err)
		os.Exit(1)
	}

	logrus.Debug("Command Output:", string(output))
}

func StopRabbitMQ() {
	mu.Lock()
	if upFlag {
		StopRabbitMQWithoutLock()
		upFlag = false
	}
	mu.Unlock()
}
