package main

import (
	"banditsecret/internal/app"
	searcher "banditsecret/internal/search"
	"banditsecret/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type CaptionRepository = storage.CaptionRepository
type CaptionSearchRepository = searcher.CaptionSearchRepository

func main() {

	// Load env vars
	err := godotenv.Load("./envs/app.env")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: do a readiness check with db prior to continuing
	time.Sleep(20 * time.Second)

	// Init db connection
	db, err := storage.InitDb()
	if err != nil {
		log.Fatalf("failed to init db: %w", err)
	}
	defer db.Close()
	var captionRepo CaptionRepository = storage.NewSQLCaptionRepository(db)

	// Init search engine connection
	esClient, err := searcher.InitEsClient()
	if err != nil {
		log.Fatalf("Failed to init Elasticsearch client: %v", err)
	}
	var captionSearchRepo CaptionSearchRepository = searcher.NewElasticSearchRepository(esClient)

	// Init app services
	appServices, err := app.NewApplicationServices(captionRepo, captionSearchRepo)
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

	v1.GET("/search", func(c *gin.Context) {
		queryHandler(c, appServices)
	})

	router.Run("0.0.0.0:" + os.Getenv("SERVER_PORT"))
}

func queryHandler(c *gin.Context, s *app.ApplicationServices) {
	query := c.Query("query")

	ctx := c.Request.Context()
	res, err := s.Searcher.SearchCaptions(ctx, os.Getenv("CAPTIONS_INDEX"), query)

	if err != nil {
		log.Printf("query failed %s", err)
		return
	}

	c.JSON(http.StatusOK, res)

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
