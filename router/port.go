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
		fmt.Println("Trying out for port", port)
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			fmt.Print("Error in listening to port ", port, " continuing")
			continue
		}
		err = ln.Close()

		if err != nil {
			fmt.Println("Port", port, "is unavailable/taken")
			continue
		}
		return port, true
	}
	return 0, false
}

func GetNextAvailablePort2(start_port, end_port int) (int, bool) {
	for port := start_port; port <= end_port; port++ {
		fmt.Println("Trying out for port", port)
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			fmt.Println("TImeout error conneting", port, "ignoring and continuing...")
			continue
		}
		if err != nil {
			fmt.Println("port error", err)
			continue
		}
		if conn != nil {
			defer conn.Close()
			continue
		}
		return port, true
	}
	return 0, false
}
