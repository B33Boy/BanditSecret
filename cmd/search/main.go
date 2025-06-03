package main

import (
	"banditsecret/internal/parser"
	"fmt"
)

func main() {
	// TODO: This should be the entrypoint for the server

	// TEMP: Testing end to end with sample
	fmt.Println("Starting Extraction")
	CaptionResult, err := parser.ExtractCaptions("https://youtu.be/jVpsLMCIB0Y")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*CaptionResult)

}
