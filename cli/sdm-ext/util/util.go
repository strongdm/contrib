package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

var MapCommandArguments = func(arguments []string, flags []cli.Flag) map[string]string {
	argsMapping := map[string]string{}

	previousArgIsFlag := false

	for index, arg := range arguments {
		foundFlag := false
		for _, flag := range flags {
			if arg[0] == '-' && FlagHasName(flag, arg) {
				argsMapping[arg] = ""
				foundFlag = true
				break
			}
		}
		if !foundFlag && previousArgIsFlag {
			argsMapping[arguments[index-1]] = arg
		}
		previousArgIsFlag = foundFlag
	}

	return argsMapping
}

func FlagHasName(flag cli.Flag, argKey string) bool {
	foundFlag := false
	for _, flagName := range strings.Split(flag.GetName(), ",") {
		if argKey == "-"+flagName || argKey == "--"+flagName {
			foundFlag = true
			break
		}
	}
	return foundFlag
}

var ExtractValuesFromJson = func(file string) ([]map[string]interface{}, error) {
	readFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal([]byte(readFile), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetUserInput() ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

var ConvertStrSliceToStr = func(strList []string) string {
	fmt.Print() // Needed because mock
	strs := ""
	for i, str := range strList {
		strs += str
		if i < len(strList)-1 {
			strs += " "
		}
		i++
	}

	return strs
}

var CheckRegexMatch = func(regexList []string, arguments string) (bool, error) {
	var matched bool
	var err error
	for _, regex := range regexList {
		matched, err = regexp.MatchString(regex, strings.TrimSpace(arguments))
		if err != nil {
			return false, err
		}
		if matched {
			break
		}
	}

	return matched, nil
}
