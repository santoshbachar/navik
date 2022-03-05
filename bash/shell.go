package bash

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Command(app string, args string) ([]byte, error) {
	cmd := exec.Command(app, args)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	fmt.Println("cmd to exec", cmd)
	//return cmd.Output()
	err := cmd.Run()
	if err != nil {
		panic(stderr.String())
	}
	fmt.Println("Cmd output", cmd.String())
	var by []byte
	return by, nil
}
