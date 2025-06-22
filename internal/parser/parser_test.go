package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {

	want := []CaptionEntry{
		{
			VideoId: "SampleVideoId",
			Start:   0,
			End:     5000,
			Text:    "SampleText1",
		},
		{
			VideoId: "SampleVideoId",
			Start:   6000,
			End:     10000,
			Text:    "SampleText2",
		},
	}
	inpFile := "testdata/sample.json"
	parserService := NewParserService()

	got, err := parserService.ParseJSON(inpFile)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if len(got) < len(want) {
		t.Fatal("ParseJSON failed: Not enough entries in parsed struct")
	}

	for i, entry := range got {
		if entry != want[i] {
			t.Fatalf("ParseJSON failed: parsed entry does not match expected:\nreceived: %+v\nexpected: %+v", entry, want[i])
		}
	}
}

func TestParseJSONFileDoesNotExist(t *testing.T) {
	inpFile := "testdata/does_not_exist.json"
	parserService := NewParserService()

	_, err := parserService.ParseJSON(inpFile)
	if err == nil {
		t.Fatalf("TestParseJSONFileDoesNotExist failed: expected an error, got nil")
	}

	expected := fmt.Sprintf("JSON caption file not found at %s", inpFile)
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("TestParseJSONFileDoesNotExist failed: expected error to contain %s, but got %s", expected, err.Error())
	}
}

func TestParseJSONFileCannotBeUnmarshalled(t *testing.T) {
	inpFile := "testdata/badsample.json"
	parserService := NewParserService()

	_, err := parserService.ParseJSON(inpFile)
	if err == nil {
		t.Fatalf("TestParseJSONFileCannotBeUnmarshalled failed: expected an error, got nil")
	}

	expected := fmt.Sprintf("failed to unmarshal JSON data from %s: json: cannot unmarshal object into Go value of type []parser.CaptionEntry", inpFile)

	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("TestParseJSONFileDoesNotExist failed: expected error to contain \"%s\", but got \"%s\"", expected, err.Error())
	}
}
