package main

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/search"
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
		log.Fatalf("ExtractCaptions failed: %v", err)
	}

	fmt.Println(*captionMetadata)

	err = godotenv.Load("./envs/app.env")
	if err != nil {
		log.Fatal(err)
	}

	// Load JSON to DB
	// captionMetadata.CaptionPath
	captions, err := storage.LoadCaptionsFromJson(parsedDir + "jVpsLMCIB0Y.en.json")

	if err != nil {
		log.Fatalf("LoadCaptionsFromJson failed: %v", err)
	}

	db, err := storage.InitDb()

	if err != nil {
		log.Fatalf("InitDB failed: %v", err)
	}

	storage.StoreVideoInfoToDb(db, captionMetadata)

	storage.StoreCaptionsToDb(db, captions)

	search.StoreCaptionsToSearchEngine(captions)

	// for _, cap := range captionsList {
	// 	fmt.Printf("ID: %s, Start: %s, End: %s, Text: %s, \n", cap.Id, cap.Start, cap.End, cap.Text)
	// }

	// search.StoreCaptionsToSearchEngine()
}
