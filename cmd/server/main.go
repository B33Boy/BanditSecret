package main

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/pkg/captionconverter"
	"banditsecret/internal/pkg/cmdutil"
	"banditsecret/internal/pkg/ytdlp"
	"banditsecret/internal/search"
	"banditsecret/internal/storage"
	"context"
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
		log.Fatalf("Failed to init db: %v", err)
	}
	defer db.Close()

	// Init esclient
	esClient, err := search.InitEsClient()
	if err != nil {
		log.Fatalf("Failed to init Elasticsearch client: %v", err)
	}

	// Get project root using executable location (banditsecret/bin/)
	// TODO: Containerize so we don't need to rely on PYTHON_LOC in venv
	// TODO: Use yt-dlp docker container rather than relying on YTDLP_EXECUTABLE
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	projectRoot := filepath.Dir(filepath.Dir(exePath))

	pythonExecutable := filepath.Join(projectRoot, os.Getenv("PYTHON_LOC"))
	converterScriptPath := filepath.Join(projectRoot, os.Getenv("CONVERTER_SCRIPT_PATH"))

	// Initialize all services
	cmdRunner := cmdutil.NewDefaultCmdRunner()
	fetchYTService, err = ytdlp.NewFetchYTService(os.Getenv("YTDLP_EXECUTABLE"), cmdRunner)
	if err != nil {
		log.Fatalf("NewFetchYTService failed: %w", err)
	}

	converterService, err = captionconverter.NewConverterService(pythonExecutable, converterScriptPath, cmdRunner)
	if err != nil {
		log.Fatalf("NewConverterService failed: %w", err)
	}

	parserService = parser.NewParserService()

	loaderService = storage.NewLoaderService(db)

	searchService = search.NewSearchService(esClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchService.CreateIndex(ctx, os.Getenv("CAPTIONS_INDEX"))

	// set up API server
	router := gin.Default()
	router.POST("/captions", ingestVideo)
	router.Run("localhost:6969")
}

func ingestVideo(c *gin.Context) {
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
	meta, err := fetchYTService.GetMetadata(url, os.Getenv("JSON_CAPTIONS_DIR"))
	if err != nil {
		log.Fatalf("Failed to get video metadata: %v", err)
	}

	vttCaptionsDir := os.Getenv("VTT_CAPTIONS_DIR")
	vttCaptionsFile, err := fetchYTService.DownloadCaptions(meta.VideoId, url, vttCaptionsDir)
	if err != nil {
		log.Fatalf("Failed to download captions: %v", err)
	}

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
