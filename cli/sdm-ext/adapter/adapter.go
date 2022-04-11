package adapter

import (
	"ext/service"
	"fmt"
)

func Servers(commandName string, mappedOptions map[string]string) error {
	adminService := service.NewAdminService()

	switch commandName {
	case "add":
		return adminService.AdminServersAdd(mappedOptions)
	default:
		return fmt.Errorf("unknown command name: %s", commandName)
	}
}
