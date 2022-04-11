package util

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp/syntax"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

const (
	FILE_REGEX_PATTERN      = `^--file [\.\/\w,-]+\.[A-Za-z]+$`
	F_REGEX_PATTERN         = `^-f [\.\/\w,-]+\.[A-Za-z]+$`
	STDIN_REGEX_PATTERN     = `^--stdin$`
	I_REGEX_PATTERN         = `^-i$`
	DEFAULT_JSON_FILENAME   = "file.json"
	DEFAULT_PATTERN_TMPFILE = "tmpfile"
)

func TestMapCommandArguments(t *testing.T) {
	tests := mapCommandArgumentsTests{}
	t.Run("Test MapCommandArguments when it is successful",
		tests.testWhenItIsSuccessful)
	t.Run("Test MapCommandArguments when it returns a map empty",
		tests.testWhenItReturnsAnEmptyMap)
}

type mapCommandArgumentsTests struct{}

func (tests mapCommandArgumentsTests) testWhenItIsSuccessful(t *testing.T) {
	cliFlagList := getCliFlagList()
	strList := getCorrectStrList()
	actualArgsMapping := MapCommandArguments(strList, cliFlagList)

	expectedArgsMapping := getArgsMapping()

	assert.Equal(t, expectedArgsMapping, actualArgsMapping)
}

func (tests mapCommandArgumentsTests) testWhenItReturnsAnEmptyMap(t *testing.T) {
	cliFlagList := getCliFlagList()
	strList := getIncorrectStrList()
	actualArgsMapping := MapCommandArguments(strList, cliFlagList)

	assert.Empty(t, actualArgsMapping)
}

func TestFlagHasName(t *testing.T) {
	tests := flagHasNameTests{}
	t.Run("Test FlagHasName when it is successful",
		tests.testWhenItIsSuccessful)
	t.Run("Test FlagHasName when it can not find the flag",
		tests.testWhenItCanNotFindTheFlag)
}

type flagHasNameTests struct{}

func (tests flagHasNameTests) testWhenItIsSuccessful(t *testing.T) {
	cliFlag := getCliFlag()
	flag := "--file"
	actualFoundFlag := FlagHasName(cliFlag, flag)

	assert.True(t, actualFoundFlag)
}

func (tests flagHasNameTests) testWhenItCanNotFindTheFlag(t *testing.T) {
	cliFlag := getCliFlag()
	flag := "--files"
	actualFoundFlag := FlagHasName(cliFlag, flag)

	assert.False(t, actualFoundFlag)
}

func TestExtractValuesFromJson(t *testing.T) {
	tests := extractValuesFromJsonTests{}
	t.Run("Test ExtractValuesFromJson when the data is successfully extracted",
		tests.testWhenTheDataIsSuccessfullyExtracted)
	t.Run("Test ExtractValuesFromJson when an nonexistent file is passed",
		tests.testWhenAnNonexistentFileIsPassed)
	t.Run("Test ExtractValuesFromJson when an existent file with a wrong content is passed",
		tests.testWhenAnExistentFileWithAWrongContentIsPassed)
}

type extractValuesFromJsonTests struct{}

func (tests extractValuesFromJsonTests) testWhenTheDataIsSuccessfullyExtracted(t *testing.T) {
	file, _ := os.Create(DEFAULT_JSON_FILENAME)
	defer os.Remove(file.Name())
	file.WriteString(getCorrectJsonFileContent())

	actualData, actualErr := ExtractValuesFromJson(file.Name())

	expectedData := getJsonData()

	assert.Equal(t, expectedData, actualData)
	assert.Nil(t, actualErr)
}

func (tests extractValuesFromJsonTests) testWhenAnNonexistentFileIsPassed(t *testing.T) {
	actualData, actualErr := ExtractValuesFromJson(DEFAULT_JSON_FILENAME)

	expectedErr := errors.New("open file.json: no such file or directory")

	assert.Nil(t, actualData)
	assert.Equal(t, expectedErr.Error(), actualErr.Error())
}

func (tests extractValuesFromJsonTests) testWhenAnExistentFileWithAWrongContentIsPassed(t *testing.T) {
	file, _ := os.Create(DEFAULT_JSON_FILENAME)
	defer os.Remove(file.Name())
	file.WriteString(getIncorrectJsonFileContent())

	actualData, actualErr := ExtractValuesFromJson(file.Name())

	expectedErr := errors.New("invalid character '}' looking for beginning of object key string")

	assert.Nil(t, actualData)
	assert.Equal(t, expectedErr.Error(), actualErr.Error())
}

func TestGetUserInput(t *testing.T) {
	tests := getUserInputTests{}
	t.Run("Test GetUserInput when the user input is valid",
		tests.testWhenTheUserInputIsValid)
	t.Run("Test GetUserInput when the user input a badly formatted json",
		tests.testWhenTheUserInputABadlyFormattedJson)
}

type getUserInputTests struct{}

