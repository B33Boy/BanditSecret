package main

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/pkg/captionconverter"
	"banditsecret/internal/pkg/cmdutil"
	"banditsecret/internal/pkg/ytdlp"
	"banditsecret/internal/search"
	"banditsecret/internal/storage"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	fetchYTService   *ytdlp.FetchYTService
	converterService *captionconverter.ConverterService
	parserService    *parser.ParserService
	loaderService    *storage.LoaderService
	searchService    *search.SearchService
)

func main() {

	// Load env vars
	err := godotenv.Load("./envs/app.env")
	if err != nil {
		log.Fatal(err)
	}

	// Init db connection
	db, err := storage.InitDb()
	if err != nil {
		log.Fatalf("Failed to init db: %w", err)
	}
	defer db.Close()

	// Init esclient
	esClient, err := search.InitEsClient()
	if err != nil {
		log.Fatalf("Failed to init Elasticsearch client: %w", err)
	}

	// Get paths for python and ytdlp executables
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	projectRoot := filepath.Dir(filepath.Dir(exePath))
	fmt.Println("===> ", projectRoot)

	ytdlpExecutable := "yt-dlp"
	pythonExecutable := filepath.Join(projectRoot, "venv", "Scripts", "python")
	converterScriptPath := filepath.Join(projectRoot, "scripts", "extract_captions.py")

	// Initialize all services
	cmdRunner := cmdutil.NewDefaultCmdRunner()

	fetchYTService, err = ytdlp.NewFetchYTService(ytdlpExecutable, cmdRunner)
	if err != nil {
		log.Fatalf("NewFetchYTService failed: %w", err)
	}

	converterService, err = captionconverter.NewConverterService(pythonExecutable, converterScriptPath, cmdRunner)
	if err != nil {
		log.Fatalf("NewConverterService failed: %w", err)
	}

	parserService = parser.NewParserService()

	loaderService = storage.NewLoaderService(db)

	// TODO: store as env var
	captionsIndex := "captions"
	searchService = search.NewSearchService(esClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchService.CreateIndex(ctx, captionsIndex)

	// set up API server
	router := gin.Default()
	router.POST("/captions", getCaptionsForVideo)
	// router.POST("/captions", handleYoutubeVideo)
	router.Run("localhost:6969")
}

func getCaptionsForVideo(c *gin.Context) {
	// Get request and look for URL
	// Request will have raw yt url in body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "could not read body"})
		return
	}
	url := string(body)
	// c.String(200, "You sent: %s", string(body))

	// TODO: Should set this as an env var
	jsonCaptionsDir := "tmp/captions_parsed/"

	// Note metadata's CaptionPath refers to the to-be generated json file
	// TODO: Add context to fetcher functions
	meta, err := fetchYTService.GetMetadata(url, jsonCaptionsDir)
	if err != nil {
		log.Fatalf("Failed to get video metadata: %v", err)
	}

	// TODO: Should set these as env vars
	vttCaptionsDir := "tmp/captions/"
	vttCaptionsFile := vttCaptionsDir + meta.VideoId + ".en.vtt"

	err = fetchYTService.DownloadCaptions(meta.VideoId, url, vttCaptionsDir)
	if err != nil {
		log.Fatalf("Failed to download captions: %v", err)
	}

	// Call python script to clean captions and export to json
	err = converterService.ConvertVTTToJSON(vttCaptionsFile, meta.CaptionPath)
	if err != nil {
		log.Fatalf("Failed to convert VTT file to JSON: %s", err)
	}

	captions, err := parserService.ParseJSON(meta.CaptionPath)

	if err != nil {
		log.Fatalf("Failed to Parse JSON captions: %s", err)
	}

	err = loaderService.LoadCaptions(meta, captions)
	if err != nil {
		log.Fatalf("Failed to load captions to db: %s", err)
	}

	err = searchService.IndexCaptions(c.Request.Context(), meta, captions)
	if err != nil {
		log.Fatalf("Failed to index captions to Elastic Search: %s", err)
	}
}

// func handleProcessYoutubeVideo(c *gin.Context)
// {

// }
