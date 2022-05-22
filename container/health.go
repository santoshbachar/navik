package container

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type ActuatorHealth struct {
	status     string      `yaml:"status"`
	components interface{} `yaml:"components"`
}

func IsUp(host string, port int) bool {
	addr := "http://" + host + ":" + strconv.Itoa(port) + "/actuator/health"
	color.Blue("addr = " + addr)
	resp, err := http.Get(addr)

	fmt.Println("About to check errors")
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
		color.Red("port is connection refused/reset. red flag. Shall I panic?")
		return false
	}

	// server closes the connection without indicating so
	resp.Close = true

	if err != nil {
		color.Red("actuator response error")
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		color.Magenta("status code is " + strconv.Itoa(resp.StatusCode))
		return false
	}

	fmt.Println("About to read the json")

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("About to print json")

	color.Green(string(json))

	ah := unmarshalActuator(json)
	return ah.status == "UP"
}

func unmarshalActuator(resp []byte) ActuatorHealth {
	var ah ActuatorHealth
	err := yaml.Unmarshal([]byte(resp), &ah)
	if err != nil {
		panic(err)
	}
	return ah
}
