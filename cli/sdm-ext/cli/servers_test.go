package cli

import (
	"ext/adapter"
	"ext/util"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestAdminServersAddAction(t *testing.T) {
	tests := adminServersAddActionTests{}
	t.Run("Test adminServersAddAction when the passed command is valid",
		tests.testWhenThePassedCommandIsValid)
	t.Run("Test adminServersAddAction when there is no arguments",
		tests.testWhenThereIsNoArguments)
	t.Run("Test adminServersAddAction when the passed flag does not exist in sdm-ext cli",
		tests.testWhenThePassedFlagDoesNotExistInSdmExtCli)
	t.Run("Test adminServersAddAction when a subcommand is passed between add command and flag",
		tests.testWhenASubcommandIsPassedBetweenAddCommandAndFlag)
	t.Run("Test adminServersAddAction when a subcommand is passed after flag value",
		tests.testWhenTheSubcommandIsPassedAfterFlagValue)
	t.Run("Test adminServersAddAction when a subcommand is passed after flag",
		tests.testWhenTheSubcommandIsPassedAfterFlag)
}

type adminServersAddActionTests struct{}

func (tests adminServersAddActionTests) testWhenThePassedCommandIsValid(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenMatchesMock)
	monkey.Patch(util.MapCommandArguments, mapCommandArgumentsMock)
	monkey.Patch(adapter.Servers, serversMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func (tests adminServersAddActionTests) testWhenThereIsNoArguments(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getEmptyArgsMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenMatchesMock)
	monkey.Patch(util.MapCommandArguments, mapCommandArgumentsMock)
	monkey.Patch(adapter.Servers, serversMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func (tests adminServersAddActionTests) testWhenThePassedFlagDoesNotExistInSdmExtCli(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsWithWrongFlagMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrWithWrongFlagMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenDoesNotMatchesMock)
	monkey.Patch(getSdmCommand, getSdmCommandMock)
	monkey.Patch(getAppName, getAppNameMock)
	monkey.Patch(getCommandName, getCommandNameMock)
	monkey.Patch(commandNotFound, commandNotFoundMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func (tests adminServersAddActionTests) testWhenASubcommandIsPassedBetweenAddCommandAndFlag(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsWithSubcommandBetweenCommandAndFlagMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrWithSubcommandBetweenCommandAndFlagMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenDoesNotMatchesMock)
	monkey.Patch(getSdmCommand, getSdmCommandMock)
	monkey.Patch(getAppName, getAppNameMock)
	monkey.Patch(getCommandName, getCommandNameMock)
	monkey.Patch(commandNotFound, commandNotFoundMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func (tests adminServersAddActionTests) testWhenTheSubcommandIsPassedAfterFlagValue(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsWithSubcommandAfterFlagValueMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrWithSubcommandAfterFlagValueMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenDoesNotMatchesMock)
	monkey.Patch(getSdmCommand, getSdmCommandMock)
	monkey.Patch(getAppName, getAppNameMock)
	monkey.Patch(getCommandName, getCommandNameMock)
	monkey.Patch(commandNotFound, commandNotFoundMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func (tests adminServersAddActionTests) testWhenTheSubcommandIsPassedAfterFlag(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsWithSubcommandAfterFlagMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrWithSubcommandAfterFlagMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchWhenDoesNotMatchesMock)
	monkey.Patch(getSdmCommand, getSdmCommandMock)
	monkey.Patch(getAppName, getAppNameMock)
	monkey.Patch(getCommandName, getCommandNameMock)
	monkey.Patch(commandNotFound, commandNotFoundMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
}

func TestGetSdmCommand(t *testing.T) {
	tests := getSdmCommandTests{}
	t.Run("Test getSdmCommand when it is successful",
		tests.testWhenItIsSucessful)
}

type getSdmCommandTests struct{}

func (tests getSdmCommandTests) testWhenItIsSucessful(t *testing.T) {
	appName := "sdm-ext admin servers"
	commandName := "add"
	arguments := "--files file.json"
	actualSdmCommand := getSdmCommand(appName, commandName, arguments)

	expectedSdmCommand := "admin servers add --files file.json"

	assert.Equal(t, actualSdmCommand, expectedSdmCommand)
}

func TestRemoveSdmExt(t *testing.T) {
	tests := removeSdmExtTests{}
	t.Run("Test removeSdmExt when it is successful",
		tests.testWhenItIsSucessful)
	t.Run("Test removeSdmExt when it does not contain sdm-ext",
		tests.testWhenItDoesNotContainSdmExt)
}

type removeSdmExtTests struct{}

func (tests removeSdmExtTests) testWhenItIsSucessful(t *testing.T) {
	appName := "sdm-ext admin servers"
	actualNewAppName := removeSdmExt(appName)

	expectedNewAppName := "admin servers"

	assert.Equal(t, actualNewAppName, expectedNewAppName)
}

func (tests removeSdmExtTests) testWhenItDoesNotContainSdmExt(t *testing.T) {
	appName := "sdm admin servers"
	actualNewAppName := removeSdmExt(appName)

	expectedNewAppName := "sdm admin servers"

	assert.Equal(t, expectedNewAppName, actualNewAppName)
}

func getArgsMock(ctx *cli.Context) cli.Args {
	return cli.Args{"--file", "file.json"}
}

func getArgsWithWrongFlagMock(ctx *cli.Context) cli.Args {
	return cli.Args{"--files", "file.json"}
}

func getArgsWithSubcommandBetweenCommandAndFlagMock(ctx *cli.Context) cli.Args {
	return cli.Args{"rdp", "--file", "file.json"}
}

func getArgsWithSubcommandAfterFlagValueMock(ctx *cli.Context) cli.Args {
	return cli.Args{"--file", "file.json", "rdp"}
}

func getArgsWithSubcommandAfterFlagMock(ctx *cli.Context) cli.Args {
	return cli.Args{"--stdin", "rdp"}
}

func getEmptyArgsMock(ctx *cli.Context) cli.Args {
	return cli.Args{}
}

func convertStrSliceToStrMock(strList []string) string {
	return "--file file.json"
}

func convertStrSliceToStrWithSubcommandBetweenCommandAndFlagMock(strList []string) string {
	return "rdp --file file.json"
}

func convertStrSliceToStrWithSubcommandAfterFlagValueMock(strList []string) string {
	return "--file file.json rdp"
}

func convertStrSliceToStrWithSubcommandAfterFlagMock(strList []string) string {
	return "--stdin rdp"
}

func convertStrSliceToStrWithWrongFlagMock(strList []string) string {
	return "--files file.json"
}

func checkRegexMatchWhenMatchesMock(regexList []string, arguments string) (bool, error) {
	return true, nil
}

func checkRegexMatchWhenDoesNotMatchesMock(regexList []string, arguments string) (bool, error) {
	return false, nil
}

func getSdmCommandMock(appName, commandName, arguments string) string {
	return ""
}

func getAppNameMock(ctx *cli.Context) string {
	return "sdm-ext admin servers"
}

func getCommandNameMock(ctx *cli.Context) string {
	return "add"
}

func commandNotFoundMock(ctx *cli.Context, command string) {}

func mapCommandArgumentsMock(arguments []string, flags []cli.Flag) map[string]string {
	return map[string]string{"--file": "file.json"}
}

func serversMock(commandName string, mappedOptions map[string]string) error {
	return nil
}
