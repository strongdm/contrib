package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminServersAdd(t *testing.T) {
	tests := adminServersAddTests{}
	t.Run("Test AdminServersAdd when it is successful",
		tests.testWhenItIsSucessful)
	t.Run("Test AdminServersAdd when find flag returns empty",
		tests.testWhenFindFlagReturnsEmpty)
	t.Run("Test AdminServersAdd when an nonexistent file is passed",
		tests.testWhenAnNonexistentFileIsPassed)
	t.Run("Test AdminServersAdd when an existent file with a wrong content is passed",
		tests.testWhenAnExistentFileWithAWrongContentIsPassed)
	t.Run("Test AdminServersAdd when one of the servers is not registered",
		tests.testWhenOneOfTheServersIsNotRegistered)
}

type adminServersAddTests struct{}

func (tests adminServersAddTests) testWhenItIsSucessful(t *testing.T) {
	adminServiceImpl := NewAdminService()
	adminServiceImpl.patchFindFlag(findFlagMock)
	adminServiceImpl.patchExtractValuesFromJson(extractValuesFromJsonMock)
	adminServiceImpl.patchExecute(executeWithSuccessMock)

	options := map[string]string{"--file": "file.json"}
	actualErr := adminServiceImpl.AdminServersAdd(options)

	assert.Nil(t, actualErr)
}

func (tests adminServersAddTests) testWhenFindFlagReturnsEmpty(t *testing.T) {
	adminServiceImpl := NewAdminService()
	adminServiceImpl.patchFindFlag(findFlagReturningEmptyMock)

	options := map[string]string{"--files": "file.json"}
	actualErr := adminServiceImpl.AdminServersAdd(options)

	assert.Nil(t, actualErr)
}

func (tests adminServersAddTests) testWhenAnNonexistentFileIsPassed(t *testing.T) {
	adminServiceImpl := NewAdminService()
	adminServiceImpl.patchFindFlag(findFlagMock)
	adminServiceImpl.patchExtractValuesFromJson(extractValuesFromJsonWithABaddlyFormattedJsonMock)

	options := map[string]string{"--file": "file.json"}

	actualErr := adminServiceImpl.AdminServersAdd(options)

	expectedErr := errors.New("invalid character '}' looking for beginning of object key string")

	assert.Equal(t, expectedErr.Error(), actualErr.Error())
}

func (tests adminServersAddTests) testWhenAnExistentFileWithAWrongContentIsPassed(t *testing.T) {
	adminServiceImpl := NewAdminService()
	adminServiceImpl.patchFindFlag(findFlagMock)
	adminServiceImpl.patchExtractValuesFromJson(extractValuesFromJsonWithANonExistentJsonFileMock)

	options := map[string]string{"--file": "file.json"}
	actualErr := adminServiceImpl.AdminServersAdd(options)

	expectedErr := errors.New("open file.json: no such file or directory")

	assert.Equal(t, expectedErr.Error(), actualErr.Error())
}

func (tests adminServersAddTests) testWhenOneOfTheServersIsNotRegistered(t *testing.T) {
	adminServiceImpl := NewAdminService()
	adminServiceImpl.patchFindFlag(findFlagMock)
	adminServiceImpl.patchExtractValuesFromJson(extractValuesFromJsonMock)
	adminServiceImpl.patchExecute(executeWithoutSuccessMock)

	options := map[string]string{"--file": "file.json"}
	actualErr := adminServiceImpl.AdminServersAdd(options)

	assert.Nil(t, actualErr)
}
func TestGetOptions(t *testing.T) {
	tests := getOptionsTests{}
	t.Run("Test getOptions when it is successful",
		tests.testWhenItIsSucessful)
	t.Run("Test getOptions when the server contains tags",
		tests.testWhenTheServerContainsTags)
	t.Run("Test getOptions when the server contains private key",
		tests.testWhenTheServerContainsPrivateKey)
	t.Run("Test getOptions when the server not contain a name",
		tests.testWhenTheServerNotContainAName)
	t.Run("Test getOptions when the server not contain a type",
		tests.testWhenTheServerNotContainAType)
	t.Run("Test getOptions when the server not contain a name and type",
		tests.testWhenTheServerNotContainANameAndType)
	t.Run("Test getOptions when the server contain an attribute empty",
		tests.testWhenTheServerContainAnAttributeEmpty)
	t.Run("Test getOptions when server is empty",
		tests.testWhenTheServerIsEmpty)
}

type getOptionsTests struct{}

