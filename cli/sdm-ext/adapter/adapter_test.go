package adapter

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestServers(t *testing.T) {
	defer monkey.UnpatchAll()

	tests := serversTests{}
	t.Run("Test Servers when command name is add",
		tests.testWhenCommandNameIsAdd)
	t.Run("Test Servers when the passed commad name is unknown",
		tests.testWhenThePassedCommandNameIsUnknown)
}

type serversTests struct{}

func (tests serversTests) testWhenCommandNameIsAdd(t *testing.T) {
	commandName := "add"
	mappedOptions := map[string]string{}

	actualErr := Servers(commandName, mappedOptions)

	assert.Nil(t, actualErr)
}

func (tests serversTests) testWhenThePassedCommandNameIsUnknown(t *testing.T) {
	commandName := "adx"
	mappedOptions := map[string]string{}

	actualErr := Servers(commandName, mappedOptions)

	expectedErr := errors.New("unknown command name: adx")

	assert.Equal(t, expectedErr, actualErr)
}
