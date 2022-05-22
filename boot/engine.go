package boot

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	_ "sync"
	"time"

	"github.com/fatih/color"
	"github.com/santoshbachar/navik/constants"
	"github.com/santoshbachar/navik/container"
	"github.com/santoshbachar/navik/router"
)

type cPair struct {
	assignedPort int
	forwardPort  int // container port
}

// map[string]*router.Config might be more efficient
// would be great if I could have used slice since all the images and containers
// needed is already known in advance, but then we can't give users the ability
// to dynaically reload configs at runtime
var RouterMap map[string]*router.Config
var PortManager router.PortManager
var ContainerPortMonitorList []string
var ContainerNameMonitorList []string

func Start(signal chan os.Signal) {

	RouterMap = make(map[string]*router.Config)

	Boot(&RouterMap, &PortManager)

	spinRouters(&RouterMap)

	fmt.Println("RouterMap after spinRouters ", RouterMap)

	spinContainers(&RouterMap)

	fmt.Println("RouterMap after spinContainers ", RouterMap)

	//for {

	// some wg group magic would be good if at all possible
	// or, I think both these should run as goroutine with for{}
	// go monitorContainers()
	// go monitorPorts()
	//}

	<-signal

	fmt.Println("hello")

}

func GetRouterMap() map[string]*router.Config {
	return RouterMap
}

func GetMonitorLists() ([]string, []string) {
	return ContainerNameMonitorList, ContainerPortMonitorList
}

func monitorContainers() {
	for {
		time.Sleep(1 * time.Second)
		for _, name := range ContainerNameMonitorList {
			fmt.Println("checking for activity, container name", name)
		}
	}
}

func addPortToMonitorList(port int) {
	portToAdd := strconv.Itoa(port)
	ContainerPortMonitorList = append(ContainerPortMonitorList, portToAdd)
}

func monitorPorts() {
	for {
		time.Sleep(5 * time.Second)
		for _, port := range ContainerPortMonitorList {
			fmt.Println("checking for activity, port", port)
		}
	}
}

func spinContainers(routerMap *map[string]*router.Config) {

	fmt.Println("Enter spinContainers pm is ", PortManager)

	for image, c := range *routerMap {
		//c := (*routerMap)[image]
		fmt.Println("c ", c)

		for i := 0; i < c.GetMinimumContainers(); i++ {
			instanceName := image + "-" + strconv.Itoa(i)
			port, ok := PortManager.GetNextAvailablePort()
			if !ok {
				panic("No more ports available to continue, Exiting Navik")
			}
			// cannot be used now.
			// 1. docker run
			// 2. then, add it to router.Config
			// args := c.GetContainerAddr(i)
			cArgs := c.GetContainerArgs()
			ok = replacePortOutFromArgs(cArgs, port)
			if !ok {
				fmt.Println("oh no, port error when about to run conatiner")
				continue
			}
			finalArgs := container.PrepareStart(constants.GetCommonArgs(), cArgs)
			id, ok := container.Start(image, instanceName, finalArgs)

			if !ok {
				fmt.Println("Unable to start container. Might handle this in monitoring")
			}

			color.Green("docker container get started.")

			counter := 0
			for {
				counter++
				color.Yellow("waiting for " + image + " #" + strconv.Itoa(i) + " to start")
				up := container.IsUp(c.GetHost(), port)
				if up {
					break
				}
				time.Sleep(time.Second * 1)
			}

			c.AddInitialRouteInfo(i, port, id)
			// this needs to changed, need map[string]*router.Config
			(*routerMap)[image] = c
			fmt.Println("InitialRouteInfo added", c)
			addPortToMonitorList(port)
			addNameToMonitorList(instanceName)
		}

		fmt.Println("routes -> ", c.GetRoutes())

		//for i := 0; i < c.GetTotalRoutes(); i++ {
		//	port := 2500
		//	id, ok := container.Start(image, "some args")
		//
		//	if !ok {
		//		fmt.Println("Unable to start container. Might handle this in monitoring")
		//	}
		//
		//	fmt.Println("Adding routes for", i)
		//	c.AddInitialRouteInfo(i, port, id)
		//}

		//c.ModifyRoutes()
		fmt.Println("spinContainers", c)
	}
	fmt.Println("Debug")
}

func addNameToMonitorList(name string) {
	ContainerNameMonitorList = append(ContainerNameMonitorList, name)
}

func spinRouters(routerMap *map[string]*router.Config) {
	//mu := sync.Mutex{}
	i := 1
	fmt.Println("spinning len - ", len(RouterMap))
	for _, router := range *routerMap {
		//mu.Lock()
		fmt.Println("spinning for ", i)
		serverMux := http.NewServeMux()
		fmt.Println("router in for loop", router)
		//time.Sleep(5 * time.Second)
		go router.Spin(i, serverMux)
		i++
		//mu.Unlock()
	}

	// in linux, commenting them
	// goalsConf := (*routerMap)["demo"]
	// fmt.Println("goalsConf", goalsConf)
	// goalsConf.Stop()
	// fmt.Println("demoConf stop")

	//c := make(chan bool)
	//for i := 0; i < 3; i++ {
	//	go router.CTest(i, c)
	//}
	//
	//for i := 0; i < 3; i++ {
	//	time.Sleep(5 * time.Second)
	//	c <- true
	//}
}

func AddMaintain(image string, count int) {
	// not required as director does the job once even a single instance is running
	// for i := 0; i < count; i++ {
	// 	serverMux := http.NewServeMux()
	// 	RouterMap[image].Spin(i, serverMux)
	// }
	RouterMap[image].AddMaintain(count)
}

func RemoveMaintain(image string, count int) {
	// not required as director does the job once even a single instance is running
	// for i := 0; i < count; i++ {
	// 	serverMux := http.NewServeMux()
	// 	RouterMap[image].Spin(i, serverMux)
	// }
	RouterMap[image].RemoveMaintain(count)
}

func AddContainer(image string, count int, current int) {
	for i := 0; i < count; i++ {
		num := i + current
		instanceName := image + "-" + strconv.Itoa(num)
		port, ok := PortManager.GetNextAvailablePort()
		if !ok {
			panic("No more ports available to continue, Exiting Navik")
		}
		// cannot be used now.
		// 1. docker run
		// 2. then, add it to router.Config
		// args := c.GetContainerAddr(i)
		c := RouterMap[image]
		cArgs := c.GetContainerArgs()
		ok = replacePortOutFromArgs(cArgs, port)
		if !ok {
			fmt.Println("oh no, port error when about to run conatiner")
			continue
		}
		finalArgs := container.PrepareStart(constants.GetCommonArgs(), cArgs)
		id, ok := container.Start(image, instanceName, finalArgs)

		if !ok {
			fmt.Println("Unable to start container. Might handle this in monitoring")
		}

		c.AddInitialRouteInfo(i, port, id)
		// this needs to changed, need map[string]*router.Config
		// commented for this function, copied from SpinContainers()
		// (*routerMap)[image] = c
		fmt.Println("InitialRouteInfo added", c)
		addPortToMonitorList(port)
		addNameToMonitorList(instanceName)
	}
}

func RemoveContainer(image string, count int) {
	c := RouterMap[image]
	for i := 0; i < count; i++ {
		c.Stop()
	}
}