func (tests getOptionsTests) testWhenItIsSucessful(t *testing.T) {
	server := getServer()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMap()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerContainsTags(t *testing.T) {
	server := getServerWithTags()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMapWithTags()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerContainsPrivateKey(t *testing.T) {
	server := getServerWithPrivateKey()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMapWithPrivateKey()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerNotContainAName(t *testing.T) {
	server := getServerWithoutName()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMap()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerNotContainAType(t *testing.T) {
	server := getServerWithoutType()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMap()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerNotContainANameAndType(t *testing.T) {
	server := getServerWithoutNameAndType()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMap()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerContainAnAttributeEmpty(t *testing.T) {
	server := getServerWithAnAttributeEmpty()
	actualOptionsMap := getOptions(server)
	expectedOptionsMap := getOptionsMap()
	assert.Equal(t, expectedOptionsMap, actualOptionsMap)
}

func (tests getOptionsTests) testWhenTheServerIsEmpty(t *testing.T) {
	server := map[string]interface{}{}
	actualOptionsMap := getOptions(server)
	assert.Empty(t, actualOptionsMap)
}

func TestTreatKey(t *testing.T) {
	tests := treatKeyTests{}
	t.Run("Test treatKey when it is successful",
		tests.testWhenItIsSucessful)
	t.Run("Test treatKey when the key is not modified",
		tests.testWhenTheKeyIsNotModified)
}

type treatKeyTests struct{}

func (tests treatKeyTests) testWhenItIsSucessful(t *testing.T) {
	key := "portOverride"
	actualKey := treatKey(key)
	expectedKey := "port-override"
	assert.Equal(t, expectedKey, actualKey)
}

func (tests treatKeyTests) testWhenTheKeyIsNotModified(t *testing.T) {
	key := "hostname"
	actualKey := treatKey(key)
	expectedKey := "hostname"
	assert.Equal(t, expectedKey, actualKey)
}

func TestTreatTags(t *testing.T) {
	tests := treatTagsTests{}
	t.Run("Test treatTags when it is successful",
		tests.testWhenItIsSucessful)
	t.Run("Test treatTags when tags map contain many tags",
		tests.testWhenTagsMapContainManyTags)
	t.Run("Test treatTags when the tags map is empty",
		tests.testWhenTheTagsMapIsEmpty)
}

type treatTagsTests struct{}

func (tests treatTagsTests) testWhenItIsSucessful(t *testing.T) {
	tagsMap := getTagsMap()
	actualTags := treatTags(tagsMap)
	expectedTags := "key1=value1"
	assert.Equal(t, expectedTags, actualTags)
}

func (tests treatTagsTests) testWhenTagsMapContainManyTags(t *testing.T) {
	tagsMap := getTagsMapWithManyTags()
	actualTags := treatTags(tagsMap)
	expectedTags := "key1=value1,key2=value2"
	assert.Equal(t, expectedTags, actualTags)
}

func (tests treatTagsTests) testWhenTheTagsMapIsEmpty(t *testing.T) {
	tagsMap := map[string]interface{}{}
	actualTags := treatTags(tagsMap)
	assert.Empty(t, actualTags)
}

func findFlagMock(flagList []string, optionsMap map[string]string) string {
	return "--file"
}

func findFlagReturningEmptyMock(flagList []string, optionsMap map[string]string) string {
	return ""
}

func extractValuesFromJsonMock(file string) ([]map[string]interface{}, error) {
	return getMapDataList(), nil
}

func extractValuesFromJsonWithANonExistentJsonFileMock(file string) ([]map[string]interface{}, error) {
	return nil, errors.New("open file.json: no such file or directory")
}

func extractValuesFromJsonWithABaddlyFormattedJsonMock(file string) ([]map[string]interface{}, error) {
	return nil, errors.New("invalid character '}' looking for beginning of object key string")
}

func executeWithSuccessMock(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder) {
	stdout := new(strings.Builder)
	stdout.WriteString("")
	stderr := new(strings.Builder)
	stderr.WriteString("")
	return *stdout, *stderr
}

func executeWithoutSuccessMock(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder) {
	stdout := new(strings.Builder)
	stdout.WriteString("")
	stderr := new(strings.Builder)
	stderr.WriteString("stderr")
	return *stdout, *stderr
}

func getMapDataList() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"hostname": "192.168.101.192",
			"name":     "Example Raw TCP",
			"type":     "rawtcp",
		},
	}
}

func getServer() map[string]interface{} {
	return map[string]interface{}{
		"name":     "Example Server",
		"hostname": "hostname",
		"type":     "type",
	}
}

func getServerWithTags() map[string]interface{} {
	return map[string]interface{}{
		"name":     "Example Server",
		"hostname": "hostname",
		"type":     "type",
		"tags": map[string]interface{}{
			"key1": "value1",
		},
	}
}

func getServerWithPrivateKey() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Example Server",
		"hostname":   "hostname",
		"type":       "type",
		"privateKey": "private key",
	}
}

func getServerWithoutName() map[string]interface{} {
	return map[string]interface{}{
		"hostname": "hostname",
		"type":     "type",
	}
}

func getServerWithoutType() map[string]interface{} {
	return map[string]interface{}{
		"name":     "Example Server",
		"hostname": "hostname",
	}
}

func getServerWithoutNameAndType() map[string]interface{} {
	return map[string]interface{}{
		"hostname": "hostname",
	}
}

func getServerWithAnAttributeEmpty() map[string]interface{} {
	return map[string]interface{}{
		"name":     "Example Server",
		"hostname": "hostname",
		"type":     "type",
		"username": "",
	}
}

func getOptionsMap() map[string]string {
	return map[string]string{
		"--hostname": "hostname",
	}
}

func getOptionsMapWithTags() map[string]string {
	return map[string]string{
		"--hostname": "hostname",
		"--tags":     "key1=value1",
	}
}

func getOptionsMapWithPrivateKey() map[string]string {
	return map[string]string{
		"--hostname":     "hostname",
		"--private-key=": `"private key"`,
	}
}

func getTagsMap() map[string]interface{} {
	return map[string]interface{}{
		"key1": "value1",
	}
}

func getTagsMapWithManyTags() map[string]interface{} {
	return map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
}
