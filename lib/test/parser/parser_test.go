package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSuccessObjectToBytes(t *testing.T) {
	// Create a sample success response object
	expected := struct {
		Status string
		Data   struct {
			Message string
			Value   int
		}
	}{
		Status: "",
		Data: struct {
			Message string
			Value   int
		}{
			Message: "",
			Value:   0,
		},
	}

	// Convert the success response to bytes using the ParseSuccessObjectToBytes function
	bytes, err := ParseSuccessObjectToBytes(expected)

	// Check if there was an error during conversion
	assert.NoError(t, err)

	// Unmarshal the bytes back into a response model
	var actual struct {
		Status string
		Data   struct {
			Message string
			Value   int
		}
	}
	err = json.Unmarshal(bytes, &actual)

	// Check if there was an error during unmarshaling
	assert.NoError(t, err)

	// Assert that the unmarshaled result matches the expected value
	assert.Equal(t, expected, actual)
}
