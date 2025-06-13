// Package converter provides tools to convert VTT files to desired JSON format
package captionconverter

import (
	"errors"
	"fmt"
	"log"

	cmdutil "banditsecret/internal/pkg/cmdutil"
)

type Converter interface {
	ConvertVTTToJSON(vttPath, jsonPath string) error
}

// Concrete ConverterService implements Converter
type ConverterService struct {
	pythonExecutable    string
	converterScriptPath string
	cmdRunner           cmdutil.CmdRunner
}

// Factory to create a concrete ConverterService
func NewConverterService(pythonExecutable, converterScriptPath string, cmdRunner cmdutil.CmdRunner) (*ConverterService, error) {

	if pythonExecutable == "" || converterScriptPath == "" {
		return nil, errors.New("pythonExecutable or converterScriptPath cannot be empty")
	}

	return &ConverterService{
		pythonExecutable:    pythonExecutable,
		converterScriptPath: converterScriptPath,
		cmdRunner:           cmdRunner,
	}, nil
}

func (cs *ConverterService) ConvertVTTToJSON(vttFilePath, jsonFilePath string) error {

	log.Printf("Attempting to convert %s to %s", vttFilePath, jsonFilePath)

	if cmdutil.FileExists(jsonFilePath) {
		log.Printf("Json file %s already exists! Skipping download\n", jsonFilePath)
		return nil
	}

	log.Printf("Converting VTT file: %s to JSON.", vttFilePath)
	cmdOutput, err := cs.cmdRunner.CombinedOutput(
		cs.pythonExecutable,
		cs.converterScriptPath,
		vttFilePath,
		jsonFilePath)

	if err != nil {
		return fmt.Errorf("failed to run python script convert captions to json: %w\nOutput: %s", err, cmdOutput)
	}

	log.Println("Done extracting captions (Python script output):", string(cmdOutput))
	return nil
}
