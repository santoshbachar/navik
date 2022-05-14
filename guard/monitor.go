package guard

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/santoshbachar/navik/boot"
	"github.com/santoshbachar/navik/dockerDriver"
)

func GetOverallStatus() string {
	res := "Status:\n"

	names, _ := boot.GetMonitorLists()
	fmt.Println("names", names)

	total := dockerDriver.GetContainerCountWithNewClient(context.Background())
	res += "\nTotal containers: " + strconv.Itoa(total) + "\n"

	list := dockerDriver.ListContainersWithNewClient(context.Background())
	fmt.Println(list)

	for _, name := range names {
		res += name
		ok, id := checkNameInList(name, &list)
		if ok {
			res += " - OK - " + id
		} else {
			res += " - NOT OK"
		}
		res += "\n"
	}

	// for _, name := range names {
	// 	id, err := dockerDriver.SearchContainer(context.Background(), name)
	// 	res += name
	// 	if err {
	// 		res += " - NOT OK"
	// 	} else {
	// 		res += " - OK - " + id
	// 	}
	// 	res += "\n"
	// }

	fmt.Println("return ->", res)

	return res

}

func checkNameInList(instanceName string, list *[]types.Container) (bool, string) {
	for _, container := range *list {
		if checkNameInNames(instanceName, container.Names) {
			return true, container.ID
		}
	}
	return false, ""
}

func checkNameInNames(instanceName string, names []string) bool {
	for _, name := range names {
		if "/"+instanceName == name {
			return true
		}
	}
	return false
}
