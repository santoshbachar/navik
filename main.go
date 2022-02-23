package main

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
	// listeningPorts[1] = "1003"

	for _, listen := range listeningPorts {

		proxy.Start(listen, forwardPorts[:])

	}

}
