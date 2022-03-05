package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"os"
)

type Bearing struct {
	Addr string
	Port string
}
type Proxy struct {
	Bearing       Bearing
	ListeningPort string
	mu            sync.Mutex
	ReverseProxy  *httputil.ReverseProxy
}

func (p *Proxy) Start(instance int, signal chan os.Signal) {

	uri := "http://localhost:" + p.Bearing.Port + "/"

	origin, _ := url.Parse(uri)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	p.mu.Lock()
	p.ReverseProxy = &httputil.ReverseProxy{Director: director}
	p.mu.Unlock()

	if p.ReverseProxy == nil {
		panic("Nil Reverse Proxy")
	}

	handleURI := "/" + string(instance) + "/"
	http.HandleFunc(handleURI, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		p.ReverseProxy.ServeHTTP(w, r)
	})

	listeningPort := ":" + p.ListeningPort
	go log.Fatalln(http.ListenAndServe(listeningPort, nil))
	<- signal
	fmt.Println("line after ListenAndServe")
}

func (p *Proxy) Check() bool {

	return false
}

func Shutdown() {
	req := http.Request{}
	req.Header.Add("NAVIK-SIGNAL", "0")
	req.URL.Scheme = "http"
}
