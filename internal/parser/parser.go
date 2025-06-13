package parser

import (
	"encoding/json"
	"fmt"
	"os"

	cmdutil "banditsecret/internal/pkg/cmdutil"
)

type Parser interface {
	ParseJSON(jsonFile string) ([]CaptionEntry, error)
}

type ParserService struct{}

func NewParserService() *ParserService {
	return &ParserService{}
}

func (s *ParserService) ParseJSON(jsonFile string) ([]CaptionEntry, error) {

	if !cmdutil.FileExists(jsonFile) {
		return nil, fmt.Errorf("JSON caption file not found at %s", jsonFile)
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file %s: %w", jsonFile, err)
	}

	var captions []CaptionEntry

	err = json.Unmarshal(data, &captions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", jsonFile, err)
	}
	return captions, nil
}
