package drudapi

import (
	"fmt"
	"os"
	"os/exec"
)

// DockerCompose serves as a wrapper to docker-compose
func DockerCompose(arg ...string) {
	proc := exec.Command("docker-compose", arg...)
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	err := proc.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
