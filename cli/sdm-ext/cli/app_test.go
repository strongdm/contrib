package cli

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteWithSdm(t *testing.T) {
	tests := executeWithSdmTests{}
	t.Run("TestExecuteWithSdm when command is successful",
		tests.testWhenSdmCommandIsSuccessful)
	t.Run("TestExecuteWithSdm when command fails",
		tests.testWhenSdmCommandFails)
}

type executeWithSdmTests struct{}

func (tests executeWithSdmTests) testWhenSdmCommandIsSuccessful(t *testing.T) {
	runCommandBackup := runCommand
	runCommand = runSuccessfulCommandMock
	defer func() {
		runCommand = runCommandBackup
	}()

	actualStdout, actualStderr := executeWithSdm()

	expectedStdout := new(strings.Builder)
	expectedStdout.Write(getFilledStdout())
	expectedStderr := new(strings.Builder)
	expectedStderr.Write(getEmptyStderr())

	assert.Equal(t, expectedStdout, actualStdout)
	assert.Equal(t, expectedStderr, actualStderr)
}

func (tests executeWithSdmTests) testWhenSdmCommandFails(t *testing.T) {
	runCommandBackup := runCommand
	runCommand = runFailedCommandMock
	defer func() {
		runCommand = runCommandBackup
	}()

	actualStdout, actualStderr := executeWithSdm()

	expectedStdout := new(strings.Builder)
	expectedStdout.Write(getEmptyStdout())
	expectedStderr := new(strings.Builder)
	expectedStderr.Write(getFilledStderr())

	assert.Equal(t, expectedStdout, actualStdout)
	assert.Equal(t, expectedStderr, actualStderr)
}

func runSuccessfulCommandMock(cmd *exec.Cmd) {
	cmd.Stdout.Write(getFilledStdout())
	cmd.Stderr.Write(getEmptyStderr())
}

func runFailedCommandMock(cmd *exec.Cmd) {
	cmd.Stdout.Write(getEmptyStdout())
	cmd.Stderr.Write(getFilledStderr())
}

func getFilledStdout() []byte {
	return []byte("stdout")
}

func getEmptyStdout() []byte {
	return []byte("")
}

func getFilledStderr() []byte {
	return []byte("stderr")
}

func getEmptyStderr() []byte {
	return []byte("")
}
