package service

import (
	"fmt"
	"os/exec"
	"strings"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func execute(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder) {
	opts := append(optionsToArguments(options), postOptions)
	commandsAndOptions := append(strings.Split(commands, " "), opts...)

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	cmd := exec.Command("sdm", commandsAndOptions...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runCommand(cmd)

	return *stdout, *stderr
}

func runCommand(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func optionsToArguments(options map[string]string) []string {
	strOptions := []string{}

	for key, value := range options {
		if key[len(key)-1:] == "=" {
			key += value
			value = ""
		}
		strOptions = append(strOptions, key)
		if len(value) > 0 {
			strOptions = append(strOptions, value)
		}
	}

	return strOptions
}
