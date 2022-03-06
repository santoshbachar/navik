package bash

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Command(app string, args []string) ([]byte, error) {
	var cmd *exec.Cmd
	if args == nil {
		cmd = exec.Command(app)
	} else {
		cmd = exec.Command(app, args...)
	}
	// cmd = exec.Command("docker", "run", "--detach", "--rm", "-p", "9001:8080", "demo")
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
	fmt.Println("Cmd string", cmd.String())
	fmt.Println("out output", out.String())
	var by []byte
	return by, nil
}
