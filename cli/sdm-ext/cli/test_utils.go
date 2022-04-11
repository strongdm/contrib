package cli

import (
	"os/exec"
)

func NewSdmMock(runCommandMock func(cmd *exec.Cmd)) *sdmImpl {
	return &sdmImpl{runCommandMock}
}
