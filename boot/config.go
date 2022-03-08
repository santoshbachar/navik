package boot

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/santoshbachar/navik/constants"
	"github.com/santoshbachar/navik/proxy"
	"github.com/santoshbachar/navik/router"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Driver        string `yaml:"driver"`
	Using         string `yaml:"using"`
	CommonArgs    string `yaml:"common_args"`
	PortPoolRange string `yaml:"port_pool_range"`
	Containers    []struct {
		Image string `yaml:"image"`
		State struct {
			Min int `yaml:"min"`
			Max int `yaml:"max"`
		} `yaml:"state"`
		Args []string `yaml:"args"`
	} `yaml:"containers"`
}

var routers []router.Container

func Bootstrap() {
	var config Config

	dat, err := os.ReadFile(constants.ResourceDir + "navik.containers.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		panic(err)
	}

	pool_min, pool_max, ok := getColonItems(config.PortPoolRange)
	if !ok {
		panic("port_pool_range format error")
	}
	current_port := pool_min

	for _, container := range config.Containers {

		available_port, ok := getNextAvailablePort(current_port, pool_max)
		if !ok {
			panic("not enough port left in the pool to start the container " + container.Image)
		}

		fmt.Println("available port", available_port)

		docker_args := ""
		for _, arg := range container.Args {
			docker_args = docker_args + " " + arg
		}
		docker_args = config.CommonArgs + docker_args
		fmt.Println("docker_args", docker_args)

		r := router.Container{}
		ok = r.Start(container.Image, docker_args, container.State.Min, container.State.Max)
		if !ok {
			panic("some port is occupied for no reason. try again")
		}

		for i := 0; i < container.State.Max; i++ {
			available_port, ok = getNextAvailablePort(current_port, pool_max)
			if !ok {
				panic("not enough port left in the pool to start the container " + container.Image)
			}

			p := proxy.ReverseProxyController{}
			p.Start()
		}

		routers = append(routers, r)
	}

}

func handleArgs(args *[]string) string {
	for i, v := range *args {
		firstTwo := v[:2]
		var ports string
		if firstTwo == "-p" {
			ports = strings.TrimSpace(v[2:])
			(*args)[i] = "-p " + getHostPort(ports)
			fmt.Println(ports)
		}
	}
	fmt.Println(*args)
	return "hey"
}

func getHostPort(ports string) string {
	pos := strings.Index(ports, ":")
	if pos == -1 {
		return ""
	}
	return ports[0:pos]
}

func getColonItems(raw string) (int, int, bool) {
	pos := strings.Index(raw, ":")
	if pos == -1 {
		return 0, 0, false
	}
	f, err := strconv.Atoi(raw[0:pos])
	if err != nil {
		panic("first port format error")
	}
	s, err := strconv.Atoi(raw[pos+1 : len(raw)])
	if err != nil {
		panic("second Port format error")
	}
	return f, s, true

}

func getNextAvailablePort(start_port, end_port int) (int, bool) {
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

func checkPort(port int) bool {
	//cmd := exec.Command("")
	return false
}
