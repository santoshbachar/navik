package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Proxy struct {
	mu    sync.Mutex
	proxy *httputil.ReverseProxy
}

func (p *Proxy) Start(listen string, ports []string) {

	for _, port := range ports {

		uri := "http://localhost:" + port + "/"

		origin, _ := url.Parse(uri)

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		p = p.proxy{Director: director}

		// handleURI := "/" + port + "/"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		p.proxy.ServeHTTP(w, r)
	})

	listeningPort := ":" + listen
	log.Fatalln(http.ListenAndServe(listeningPort, nil))

}
