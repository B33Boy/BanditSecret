package main

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/storage"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// TODO: This should be the entrypoint for the server

	// TEMP: Testing end to end with sample
	fmt.Println("Starting Extraction")
	parsedDir := "tmp/captions_parsed/"
	captionMetadata, err := parser.ExtractCaptions("https://youtu.be/jVpsLMCIB0Y", parsedDir)

	if err != nil {
		log.Fatalf("ExtractCaptions failed: %w", err)
	}

	fmt.Println(*captionMetadata)

	// Load .env
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Load JSON to DB
	// captionMetadata.CaptionPath
	captionsList, err := storage.LoadCaptionsFromJson(parsedDir + "jVpsLMCIB0Y.en.json")

	if err != nil {
		log.Fatalf("LoadCaptionsFromJson failed: %w", err)
	}

	db, err := storage.InitDb()

	storage.StoreVideoInfoToDb(db, captionMetadata)

	storage.StoreCaptionsToDb(db, captionsList)

	// for _, cap := range captions {
	// 	fmt.Printf("ID: %s, Start: %s, End: %s, Text: %s, \n", cap.Id, cap.Start, cap.End, cap.Text)
	// }

}
