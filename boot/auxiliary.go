package boot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/santoshbachar/navik/constants"
	"github.com/santoshbachar/navik/router"
	"gopkg.in/yaml.v2"
)

func AuxBootstrap(routerMap map[string]*router.Config) bool {
	dat, err := os.ReadFile(constants.ResourceDir + "navik.containers.yaml")
	if err != nil {
		panic(err)
	}

	var config *Config
	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		panic(err)
	}

	for _, proposedContainer := range config.Containers {
		proposedMin := proposedContainer.State.Min
		currentMin := (routerMap)[proposedContainer.Image].GetMinimumContainers()
		if proposedMin == currentMin {
			fmt.Println(proposedContainer.Image + " is intact")
		} else {
			if proposedMin > currentMin {
				add(proposedContainer.Image, proposedMin-currentMin, currentMin)
			} else {
				reduce(proposedContainer.Image, currentMin-proposedMin)
			}
		}
	}
	return true
}

func add(image string, count int, current int) {
	fmt.Println("yaml changed, adding " + strconv.Itoa(count) + " new routers, containers")
	AddMaintain(image, count)
	AddContainer(image, count, current)
}

func reduce(image string, count int) {
	fmt.Println("yaml changed, reducing " + strconv.Itoa(count) + " existing routers, containers")
	RemoveMaintain(image, count)
	RemoveContainer(image, count)
}
