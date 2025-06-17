package main

import (
	"banditsecret/internal/app"
	"fmt"
	"log"
	"os"
	"time"

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

	// TODO: do a readiness check with db prior to continuing
	time.Sleep(8 * time.Second)

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

	v1.GET("/test_get_metadata", func(c *gin.Context) {

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": "could not read body"})
			return
		}
		url := string(body)

		meta, err := appServices.Fetcher.GetMetadata(url, os.Getenv("JSON_CAPTIONS_DIR"))
		if err != nil {
			log.Printf("Failed to get video metadata: %v", err)
			return
		}
		fmt.Println(*meta)
	})

	v1.GET("/test_get_captions", func(c *gin.Context) {

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": "could not read body"})
			return
		}
		url := string(body)

		meta, err := appServices.Fetcher.GetMetadata(url, os.Getenv("JSON_CAPTIONS_DIR"))
		if err != nil {
			log.Printf("Failed to get video metadata: %v", err)
			return
		}

		output, err := appServices.Fetcher.DownloadCaptions(meta.VideoId, url, os.Getenv("VTT_CAPTIONS_DIR"))
		if err != nil {
			log.Printf("Failed to get video metadata: %v", err)
			return
		}

		fmt.Println(output)
	})

	router.Run("0.0.0.0:" + os.Getenv("SERVER_PORT"))
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
		log.Printf("Failed to download captions: %v", err)
		return
	}

	err = appServices.Converter.ConvertVTTToJSON(vttCaptionsFile, meta.CaptionPath)
	if err != nil {
		log.Printf("Failed to convert VTT file to JSON: %s", err)
		return
	}

	captions, err := appServices.Parser.ParseJSON(meta.CaptionPath)
	if err != nil {
		log.Printf("Failed to Parse JSON captions: %s", err)
		return
	}

	ctx := c.Request.Context()
	err = appServices.Loader.LoadCaptions(ctx, meta, captions)
	if err != nil {
		log.Printf("Failed to load captions to db: %s", err)
		return
	}

	err = appServices.Searcher.IndexCaptions(ctx, meta, captions)
	if err != nil {
		log.Printf("Failed to index captions to Elastic Search: %s", err)
		return
	}
}
