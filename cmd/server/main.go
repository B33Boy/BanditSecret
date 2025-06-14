package main

import (
	"banditsecret/internal/app"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// type ApplicationServices app.ApplicationServices

func main() {

	// Load env vars
	err := godotenv.Load("./envs/app.env")
	if err != nil {
		log.Fatal(err)
	}

	// Get connections to external services (e.g. mysql db, elastic search)
	db, esClient, err := app.InitConnections()
	if err != nil {
		log.Fatalf("Application connections failed to initialize: %v", err)
	}
	defer db.Close()

	// Init app services
	appServices, err := app.NewApplicationServices(db, esClient)
	if err != nil {
		log.Fatalf("Failed to initialize application services: %v", err)
	}

	startServer(appServices)
}

func startServer(appServices *app.ApplicationServices) {
	// set up API server
	router := gin.Default()
	v1 := router.Group("/v1")

	v1.POST("/captions", func(c *gin.Context) {
		ingestVideoHandler(c, appServices)
	})

	router.Run("localhost:" + os.Getenv("SERVER_PORT"))
}

func ingestVideoHandler(c *gin.Context, appServices *app.ApplicationServices) {
	// Get request and look for URL
	// Request will have raw yt url in body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "could not read body"})
		return
	}
	url := string(body)

	// Note metadata's CaptionPath refers to the to-be generated json file
	// TODO: Add context to fetcher functions
	meta, err := appServices.Fetcher.GetMetadata(url, os.Getenv("JSON_CAPTIONS_DIR"))
	if err != nil {
		log.Fatalf("Failed to get video metadata: %v", err)
	}

	vttCaptionsDir := os.Getenv("VTT_CAPTIONS_DIR")
	vttCaptionsFile, err := appServices.Fetcher.DownloadCaptions(meta.VideoId, url, vttCaptionsDir)
	if err != nil {
		log.Fatalf("Failed to download captions: %v", err)
	}

	err = appServices.Converter.ConvertVTTToJSON(vttCaptionsFile, meta.CaptionPath)
	if err != nil {
		log.Fatalf("Failed to convert VTT file to JSON: %s", err)
	}

	captions, err := appServices.Parser.ParseJSON(meta.CaptionPath)
	if err != nil {
		log.Fatalf("Failed to Parse JSON captions: %s", err)
	}

	ctx := c.Request.Context()
	err = appServices.Loader.LoadCaptions(ctx, meta, captions)
	if err != nil {
		log.Fatalf("Failed to load captions to db: %s", err)
	}

	err = appServices.Searcher.IndexCaptions(ctx, meta, captions)
	if err != nil {
		log.Fatalf("Failed to index captions to Elastic Search: %s", err)
	}
}
