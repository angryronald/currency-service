package config

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	// Call the GetLogger function
	logger := GetLogger()

	// Ensure that the logger is not nil
	if logger == nil {
		t.Error("Logger is nil")
	}

	// Check the formatter
	formatter, ok := logger.Formatter.(*logrus.TextFormatter)
	if !ok {
		t.Error("Invalid formatter")
	}
	if !formatter.FullTimestamp {
		t.Error("FullTimestamp is not set to true")
	}

	// Check the output
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)
	logger.Info("Test log message")
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Test log message") {
		t.Error("Log output doesn't contain expected message")
	}
}

func TestGetValue(t *testing.T) {
	// Set up a mock secret
	os.Setenv("test_key", "test_value")

	// Test retrieving an existing value
	result := GetValue("test_key")
	assert.Equal(t, "test_value", result)
}

func TestGetAllowedClients(t *testing.T) {
	// Test getting the working directory
	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")
	result := GetAllowedClients()
	assert.NotNil(t, result)
}
