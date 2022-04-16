package boot

import (
	"net"
	"strconv"
	"time"
)

func isPortAvailable(port int) bool{
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
	if err != nil {
		return true
	}
	if conn!=nil {
		defer conn.Close()
	}
	return false
}

func getAvailablePortCount(start_port, end_port int) int {
	total := 0
	for port := start_port; port <= end_port; port++ {
		if isPortAvailable(port) {
			total++
		}
	}
	return total
}
