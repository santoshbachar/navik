package main

import (
	"navik.com/m/v1/yaml"
)

func main() {
	// ctx := context.Background()

	// cli := dockerDriver.Connect(ctx)

	// count := dockerDriver.ListContainers(ctx, cli)
	// if count == 0 {
	// 	dockerDriver.StartContainerFromExistingImage(ctx, cli, "hello-world", "")
	// }

	yaml.LoadConfig()

}
