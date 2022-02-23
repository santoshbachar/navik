package main

import (
	"fmt"
	"github.com/santoshbachar/navik/proxy"
)

func main() {
	// ctx := context.Background()

	// cli := dockerDriver.Connect(ctx)

	// count := dockerDriver.ListContainers(ctx, cli)
	// if count == 0 {
	// 	dockerDriver.StartContainerFromExistingImage(ctx, cli, "hello-world", "")
	// }

	// yaml.LoadConfig()

	var forwardPorts [2]string
	forwardPorts[0] = "9001"
	forwardPorts[1] = "9002"

	var listeningPorts [2]string
	listeningPorts[0] = "2002"
	listeningPorts[1] = "2003"

	var proxies []proxy.Proxy

	for k, listen := range listeningPorts {

		fmt.Println("enter listeningPorts on", listen)
		proxy := proxy.Proxy{Bearing: proxy.Bearing{Addr: "localhost", Port: forwardPorts[k]}, ListeningPort: listen}
		proxies = append(proxies, proxy)
		fmt.Println("Start on", proxy.Bearing)

		go proxy.Start(k)

	}

	fmt.Scanln()

}
