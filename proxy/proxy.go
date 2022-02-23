package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
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

func (p *Proxy) Start(instance int) {

	uri := "http://localhost:" + p.Bearing.Port + "/"

	origin, _ := url.Parse(uri)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	p.ReverseProxy = &httputil.ReverseProxy{Director: director}

	handleURI := "/" + string(instance) + "/"
	http.HandleFunc(handleURI, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		p.ReverseProxy.ServeHTTP(w, r)
	})

	listeningPort := ":" + p.ListeningPort
	log.Fatalln(http.ListenAndServe(listeningPort, nil))
}
