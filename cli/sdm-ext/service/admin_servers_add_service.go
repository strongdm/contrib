package service

import (
	"ext/util"
	"fmt"
	"strings"
	"unicode"
)

func (a AdminService) AdminServersAdd(options map[string]string) error {
	flagList := []string{"--file", "-f", "--stdin", "-i"}
	flag := util.FindFlag(flagList, options)

	var servers []map[string]interface{}
	var err error

	if flag == "--file" || flag == "-f" {
		servers, err = util.ExtractValuesFromJson(options[flag])
		if err != nil {
			return err
		}
	} else if flag == "--stdin" || flag == "-i" {
		servers, err = util.GetUserInput()
		if err != nil {
			return err
		}
	}

	for _, server := range servers {
		serverName := fmt.Sprint(server["name"])
		serverType := fmt.Sprint(server["type"])

		_, stderr := execute(
			fmt.Sprintf("admin servers add %s", serverType),
			getOptions(server),
			serverName,
		)

		if stderr.String() == "" {
			fmt.Printf("Server \"%s\" successfully registered\n", serverName)
		} else {
			fmt.Printf("There was an error registering the \"%s\" server\n", serverName)
		}
	}

	return nil
}

func getOptions(server map[string]interface{}) map[string]string {
	options := map[string]string{}

	for key, value := range server {
		if key != "name" && key != "type" {
			key = treatKey(key)
			if value != "" {
				if key == "tags" {
					value = treatTags(value.(map[string]interface{}))
				}
				if key == "private-key" {
					options["--"+key+"="] = fmt.Sprintf(`"%s"`, value)
				} else {
					options["--"+key] = fmt.Sprint(value)
				}
			}
		}
	}

	return options
}

func treatKey(key string) string {
	for _, character := range key {
		if character >= 'A' && character <= 'Z' {
			key = strings.Replace(key, string(character), "-"+string(unicode.ToLower(character)), -1)
		}
	}

	return key
}

func treatTags(tagsMap map[string]interface{}) string {
	tags := ""
	i := 0
	for key, value := range tagsMap {
		tags += fmt.Sprintf("%s=%s", key, value)
		if i < len(tagsMap)-1 {
			tags += ","
		}
		i++
	}

	return tags
}
