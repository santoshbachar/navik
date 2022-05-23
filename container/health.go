package container

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type ActuatorHealth struct {
	Status     string      `yaml:"status"`
	Components interface{} `yaml:"components"`
}

func IsUp(host string, port int) bool {

	addr := "http://" + host + ":" + strconv.Itoa(port) + "/actuator/health"

	// client := &http.Client{}
	// req, _ := http.NewRequest("GET", addr, nil)

	// req.Close = true
	// req.Header.Set("Content-Type", "application/json")
	// resp, err := client.Do(req)

	// if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
	// 	color.Red("🚩 port is connection refused/reset. Shall I panic?")
	// 	return false
	// }

	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	color.Blue("addr = " + addr)

	resp, err := http.Get(addr)

	fmt.Println("About to check errors")
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
		color.Red("🚩 port is connection refused/reset. Shall I panic?")
		return false
	}

	// server closes the connection without indicating so
	// resp.Close = true
	color.Cyan("resp.Close")

	if err != nil {
		color.Red("actuator response error")
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			color.Red("Content-Length is 0. Next time the content shall come")
			return false
		} else {
			fmt.Println("Some other error")
			panic(err)
		}
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
	fmt.Println(ah)
	isUp := false
	if ah.Status == "UP" {
		isUp = true
	}
	color.Cyan("ah.status == UP " + strconv.FormatBool(isUp))
	return ah.Status == "UP"
}

func unmarshalActuator(resp []byte) ActuatorHealth {
	var ah ActuatorHealth
	err := yaml.Unmarshal([]byte(resp), &ah)
	if err != nil {
		panic(err)
	}
	return ah
}
