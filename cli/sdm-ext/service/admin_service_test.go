package service

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := executeTests{}
	t.Run("Test execute when the passed command is valid",
		tests.testWhenThePassedCommadIsValid)
	t.Run("Test execute when command fails",
		tests.testWhenThePassedCommadFails)
}

type executeTests struct{}

func (tests executeTests) testWhenThePassedCommadIsValid(t *testing.T) {
	sdmServiceImpl := NewSdmServiceMock(runSuccessfulCommandMock)

	commands := getRawTCPServerAddCommand()
	options := getOptionsToExecute()
	postOptions := getRawTCPServerName()

	actualStdout, actualStderr := sdmServiceImpl.execute(commands, options, postOptions)

	expectedStdout := new(strings.Builder)
	expectedStdout.Write([]byte(""))
	expectedStderr := new(strings.Builder)
	expectedStderr.Write([]byte(""))

	assert.Equal(t, expectedStdout, &actualStdout)
	assert.Equal(t, expectedStderr, &actualStderr)
}

func (tests executeTests) testWhenThePassedCommadFails(t *testing.T) {
	sdmServiceImpl := NewSdmServiceMock(runFailedCommandMock)

	commands := getRawTCPServerAddCommand()
	options := getOptionsToExecute()
	postOptions := getRawTCPServerName()

	actualStdout, actualStderr := sdmServiceImpl.execute(commands, options, postOptions)

	expectedStdout := new(strings.Builder)
	expectedStdout.Write([]byte(""))
	expectedStderr := new(strings.Builder)
	expectedStderr.Write([]byte("stderr"))

	assert.Equal(t, expectedStdout, &actualStdout)
	assert.Equal(t, expectedStderr, &actualStderr)
}

func TestOptionsToArguments(t *testing.T) {
	tests := optionsToArgumentsTests{}
	t.Run("Test OptionsToArguments when the conversion is successful",
		tests.testWhenTheConversionIsSuccessful)
	t.Run("Test OptionsToArguments when the options are empty",
		tests.testWhenTheOptionsAreEmpty)
}

type optionsToArgumentsTests struct{}

func (tests optionsToArgumentsTests) testWhenTheConversionIsSuccessful(t *testing.T) {
	optionsMap := getOptionsToExecute()

	actualOptionsList := optionsToArguments(optionsMap)

	expectedOptionsList := getOptionsListToExecute()

	assert.Equal(t, expectedOptionsList, actualOptionsList)
}

func (tests optionsToArgumentsTests) testWhenTheOptionsAreEmpty(t *testing.T) {
	optionsMap := map[string]string{}

	actualOptionsList := optionsToArguments(optionsMap)

	assert.Empty(t, actualOptionsList)
}

func runSuccessfulCommandMock(cmd *exec.Cmd) {
	cmd.Stdout.Write([]byte(""))
	cmd.Stderr.Write([]byte(""))
}

func runFailedCommandMock(cmd *exec.Cmd) {
	cmd.Stdout.Write([]byte(""))
	cmd.Stderr.Write([]byte("stderr"))
}

func getOptionsToExecute() map[string]string {
	return map[string]string{
		"--hostname":      "192.168.101.192",
		"--port":          "49150",
		"--port-override": "13392",
		"--tags":          "key1=value1,key2=value2",
	}
}

func getRawTCPServerAddCommand() string {
	return "admin servers add rawtcp"
}

func getOptionsListToExecute() []string {
	return []string{
		"--hostname",
		"192.168.101.192",
		"--port",
		"49150",
		"--port-override",
		"13392",
		"--tags",
		"key1=value1,key2=value2",
	}
}

func getRawTCPServerName() string {
	return "Example Raw TCP"
}
