package service

import (
	"ext/util"
	"strings"
)

type adminService interface {
	execute(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder)
}

type AdminServiceImpl struct {
	findFlag              func(flagList []string, optionsMap map[string]string) string
	extractValuesFromJson func(file string) ([]map[string]interface{}, error)
	execute               func(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder)
}

func NewAdminService() *AdminServiceImpl {
	sdmService := NewSdmService()
	return &AdminServiceImpl{
		util.FindFlag,
		util.ExtractValuesFromJson,
		sdmService.execute,
	}
}

func (s *AdminServiceImpl) patchFindFlag(findFlag func(flagList []string, optionsMap map[string]string) string) {
	s.findFlag = findFlag
}

func (s *AdminServiceImpl) patchExtractValuesFromJson(extractValuesFromJson func(file string) ([]map[string]interface{}, error)) {
	s.extractValuesFromJson = extractValuesFromJson
}

func (s *AdminServiceImpl) patchExecute(execute func(commands string, options map[string]string, postOptions string) (strings.Builder, strings.Builder)) {
	s.execute = execute
}
