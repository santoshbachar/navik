package boot

import (
	"fmt"
	//"github.com/santoshbachar/navik/container"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/santoshbachar/navik/constants"
	//"github.com/santoshbachar/navik/proxy"
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

//var routers []container.Container

func PreFlightCheck(config *Config) (map[string]router.Config, bool) {
	// might need to add check for docker avialability later
	requiredPorts := 0
	routeMap := make(map[string]router.Config)
	for _, container := range config.Containers {
		ok, p1, p2 := getPortsFromArgs(&container.Args)
		if !ok {
			fmt.Println("your image `" + container.Image + "` is missing ports")
			return nil, false
		}
		fmt.Println(p1, p2)
		portok := isPortAvailable(p1)
		if !portok {
			fmt.Println("your image with port `" + strconv.Itoa(p1) + "` is not available")
			return nil, false
		}
		routeMap[container.Image] = router.GetInitialConfig(p1, p2, container.State.Min, container.Args)
		//routeMap[container.Image] = router.Config{p1, p2, container.State.Min}
		requiredPorts += container.State.Min
	}
	fmt.Println("Required ports by user ", requiredPorts)
	fmt.Println("Required ports by system ", requiredPorts*2)

	p1, p2, ok := getColonItems(config.PortPoolRange)
	if !ok {
		fmt.Println("invalid port-pool")
		return nil, false
	}

	actual := getAvailablePortCount(p1, p2)
	fmt.Println("Available ports in system ", actual)
	if actual < requiredPorts*2 {
		fmt.Println("There's not enough ports to continue")
		return nil, false
	}

	return routeMap, true
}

func Boot(routeMap *map[string]router.Config, portManager *router.PortManager) {
	var config Config

	*routeMap = Bootstrap(&config)
	fmt.Println("Bootstrap is done. Unpacking")

	LoadPortManager(portManager, &config)

	LoadRouterMap(routeMap, &config)
}

func LoadPortManager(pm *router.PortManager, config *Config) {
	min, max, ok := getColonItems(config.PortPoolRange)
	if !ok {
		panic("port-range invalid format")
	}
	pm.InitializePortManager(min, max)
	fmt.Println("Port Manager initialized", pm)
}

func LoadRouterMap(routeMap *map[string]router.Config, config *Config) {
	for _, container := range config.Containers {
		//p1, p2 :=
		//routeMap[container.Image] =
		i := container.Args
		fmt.Println(i)
	}
}

func Bootstrap(config *Config) map[string]router.Config {

	dat, err := os.ReadFile(constants.ResourceDir + "navik.containers.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		panic(err)
	}

	routeMap, ok := PreFlightCheck(config)
	if !ok {
		panic("Pre Flight Check failed. Check stacktrace for more info.")
	}

	//for _, container := range config.Containers {
	//	PreFlightCheck(&container.Args)
	//}

	return routeMap

}

//func BootStrapAdvanced(config *Config) {
//
//	pool_min, pool_max, ok := getColonItems(config.PortPoolRange)
//	if !ok {
//		panic("port_pool_range format error")
//	}
//	current_port := pool_min
//
//	for _, container := range config.Containers {
//
//		available_port, ok := getNextAvailablePort(current_port, pool_max)
//		if !ok {
//			panic("not enough port left in the pool to start the container " + container.Image)
//		}
//
//		fmt.Println("available port", available_port)
//
//		docker_args := ""
//		for _, arg := range container.Args {
//			docker_args = docker_args + " " + arg
//		}
//		docker_args = config.CommonArgs + docker_args
//		fmt.Println("docker_args", docker_args)
//
//		r := container.Container{}
//		ok = r.Start(container.Image, docker_args, container.State.Min, container.State.Max)
//		if !ok {
//			panic("some port is occupied for no reason. try again")
//		}
//
//		for i := 0; i < container.State.Max; i++ {
//			available_port, ok = getNextAvailablePort(current_port, pool_max)
//			if !ok {
//				panic("not enough port left in the pool to start the container " + container.Image)
//			}
//
//			p := proxy.ReverseProxyController{}
//			p.Start()
//		}
//
//		routers = append(routers, r)
//	}
//
//}

func getPortsFromArgsIfArgsWereOne(args string) (bool, int, int) {
	pos := strings.Index(args, "-p ")
	if pos == -1 {
		return false, 0, 0
	}
	newString := args[pos+3:]
	newPos := strings.Index(newString, " ")
	if pos == -1 {
		return false, 0, 0
	}
	found := getHostPort(args[pos+3 : newPos])
	one, two, ok := getColonItems(found)
	if !ok {
		panic("error")
	}
	fmt.Println(one, two)
	return true, one, two
}

func getPortsFromArgsIfArgsNeedToBeReplaced(args *[]string) (bool, int, int) {
	for i, v := range *args {
		firstTwo := v[:2]
		var ports string
		if firstTwo == "-p" {
			ports = strings.TrimSpace(v[2:])
			(*args)[i] = "-p " + getHostPort(ports)
			fmt.Println(ports)
		}
	}
	return false, 0, 0
}

func getPortsFromArgs(args *[]string) (bool, int, int) {
	for _, v := range *args {
		firstTwo := v[:2]
		var ports string
		if firstTwo == "-p" {
			ports = strings.TrimSpace(v[2:])
			//(*args)[i] = "-p " + getHostPort(ports)
			//fmt.Println(ports)
			p1, p2, ok := getColonItems(ports)
			if !ok {
				panic("something wrong port colons")
			}
			return true, p1, p2
		}
	}
	return false, 0, 0
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
