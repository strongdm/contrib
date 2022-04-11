package cli

import (
	"fmt"

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

func commandNotFound(ctx *cli.Context, command string) {
	sdmImpl := NewSdm()
	stdout, stderr := sdmImpl.execute()
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())
}
