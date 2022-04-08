package service

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type sdmService interface {
	execute(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder)
}

type sdmServiceImpl struct {
	runCommand func(cmd *exec.Cmd)
}

func NewSdmService() *sdmServiceImpl {
	return &sdmServiceImpl{runCommand}
}

func (i sdmServiceImpl) execute(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder) {
	opts := append(optionsToArguments(options), postOptions)
	commandsAndOptions := append(strings.Split(commands, " "), opts...)

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	cmd := exec.Command("sdm", commandsAndOptions...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	i.runCommand(cmd)

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

	keys := make([]string, 0, len(options))
	for key := range options {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(key, options[key])
		if key[len(key)-1:] == "=" {
			key += options[key]
			options[key] = ""
		}
		strOptions = append(strOptions, key)
		if len(options[key]) > 0 {
			strOptions = append(strOptions, options[key])
		}
	}

	return strOptions
}
