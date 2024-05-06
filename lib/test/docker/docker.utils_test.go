package docker

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextPort(t *testing.T) {
	port := 10000
	// Test sequential port allocation
	assert.Equal(t, 10001, nextPort(&port))
	assert.Equal(t, 10002, nextPort(&port))
}

func TestIsOpen(t *testing.T) {
	assert.True(t, isOpen("8888"))

	// Start a local HTTP server for testing and check if the port is open
	go startLocalServer("8081")
	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)
	assert.False(t, isOpen("8081"))
}

func TestGetAvailablePort(t *testing.T) {
	port := GetAvailablePort(8080)
	assert.NotEmpty(t, port)
	assert.True(t, isOpen(port))
}

// Helper function to start a local HTTP server for testing
func startLocalServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
	http.ListenAndServe(":"+port, nil)
}
