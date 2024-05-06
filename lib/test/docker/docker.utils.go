package docker

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var lock sync.Mutex

func nextPort(port *int) int {
	lock.Lock()
	*port = *port + 1

	defer lock.Unlock()
	return *port
}

func isOpen(port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("", port), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
		defer conn.Close()
		return false
	}
	return false
}

func GetAvailablePort(defaultPort int) string {
	if isOpen(fmt.Sprint(defaultPort)) {
		return fmt.Sprint(defaultPort)
	}
	for {
		p := fmt.Sprintf("%d", nextPort(&defaultPort))
		if isOpen(p) {
			return p
		}
	}
}
