package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	t.Run("Valid JSON File", func(t *testing.T) {
		// Create a temporary JSON file with valid JSON content
		validJSON := `{"name": "John", "age": 30}`
		tmpfile, err := ioutil.TempFile("", "test.json")
		if err != nil {
			t.Fatalf("Error creating temporary JSON file: %v", err)
		}
		defer os.Remove(tmpfile.Name()) // Clean up

		if _, err := tmpfile.Write([]byte(validJSON)); err != nil {
			t.Fatalf("Error writing to temporary JSON file: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatalf("Error closing temporary JSON file: %v", err)
		}

		// Load JSON from the temporary file
		result, err := LoadJSON(tmpfile.Name())

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		expected := map[string]interface{}{"name": "John", "age": 30}
		expectedJSON, _ := json.Marshal(expected)
		actualJSON, _ := json.Marshal(result)
		if string(expectedJSON) != string(actualJSON) {
			t.Errorf("Expected result to be %v, but got %v", expected, result)
		}
	})

	t.Run("Invalid JSON File", func(t *testing.T) {
		// Create a temporary JSON file with invalid JSON content
		invalidJSON := `{"name": "John", "age": }`
		tmpfile, err := ioutil.TempFile("", "test.json")
		if err != nil {
			t.Fatalf("Error creating temporary JSON file: %v", err)
		}
		defer os.Remove(tmpfile.Name()) // Clean up

		if _, err := tmpfile.Write([]byte(invalidJSON)); err != nil {
			t.Fatalf("Error writing to temporary JSON file: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatalf("Error closing temporary JSON file: %v", err)
		}

		// Load JSON from the temporary file
		_, err = LoadJSON(tmpfile.Name())
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	t.Run("Non-Existent JSON File", func(t *testing.T) {
		// Attempt to load a non-existent JSON file
		_, err := LoadJSON("nonexistent.json")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})
}
