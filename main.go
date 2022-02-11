package main

import (
	"context"

	"navik.com/m/v1/dockerDriver"
)

func main() {
	ctx := context.Background()

	cli := dockerDriver.Connect(ctx)

	count := dockerDriver.ListContainers(ctx, cli)
	if count == 0 {
		dockerDriver.StartContainerFromExistingImage(ctx, cli, "hello-world", "")
	}

}
