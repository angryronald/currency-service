package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestResponse(t *testing.T) {
	// Prepare a sample data for the response
	data := map[string]interface{}{
		"key": "value",
	}

	// Create a new HTTP request with a response recorder
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Call the Response function
	Response(w, http.StatusOK, data, logrus.New())

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
	}

	// Decode the response body
	var response ResponseModel
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response data
	if !reflect.DeepEqual(response.Data, data) {
		t.Errorf("expected %v, got %s", response.Data, data)
	}
}

func TestResponseError(t *testing.T) {
	// Create a new HTTP request with a response recorder
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Call the ResponseError function
	ResponseError(w, http.StatusNotFound, "", logrus.New())

	// Check the response status code
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code %d but got %d", http.StatusNotFound, w.Code)
	}

	// Decode the response body
	var response ResponseModel
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response message
	if response.Message != "" {
		t.Errorf("expected response message 'Not Found' but got %s", response.Message)
	}
}
