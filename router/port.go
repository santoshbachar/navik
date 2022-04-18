package router

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type PortManager struct {
	min     int
	max     int
	current int
}

func (pm *PortManager) InitializePortManager(min, max int) {
	pm.min = min
	pm.max = max
	pm.current = pm.min
}

func (pm *PortManager) UpdateCurrentPortByOne(port int) {
	updateTo := port + 1
	if updateTo > pm.max {
		return
	}
	pm.current = updateTo
}

func (pm *PortManager) GetNextAvailablePort() (int, bool) {
	availablePort, ok := GetNextAvailablePort(pm.current, pm.max)
	pm.UpdateCurrentPortByOne(availablePort)
	return availablePort, ok
}

func GetNextAvailablePort(start_port, end_port int) (int, bool) {
	for port := start_port; port <= end_port; port++ {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
		if err != nil {
			fmt.Println("Error conneting", port)
		}
		if conn != nil {
			defer conn.Close()
			continue
		}
		return port, true
	}
	return 0, false
}
