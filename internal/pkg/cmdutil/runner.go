package cmdutil

import (
	"os"
	"os/exec"
)

// Defines the interface to execute external commands
type CmdRunner interface {
	CombinedOutput(name string, arg ...string) ([]byte, error)
	Output(name string, arg ...string) ([]byte, error)
}

// Concrete CmdRunner
type defaultCmdRunner struct{}

func (cr *defaultCmdRunner) CombinedOutput(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.CombinedOutput()
}

func (cr *defaultCmdRunner) Output(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.Output()
}

// Factory to return a DefaultCmdRunner struct
func NewDefaultCmdRunner() CmdRunner {
	return &defaultCmdRunner{}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
