package util

import (
	"github.com/urfave/cli"
)

func GetAdminServersAddFileFlag() cli.StringFlag {
	return cli.StringFlag{
		Name:  "file,f",
		Usage: "load from a JSON file",
	}
}

func GetAdminServersAddStdinFlag() cli.BoolFlag {
	return cli.BoolFlag{
		Name:  "stdin,i",
		Usage: "load from stdin",
	}
}

func FindFlag(flagList []string, optionsMap map[string]string) string {
	for _, flag := range flagList {
		for key := range optionsMap {
			if flag == key {
				return flag
			}
		}
	}

	return ""
}
