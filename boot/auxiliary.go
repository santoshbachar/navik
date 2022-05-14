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
				add(proposedMin - currentMin)
			} else {
				reduce(currentMin - proposedMin)
			}
		}
	}
	return true
}

func add(count int) {
	fmt.Println("yaml changed, adding " + strconv.Itoa(count) + " new routers, containers")
}

func reduce(count int) {
	fmt.Println("yaml changed, reducing " + strconv.Itoa(count) + " existing routers, containers")
}
