package yaml

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

var data = `
container: MyContainer
instances:
   min: 1
   max: 3
   cpu: 80
`

type Instance struct {
	Min       int
	Max       int
	Resources int `yaml:"cpu"`
}

type NavikConfig struct {
	Container string
	Instances Instance
}

func LoadConfig() {
	config := NavikConfig{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- config: %v\n\n", config)

	containerName, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Container Name: ", string(containerName))
}
