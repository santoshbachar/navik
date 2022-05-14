package ear

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/santoshbachar/navik/boot"
	"github.com/santoshbachar/navik/guard"
)

func Eavesdrop() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "reload in progress")
		resp := boot.AuxBootstrap(boot.GetRouterMap())
		fmt.Fprintln(w, resp)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Status queried from Navik-Cli")
		resp := guard.GetOverallStatus()
		fmt.Println("resp is prepared in ", time.Now())
		fmt.Fprintln(w, resp)
		fmt.Println("file is writte in", time.Now())
	})

	log.Fatal(http.ListenAndServe(":2000", nil))

}
