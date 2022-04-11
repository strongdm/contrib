package cli

import (
	"ext/util"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const (
	FILE_REGEX_PATTERN  = `^--file [\.:\/\\\w,-]+\.[A-Za-z]+$`
	F_REGEX_PATTERN     = `^-f [\.:\/\\\w,-]+\.[A-Za-z]+$`
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
	Action:          NewSdmExt().adminServersAddAction,
	SkipFlagParsing: true,
}

var adminServersAddFlags = []cli.Flag{
	util.GetAdminServersAddFileFlag(),
	util.GetAdminServersAddStdinFlag(),
}

func (i sdmExtImpl) adminServersAddAction(ctx *cli.Context) error {
	argumentList := i.getArgs(ctx)
	arguments := i.convertStrSliceToStr(argumentList)
	matched, err := i.checkRegexMatch(getRegexList(), arguments)
	if err != nil {
		return err
	}

	if !matched {
		sdmCommand := i.getSdmCommand(i.getAppName(ctx), i.getCommandName(ctx), arguments)
		i.commandNotFound(ctx, sdmCommand)

		return nil
	}

	mappedArguments := i.mapCommandArguments(argumentList, adminServersAddFlags)
	err = i.servers(ctx.Command.Name, mappedArguments)
	if err != nil {
		return err
	}

	return nil
}

func getArgs(ctx *cli.Context) cli.Args {
	return ctx.Args()
}

func getSdmCommand(appName, commandName, arguments string) string {
	newAppName := removeSdmExt(appName)
	return fmt.Sprintf("%s %s %s", newAppName, commandName, arguments)
}

func getAppName(ctx *cli.Context) string {
	return ctx.App.Name
}

func getCommandName(ctx *cli.Context) string {
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
