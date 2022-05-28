package boot

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/santoshbachar/navik/utility"
)

func isPortAvailable(portWithProtocol string) bool {
	pNum, protocol := utility.GetNumberAndProtocolFromPort(portWithProtocol)
	timeout := time.Second
	conn, err := net.DialTimeout(protocol, net.JoinHostPort("localhost", utility.GetPortStringFromPortInt(pNum)), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
		defer conn.Close()
	}
	return false
}

func getAvailablePortCount(start_port, end_port string) int {

	start_port_i := utility.GetPortIntFromPortString(start_port)
	end_port_i := utility.GetPortIntFromPortString(end_port)

	total := 0

	for port := start_port_i; port <= end_port_i; port++ {

		port_s_wp := utility.GetPortStringFromPortInt(port) + "/tcp"

		if isPortAvailable(port_s_wp) {
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
