package container

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/santoshbachar/navik/bash"
	"github.com/santoshbachar/navik/dockerDriver"
	"github.com/santoshbachar/navik/proxy"
)

type ContainerOld struct {
	Image     string
	Min       int
	Max       int
	Instances int
	ProxyPool []proxy.ProxyPool
}

type Info struct {
	ID                string
	live              bool
	lastKnownLiveTime time.Time
}

func NewInfo() *Info {
	return &Info{}
}

func RegisterLive(i *Info) {
	i.live = true
	i.lastKnownLiveTime = time.Now()
}

func (i *Info) AddId(id string) {
	i.ID = id
	RegisterLive(i)
}

func PrepareStart(argsN string, argsAP *[]string) string {
	result := argsN
	for _, arg := range *argsAP {
		result = result + " " + arg
	}
	return result
}

func Start(image string, instanceName string, args string) (string, bool) {
	finalArgs := "run " + args + " --name " + instanceName + " " + image
	var argsSlice = strings.Fields(finalArgs)
	_, err := bash.Command("docker", argsSlice)
	if err != nil {
		fmt.Println("docker run error")
		return "", false
	}

	//id, ok := dockerDriver.MockSearchContainer(name)
	//// might be a good idea to pass this on first go
	//if !ok {
	//	fmt.Println("something is wrong with docker.")
	//	return "", false
	//}
	//id, ok := getContainerIdByName(name)
	id, ok := dockerDriver.SearchContainer(context.Background(), image)
	if !ok {
		fmt.Println("Container", image, "started but unable to detect afterwards, returning.")
		return "", false
	}
	//return id, true
	return id, true
}

func getContainerIdByName(name string) (string, bool) {
	args := []string{
		"ps",
		"|",
		"grep",
		name,
	}
	_, err := bash.Command("docker", args)
	if err != nil {
		fmt.Println("Panic: getContainerIdByName() failed to get Name")
		return "", false
	}
	return "something temprary for now", true
}

func (c *ContainerOld) Start(name string, args string, min, max int) bool {
	for i := 0; i < max; i++ {
		finalArgs := "run " + args + " " + name + " --name " + name + "-" + strconv.Itoa(i+1)
		var argsSlice = strings.Fields(finalArgs)
		fmt.Println("Commented code -> docker", argsSlice)
		//_, err := bash.Command("docker", argsSlice)
		// _, err := bash.Command("ls", "-la")
		// _, err := bash.Command("docker", "run")

		//if err != nil {
		//	fmt.Println("err ->", err)
		//	return false
		//}
		return false //forcefully for test
	}

	return true
}

//func (ph *ProxyHandler) AddProxies(reverseProxy ReverseProxy) {
//	ph.ReverseProxies = append(ph.ReverseProxies, reverseProxy)
//}
//
//func (c *Container) AddContainer(name string, reverseProxy ReverseProxy) {
//	c.Name = name
//	ph := ProxyHandler{}
//	ph.AddProxies(reverseProxy)
//	c.ProxyHolder = append(c.ProxyHolder, ph)
//}

//func (c *Container) Shutdown() {
//	for _, p := range c.ProxyPool {
//		p.Shutdown()
//	}
//}
