package boot

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func isPortAvailable(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
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

func replacePortOutFromArgs(args *[]string, newPort int) bool {
	for i, v := range *args {
		firstTwo := v[:2]
		var ports string
		if firstTwo != "-p" {
			continue
		}
		ports = strings.TrimSpace(v[2:])
		fmt.Println("ports =", ports)
		pos := strings.Index(ports, ":")
		fmt.Println("pos = ", pos)
		if pos == -1 {
			fmt.Println(": not found, this shouldn't be happening.")
			return false
		}
		preserve := ports[pos+1:]
		fmt.Println("preserve= ", preserve)
		(*args)[i] = "-p " + strconv.Itoa(newPort) + ":" + preserve
		fmt.Println(ports)
		return true
	}
	return false
}
