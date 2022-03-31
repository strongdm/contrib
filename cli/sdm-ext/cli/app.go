package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// Version is set by the build system.
var Version = ""

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "sdm-ext"
	app.Usage = "sdm-ext is an extension of sdm admin"
	app.Version = Version
	if app.Version == "" {
		app.HideVersion = true
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, adminCommand)
	app.CommandNotFound = commandNotFound

	return app
}

func executeWithSdm() (*strings.Builder, *strings.Builder) {
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	args := os.Args[1:]
	cmd := exec.Command("sdm", args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runCommand(cmd)

	return stdout, stderr
}

func runCommand(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func commandNotFound(ctx *cli.Context, command string) {
	stdout, stderr := executeWithSdm()
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())
}
