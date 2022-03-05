package containers

import (
	"fmt"
	"github.com/santoshbachar/navik/bash"
)

type Container struct {
	Name      string
	Min       int
	Max       int
	Instances int
	//ProxyPool []proxy.ProxyPool
}

func (c *Container) Start(name string, args string, min, max int) bool {
	//_, err := bash.Command("docker", "run "+args+" "+name)
	_, err := bash.Command("ls", "")

	if err != nil {
		fmt.Println("err ->", err)
		return false
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
