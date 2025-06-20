package parser

import (
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

	parserService := NewParserService()
	got, err := parserService.ParseJSON("testdata/sample.json")
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
