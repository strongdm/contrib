package cli

import (
	"ext/adapter"
	"ext/util"

	"github.com/urfave/cli"
)

type sdmExt interface {
	adminServersAddAction(ctx *cli.Context) error
}

type sdmExtImpl struct {
	getArgs              func(ctx *cli.Context) cli.Args
	convertStrSliceToStr func(strList []string) string
	checkRegexMatch      func(regexList []string, arguments string) (bool, error)
	getSdmCommand        func(appName, commandName, arguments string) string
	commandNotFound      func(ctx *cli.Context, command string)
	getAppName           func(ctx *cli.Context) string
	getCommandName       func(ctx *cli.Context) string
	mapCommandArguments  func(arguments []string, flags []cli.Flag) map[string]string
	servers              func(commandName string, mappedOptions map[string]string) error
}

func NewSdmExt() *sdmExtImpl {
	return &sdmExtImpl{
		getArgs:              getArgs,
		convertStrSliceToStr: util.ConvertStrSliceToStr,
		checkRegexMatch:      util.CheckRegexMatch,
		getSdmCommand:        getSdmCommand,
		commandNotFound:      commandNotFound,
		getAppName:           getAppName,
		getCommandName:       getCommandName,
		mapCommandArguments:  util.MapCommandArguments,
		servers:              adapter.Servers,
	}
}

func (i *sdmExtImpl) patchGetArgs(getArgs func(ctx *cli.Context) cli.Args) {
	i.getArgs = getArgs
}

func (i *sdmExtImpl) patchConvertStrSliceToStr(convertStrSliceToStr func(strList []string) string) {
	i.convertStrSliceToStr = convertStrSliceToStr
}

func (i *sdmExtImpl) patchCheckRegexMatch(checkRegexMatch func(regexList []string, arguments string) (bool, error)) {
	i.checkRegexMatch = checkRegexMatch
}

func (i *sdmExtImpl) patchGetSdmCommand(getSdmCommand func(appName, commandName, arguments string) string) {
	i.getSdmCommand = getSdmCommand
}

func (i *sdmExtImpl) patchCommandNotFound(commandNotFound func(ctx *cli.Context, command string)) {
	i.commandNotFound = commandNotFound
}

func (i *sdmExtImpl) patchGetAppName(getAppName func(ctx *cli.Context) string) {
	i.getAppName = getAppName
}

func (i *sdmExtImpl) patchGetCommandName(getCommandName func(ctx *cli.Context) string) {
	i.getCommandName = getCommandName
}

func (i *sdmExtImpl) patchMapCommandArguments(mapCommandArguments func(arguments []string, flags []cli.Flag) map[string]string) {
	i.mapCommandArguments = mapCommandArguments
}

func (i *sdmExtImpl) patchServers(servers func(commandName string, mappedOptions map[string]string) error) {
	i.servers = servers
}
