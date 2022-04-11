package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type sdm interface {
	execute() (*strings.Builder, *strings.Builder)
}

type sdmImpl struct {
	runCommand func(cmd *exec.Cmd)
}

func NewSdm() *sdmImpl {
	return &sdmImpl{runCommand}
}

func (i sdmImpl) execute() (*strings.Builder, *strings.Builder) {
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	args := os.Args[1:]
	cmd := exec.Command("sdm", args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	i.runCommand(cmd)

	return stdout, stderr
}

func runCommand(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
