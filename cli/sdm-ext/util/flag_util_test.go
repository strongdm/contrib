package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFlag(t *testing.T) {
	tests := findFlagTests{}
	t.Run("TestFindFlag when find flag is successful",
		tests.testWhenFindFlagIsSuccessful)
	t.Run("TestFindFlag when the given options are not found",
		tests.testWhenTheGivenOptionsAreNotFound)
	t.Run("TestFindFlag when options are empty",
		tests.testWhenOptionsAreEmpty)
}

type findFlagTests struct{}

func (tests findFlagTests) testWhenFindFlagIsSuccessful(t *testing.T) {
	flagList := getFlagList()
	options := map[string]string{"--file": "file.json"}
	flagFound := FindFlag(flagList, options)

	expectedflag := "--file"

	assert.Equal(t, expectedflag, flagFound)
}

func (tests findFlagTests) testWhenTheGivenOptionsAreNotFound(t *testing.T) {
	flagList := getFlagList()
	options := map[string]string{"--files": "file.json"}
	flagFound := FindFlag(flagList, options)

	assert.Empty(t, flagFound)
}

func (tests findFlagTests) testWhenOptionsAreEmpty(t *testing.T) {
	flagList := getFlagList()
	options := map[string]string{}
	flagFound := FindFlag(flagList, options)

	assert.Empty(t, flagFound)
}

func getFlagList() []string {
	return []string{"--file", "-f", "--stdin", "-i"}
}
