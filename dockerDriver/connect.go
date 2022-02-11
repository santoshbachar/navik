package dockerDriver

import (
	"context"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Connect(ctx context.Context) *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli
}

func ListContainers(ctx context.Context, cli *client.Client) int {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return len(containers)
}

func StartContainerFromExistingImage(ctx context.Context, cli *client.Client, imageName string, containerName string) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func StartContainerWithConfigFromExistingImage(ctx context.Context, cli *client.Client, imageName string, containerName string) {
	volName := "/"
	vol := map[string]struct{}{volName: {}}
	config := &container.Config{
		Volumes: vol,
	}
	_, err := cli.ContainerCreate(ctx, config, nil, nil, nil, containerName)

	if err != nil {
		panic(err)
	}

	query := url.Values{}
	query.Set("name", "somename")
}
