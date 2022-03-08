package proxy

import (
	"os"
	"os/signal"

	"github.com/docker/go-connections/proxy"
)

type ReverseProxyController struct {
	Proxy  proxy.Proxy
	Signal chan os.Signal
}

func (rpc *ReverseProxyController) Shutdown() {
	signal.Notify(rpc.Signal, os.Kill)
}

func (rpc *ReverseProxyController) Start() {

}

type ProxyPool struct {
	ReverseProxies []ReverseProxyController
}

func (pp *ProxyPool) Shutdown() {
	for _, rp := range pp.ReverseProxies {
		rp.Shutdown()
	}
}

//func (ph *ProxyHandler) AddProxies(reverseProxy ReverseProxy) {
//	ph.ReverseProxies = append(ph.ReverseProxies, reverseProxy)
//}
//
//func (c *Containers) AddContainer(name string, reverseProxy ReverseProxy) {
//	c.Name = name
//	ph := ProxyHandler{}
//	ph.AddProxies(reverseProxy)
//	c.ProxyHolder = append(c.ProxyHolder, ph)
//}
//
//func (c *Containers) AddProxyWithContainer() {
//
//}
//
//func (c *Container) Shutdown() {
//	signal.Notify(c.Signal, os.Kill)
//}
