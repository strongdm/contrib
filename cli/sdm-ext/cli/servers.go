package cli

import (
	"ext/adapter"
	"ext/util"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const (
	FILE_REGEX_PATTERN  = `^--file [\.\/\w,-]+\.[A-Za-z]+$`
	F_REGEX_PATTERN     = `^-f [\.\/\w,-]+\.[A-Za-z]+$`
	STDIN_REGEX_PATTERN = `^--stdin$`
	I_REGEX_PATTERN     = `^-i$`
)

var adminCommand = cli.Command{
	Name:        "admin",
	Usage:       "administrative commands",
	Subcommands: cli.Commands{adminServersCommand},
}

var adminServersCommand = cli.Command{
	Name:        "servers",
	Usage:       "manage servers",
	Subcommands: cli.Commands{adminServersAddCommand},
}

var adminServersAddCommand = cli.Command{
	Name:            "add",
	Aliases:         []string{"create"},
	Usage:           "add one or more server",
	Flags:           adminServersAddFlags,
	Action:          adminServersAddAction,
	SkipFlagParsing: true,
}

var adminServersAddFlags = []cli.Flag{
	util.GetAdminServersAddFileFlag(),
	util.GetAdminServersAddStdinFlag(),
}

func adminServersAddAction(ctx *cli.Context) error {
	argumentList := getArgs(ctx)
	arguments := util.ConvertStrSliceToStr(argumentList)

	matched, err := util.CheckRegexMatch(getRegexList(), arguments)
	if err != nil {
		return err
	}

	if !matched {
		sdmCommand := getSdmCommand(getAppName(ctx), getCommandName(ctx), arguments)
		commandNotFound(ctx, sdmCommand)

		return nil
	}

	mappedArguments := util.MapCommandArguments(argumentList, adminServersAddFlags)
	err = adapter.Servers(ctx.Command.Name, mappedArguments)
	if err != nil {
		return err
	}

	return nil
}

var getArgs = func(ctx *cli.Context) cli.Args {
	if ctx == nil {
		fmt.Println() // Needed because mock
	}
	return ctx.Args()
}

var getSdmCommand = func(appName, commandName, arguments string) string {
	newAppName := removeSdmExt(appName)
	return fmt.Sprintf("%s %s %s", newAppName, commandName, arguments)
}

var getAppName = func(ctx *cli.Context) string {
	if ctx == nil {
		fmt.Println() // Needed because mock
	}
	return ctx.App.Name
}

var getCommandName = func(ctx *cli.Context) string {
	if ctx == nil {
		fmt.Println() // Needed because mock
	}
	return ctx.Command.Name
}

func removeSdmExt(appName string) string {
	return strings.Replace(appName, "sdm-ext ", "", 1)
}

func getRegexList() []string {
	return []string{
		FILE_REGEX_PATTERN,
		F_REGEX_PATTERN,
		STDIN_REGEX_PATTERN,
		I_REGEX_PATTERN,
	}
}
