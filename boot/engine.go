package boot

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	_ "sync"
	"time"

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
			finalArgs := container.PrepareStart("--detach --rm", cArgs)
			id, ok := container.Start(image, instanceName, finalArgs)

			if !ok {
				fmt.Println("Unable to start container. Might handle this in monitoring")
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
