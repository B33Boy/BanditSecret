// Package converter provides tools to convert VTT files to desired JSON format
package captionconverter

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

type Converter interface {
	ConvertVTTToJSON(vttPath, jsonPath string) error
}

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

// Concrete ConverterService implements Converter
type ConverterService struct {
	pythonExecutable    string
	converterScriptPath string
	cmdRunner           CmdRunner
}

// Factory to create a concrete ConverterService
func NewConverterService(pythonExecutable, converterScriptPath string, cmdRunner CmdRunner) (*ConverterService, error) {

	if pythonExecutable == "" || converterScriptPath == "" {
		return nil, errors.New("pythonExecutable or converterScriptPath cannot be empty")
	}

	return &ConverterService{
		pythonExecutable:    pythonExecutable,
		converterScriptPath: converterScriptPath,
		cmdRunner:           cmdRunner,
	}, nil
}

func (cs *ConverterService) ConvertVTTToJSON(vttPath, jsonPath string) error {

	log.Printf("Converting VTT file: %s to JSON.", vttPath)
	cmdOutput, err := cs.cmdRunner.CombinedOutput(
		cs.pythonExecutable,
		cs.converterScriptPath,
		vttPath,
		jsonPath)

	if err != nil {
		return fmt.Errorf("failed to run python script convert captions to json: %w\nOutput: %s", err, cmdOutput)
	}

	log.Println("Done extracting captions (Python script output):", string(cmdOutput))
	return nil
}
