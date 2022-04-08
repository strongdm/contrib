package service

import "os/exec"

func NewSdmServiceMock(runCommand func(cmd *exec.Cmd)) *sdmServiceImpl {
	return &sdmServiceImpl{runCommand}
}
