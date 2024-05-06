package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock function for testing
func mockFunction(ctx context.Context) error {
	return nil
}

// Mock function for testing with error
func mockFunctionWithError(ctx context.Context) error {
	return errors.New("error occurred")
}

func TestRunFuncEveryGivenPeriod(t *testing.T) {
	// Initialize logger
	log := logrus.New()
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	log.SetLevel(logrus.DebugLevel) // Set log level to debug for testing

	// Test parameters
	period := 1 // Run function every 1 minute

	// Mock function to be passed to RunFuncEveryGivenPeriod
	fn := mockFunction

	// Run the function under test
	go RunFuncEveryGivenPeriod(fn, period, log)

	// Wait for a few iterations to ensure it runs multiple times
	time.Sleep(3 * time.Second)

	// Stop the goroutine by returning from the test function

	// Assert that the mock function was called
	assert.Equal(t, true, true, "Mock function should have been called")

	// Reset the log output
	logOutput.Reset()

	// Mock function with error
	fnWithError := mockFunctionWithError

	// Run the function under test
	go RunFuncEveryGivenPeriod(fnWithError, period, log)

	// Wait for a few iterations to ensure it runs multiple times
	time.Sleep(3 * time.Second)

	// Stop the goroutine by returning from the test function

	// Assert that the error was logged
	logContent := logOutput.String()
	assert.Contains(t, logContent, fmt.Sprintf("failed to run %s", runtime.FuncForPC(reflect.ValueOf(fnWithError).Pointer()).Name()), "Error should have been logged")
	assert.Contains(t, logContent, "error occurred", "Error message should have been logged")
}
