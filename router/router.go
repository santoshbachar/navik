package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	_ "os/signal"
	"strconv"
	"sync"

	"github.com/santoshbachar/navik/container"
)

type cPair struct {
	addr string // where to contact the container
	ptr  *container.Info
}

type Config struct {
	portOut int // main port described in -p
	portIn  int // not to be taken into consideration in v1
	// now, this has to be considered as hostport:port
	maintain      int
	containerArgs []string
	routes        []cPair // hoping not to use this technique
	//Routes []string
	mu         sync.Mutex
	stopSignal chan os.Signal
}

var counter int = 0

//func getNewDirector(director func (req *http.Request)) httputil.ReverseProxy{
//	httputil.ReverseProxy{Director: director}
//}

func (cp *cPair) AddAddr(addr int) {
	cp.addr = strconv.Itoa(addr)
}

func (c *Config) GetTotalRoutes() int {
	return len(c.routes)
}

func (c *Config) GetContainerAddr(index int) string {
	fmt.Println("c", c)
	fmt.Println("GetContainerAddr", c.routes[index].addr)
	// return (*c.GetRoutes())[index].addr
	return c.routes[index].addr
}

func (c *Config) GetContainerArgs() *[]string {
	return &c.containerArgs
}

func (c *Config) GetMinimumContainers() int {
	return c.maintain
}

func (c *Config) GetRoutes() *[]cPair {
	return &c.routes
}

func (c *Config) getAddrToContainer(index int) *cPair {
	return &c.routes[index]
}

func (c *Config) getPointerToContainer(index int) *container.Info {
	return c.routes[index].ptr
}

func GetInitialConfig(p1, p2, maintain int, args []string) *Config {
	return &Config{p1, p2, maintain, args, nil, sync.Mutex{}, make(chan os.Signal, 1)}
}

func (c *Config) ModifyRoutes(newRoutes []string) {
	//c.routes = newRoutes
}

func getCurrentCounter() int {
	return counter
}

func getNextCounter(max int) int {
	fmt.Println("Someone is asking for counter")

	if counter >= max {
		counter = 0
	}
	fmt.Println("going to return ", counter)
	counterToReturn := counter
	counter++
	return counterToReturn
}

func getNewDirector(routes *[]cPair) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)

		req.URL.Scheme = "http"
		c := getNextCounter(len(*routes))

		req.URL.Host = "localhost:" + (*routes)[c].addr
		req.Header.Add("X-Origin-Host", req.URL.Host)

		fmt.Println("Inside director function")
		fmt.Println("c=", c)
		fmt.Println("req.URL.Host=", req.URL.Host)
	}
	fmt.Println("Director is set")

	return &httputil.ReverseProxy{Director: director}
}

func (c *Config) Stop() {
	//signal.Notify(c.stopSignal, os.Interrupt)
	c.stopSignal <- os.Interrupt
}

func (conf *Config) Spin(i int, serverMux *http.ServeMux) {
	//director := func(req *http.Request) {
	//	req.URL.Scheme = "http"
	//}

	//serverMux := http.NewServeMux()

	fmt.Println("Spin() i = ", i)
	fmt.Println("conf", conf)

	(*serverMux).HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//proxy := getNewDirector(director)
		conf.mu.Lock()
		proxy := getNewDirector(&conf.routes)
		fmt.Println("/ proxy", proxy)
		fmt.Println("request is ", request)
		conf.mu.Unlock()

		proxy.ServeHTTP(writer, request)
	})

	srv := &http.Server{Addr: ":" + strconv.Itoa(conf.portOut), Handler: serverMux}

	go func() {
		<-conf.stopSignal
		fmt.Println("Server with conf ", conf, "is exiting")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server(proxy) Shutdown: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server(proxy) ListenAndServe: %v", err)
	}

}

func NewRoute() *cPair {
	return &cPair{"", container.NewInfo()}
}

func (c *Config) AddInitialRouteInfo(i int, port int, id string) {
	fmt.Println("Adding initial route info")
	if c.routes == nil {
		//c.routes = append(c.routes, cPair{"", &container.Info{}})
		//c.routes[0].AddAddr(port)
	}
	c.routes = append(c.routes, *NewRoute())
	c.getPointerToContainer(i).AddId(id)
	c.getAddrToContainer(i).AddAddr(port)
	fmt.Println("after adding intitial route info")
	fmt.Println("c", c)
	fmt.Println("c.route", c.routes)
}

func CTest(index int, ok chan bool) {
	fmt.Println("Hello from ", index)
	<-ok
	fmt.Println("Exiting from ", index)
}
