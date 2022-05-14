package main

import (
	"os"

	"github.com/santoshbachar/navik/boot"
	"github.com/santoshbachar/navik/ear"
	//"github.com/santoshbachar/navik/infrastructure"
	//"github.com/santoshbachar/navik/proxy"
	//"github.com/santoshbachar/navik/containers"
)

func main() {
	//ctx := context.Background()
	//
	//cli := dockerDriver.Connect(ctx)
	//
	//count := dockerDriver.ListContainers(ctx, cli)
	//if count == 0 {
	//	dockerDriver.StartContainerFromExistingImage(ctx, cli, "hello-world", "")
	//}

	// yaml.LoadConfig()

	//boot.Bootstrap()

	go ear.Eavesdrop()

	boot.Start(make(chan os.Signal))

	return
	//infrastructure.Provision()
	//
	//var forwardPorts [2]string
	//forwardPorts[0] = "9001"
	//forwardPorts[1] = "9002"
	//
	//var listeningPorts [2]string
	//listeningPorts[0] = "2002"
	//listeningPorts[1] = "2003"
	//
	//var proxies []proxy.Proxy
	//var containers containers.Containers
	//
	//for k, listen := range listeningPorts {
	//
	//	done := make(chan os.Signal, 1)
	//	proxy := proxy.Proxy{Bearing: proxy.Bearing{Addr: "localhost", Port: forwardPorts[k]}, ListeningPort: listen}
	//	containers.Add(containers.Con{proxy, done})
	//	containers.AddContainer(containers.C)
	//	proxies = append(proxies, proxy)
	//	fmt.Println("Start on", proxy.Bearing)
	//
	//	go proxy.Start(k, nil)
	//
	//}
	//
	//containers[1].Shutdown()
	//
	//for {
	//	for k, proxy := range proxies {
	//		time.Sleep(time.Second * 5)
	//		fmt.Println("#", k, proxy.Check())
	//	}
	//
	//}

	// fmt.Scanln()

}
