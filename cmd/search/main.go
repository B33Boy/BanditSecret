package main

import (
	"banditsecret/internal/storage"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// // TODO: This should be the entrypoint for the server

	// // TEMP: Testing end to end with sample
	// fmt.Println("Starting Extraction")
	// captionResult, err := parser.ExtractCaptions("https://youtu.be/jVpsLMCIB0Y")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(*captionResult)

	// Load .env
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load JSON to DB
	// captionResult.CaptionPath
	captions, err := storage.LoadCaptionsFromJson("tmp/captions_parsed/jVpsLMCIB0Y.en.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	storage.StoreCaptionsToDb(captions)

	// for _, cap := range captions {
	// 	fmt.Printf("ID: %s, Start: %s, End: %s, Text: %s, \n", cap.Id, cap.Start, cap.End, cap.Text)
	// }

}
