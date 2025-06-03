package main

import (
	"banditsecret/internal/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting Extraction")
	err := parser.ExtractCaptions("https://youtu.be/bkO3a50tvco")

	if err != nil {
		fmt.Println(err)
	}
}
