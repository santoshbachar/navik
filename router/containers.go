package router

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/santoshbachar/navik/bash"
	"github.com/santoshbachar/navik/proxy"
)

type Container struct {
	Image     string
	Min       int
	Max       int
	Instances int
	ProxyPool []proxy.ProxyPool
}

func (c *Container) Start(name string, args string, min, max int) bool {
	for i := 0; i < max; i++ {
		finalArgs := "run " + args + " " + name + " --name " + name + "-" + strconv.Itoa(i+1)
		var argsSlice = strings.Fields(finalArgs)
		_, err := bash.Command("docker", argsSlice)
		// _, err := bash.Command("ls", "-la")
		// _, err := bash.Command("docker", "run")

		if err != nil {
			fmt.Println("err ->", err)
			return false
		}
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