func (tests getUserInputTests) testWhenTheUserInputIsValid(t *testing.T) {
	tmpFile, _ := ioutil.TempFile("", DEFAULT_PATTERN_TMPFILE)
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()
	tmpFile.WriteString(getCorrectJsonFileContent())
	tmpFile.Seek(0, 0)

	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()
	os.Stdin = tmpFile

	actualData, actualErr := GetUserInput()

	expectedData := getJsonData()

	assert.Equal(t, expectedData, actualData)
	assert.Nil(t, actualErr)

	tmpFile.Close()
}

func (tests getUserInputTests) testWhenTheUserInputABadlyFormattedJson(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "tempfile")
	assert.Nil(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(getIncorrectJsonFileContent())
	assert.Nil(t, err)

	_, err = tempFile.Seek(0, 0)
	assert.Nil(t, err)

	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()
	os.Stdin = tempFile

	actualData, actualErr := GetUserInput()

	expectedErr := errors.New("invalid character '}' looking for beginning of object key string")

	assert.Nil(t, actualData)
	assert.Equal(t, expectedErr.Error(), actualErr.Error())

	err = tempFile.Close()
	assert.Nil(t, err)
}

func TestConvertStrSliceToStr(t *testing.T) {
	tests := convertStrSliceToStrTests{}
	t.Run("Test ConvertStrSliceToStr when it is successful",
		tests.testWhenItIsSuccessful)
	t.Run("Test ConvertStrSliceToStr when the given slice is empty",
		tests.testWhenTheGivenSliceIsEmpty)
}

type convertStrSliceToStrTests struct{}

func (tests convertStrSliceToStrTests) testWhenItIsSuccessful(t *testing.T) {
	strList := getCorrectStrList()
	actualStr := ConvertStrSliceToStr(strList)

	expectedStr := "--file file.json"

	assert.Equal(t, expectedStr, actualStr)
}

func (tests convertStrSliceToStrTests) testWhenTheGivenSliceIsEmpty(t *testing.T) {
	strList := []string{}
	actualStr := ConvertStrSliceToStr(strList)

	assert.Empty(t, actualStr)
}

func TestCheckRegexMatch(t *testing.T) {
	tests := checkRegexMatchTests{}
	t.Run("Test CheckRegexMatch when it is successful",
		tests.testWhenItIsSuccessful)
	t.Run("Test CheckRegexMatch when it does not matches",
		tests.testWhenItDoesNotMatches)
	t.Run("Test CheckRegexMatch when an invalid regex is passed",
		tests.testWhenAnInvalidRegexIsPassed)
}

type checkRegexMatchTests struct{}

func (tests checkRegexMatchTests) testWhenItIsSuccessful(t *testing.T) {
	regexList := getCorrectRegexList()
	arguments := getCorrectArguments()
	actualMatched, actualErr := CheckRegexMatch(regexList, arguments)

	assert.True(t, actualMatched)
	assert.Nil(t, actualErr)
}

func (tests checkRegexMatchTests) testWhenItDoesNotMatches(t *testing.T) {
	regexList := getCorrectRegexList()
	arguments := getIncorrectArguments()
	actualMatched, actualErr := CheckRegexMatch(regexList, arguments)

	assert.False(t, actualMatched)
	assert.Nil(t, actualErr)
}

func (tests checkRegexMatchTests) testWhenAnInvalidRegexIsPassed(t *testing.T) {
	regexList := getIncorrectRegexList()
	arguments := getIncorrectArguments()
	actualMatched, actualErr := CheckRegexMatch(regexList, arguments)

	expectedErr := syntax.Error{Code: "invalid escape sequence", Expr: "\\K"}

	assert.False(t, actualMatched)
	assert.Equal(t, &expectedErr, actualErr)
}

func getCliFlag() cli.Flag {
	return cli.StringFlag{
		Name:  "file,f",
		Usage: "load from a JSON file",
	}
}

func getCliFlagList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "file,f",
			Usage: "load from a JSON file",
		},
		cli.BoolFlag{
			Name:  "stdin,i",
			Usage: "load from stdin",
		},
	}
}

func getCorrectStrList() []string {
	return []string{"--file", "file.json"}
}

func getIncorrectStrList() []string {
	return []string{"--files", "file.json"}
}

func getCorrectJsonFileContent() string {
	return `
		[
			{
				"name": "Example Raw TCP",
				"hostname": "192.168.101.192",
				"type": "rawtcp"
			}
		]
	`
}

func getIncorrectJsonFileContent() string {
	return `
		[
			{
				"name": "Example Raw TCP",
				"hostname": "192.168.101.192",
				"type": "rawtcp",
			}
		]
	`
}

func getJsonData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"hostname": "192.168.101.192",
			"name":     "Example Raw TCP",
			"type":     "rawtcp",
		},
	}
}

func getArgsMapping() map[string]string {
	return map[string]string{"--file": "file.json"}
}

func getCorrectRegexList() []string {
	return []string{
		FILE_REGEX_PATTERN,
		F_REGEX_PATTERN,
		STDIN_REGEX_PATTERN,
		I_REGEX_PATTERN,
	}
}

func getIncorrectRegexList() []string {
	return []string{`^--file \K[\.\/\w,-]+\.[A-Za-z]+$`}
}

func getCorrectArguments() string {
	return "--file file.json"
}

func getIncorrectArguments() string {
	return "--files file.json"
}
