package utility

import (
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func GetPortStringFromPortInt(port int) string {
	return strconv.Itoa(port)
}
func GetPortIntFromPortString(port string) int {
	p_i, err := strconv.Atoi(port)
	if err != nil {
		color.Red("strconv error for port " + port)
		color.Yellow("Tip: Expecting port as number only, don't pass the complete port/protocol")
		panic(err)
	}
	return p_i
}

func GetNumberAndProtocolFromPort(port string) (int, string) {
	num := strings.Index(port, "/")
	portInt := 0
	if num == -1 {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			color.Red(port + " is not a number. conversion failure")
			panic(err)
		}
		return portInt, ""
	}

	return portInt, port[num:]
}
