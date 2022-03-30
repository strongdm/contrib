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
}

type adminServersAddActionTests struct{}

func (tests adminServersAddActionTests) testWhenThePassedCommandIsValid(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getArgsMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchMock)
	monkey.Patch(util.MapCommandArguments, mapCommandArgumentsMock)
	monkey.Patch(adapter.Servers, serversMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
	monkey.UnpatchAll()
}

func (tests adminServersAddActionTests) testWhenThereIsNoArguments(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(getArgs, getEmptyArgsMock)
	monkey.Patch(util.ConvertStrSliceToStr, convertStrSliceToStrMock)
	monkey.Patch(util.CheckRegexMatch, checkRegexMatchMock)
	monkey.Patch(util.MapCommandArguments, mapCommandArgumentsMock)
	monkey.Patch(adapter.Servers, serversMock)

	actualErr := adminServersAddAction(&cli.Context{})

	assert.Nil(t, actualErr)
	monkey.UnpatchAll()
}

func getArgsMock(ctx *cli.Context) cli.Args {
	return cli.Args{"--file", "file.json"}
}

func getEmptyArgsMock(ctx *cli.Context) cli.Args {
	return cli.Args{}
}

func convertStrSliceToStrMock(strList []string) string {
	return "--file file.json"
}

func checkRegexMatchMock(regexList []string, arguments string) (bool, error) {
	return true, nil
}

func mapCommandArgumentsMock(arguments []string, flags []cli.Flag) map[string]string {
	return map[string]string{"--file": "file.json"}
}

func serversMock(commandName string, mappedOptions map[string]string) error {
	return nil
}

func TestGetSdmCommand(t *testing.T) {
	defer monkey.UnpatchAll()

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
	defer monkey.UnpatchAll()

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
