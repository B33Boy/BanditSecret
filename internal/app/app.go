package app

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/pkg/captionconverter"
	"banditsecret/internal/pkg/cmdutil"
	"banditsecret/internal/pkg/ytdlp"
	searcher "banditsecret/internal/search"
	"banditsecret/internal/storage"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
)

type SQLCaptionRepository = storage.SQLCaptionRepository

type ApplicationServices struct {
	Fetcher   ytdlp.YTFetcher
	Converter captionconverter.Converter
	Parser    parser.Parser
	Loader    storage.Loader
	Searcher  searcher.Searcher
}

func NewApplicationServices(db storage.CaptionRepository, esClient *es.TypedClient) (*ApplicationServices, error) {

	// Get project root using executable location (banditsecret/bin/)
	// TODO: Containerize so we don't need to rely on PYTHON_LOC in venv
	// TODO: Use yt-dlp docker container rather than relying on YTDLP_EXECUTABLE
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	projectRoot := filepath.Dir(exePath)
	log.Println(projectRoot)

	// TODO: refactor, we don't need to pass this data in anymore, and we don't need projectRoot
	pythonExecutable := os.Getenv("PYTHON_LOC")
	converterScriptPath := os.Getenv("CONVERTER_SCRIPT_PATH")

	// Initialize all services
	cmdRunner := cmdutil.NewDefaultCmdRunner()

	fetchYTService, err := ytdlp.NewFetchYTService(os.Getenv("YTDLP_EXECUTABLE"), cmdRunner)
	if err != nil {
		return nil, fmt.Errorf("NewFetchYTService failed: %w", err)
	}

	converterService, err := captionconverter.NewConverterService(pythonExecutable, converterScriptPath, cmdRunner)
	if err != nil {
		return nil, fmt.Errorf("NewConverterService failed: %w", err)
	}

	parserService := parser.NewParserService()
	loaderService := storage.NewLoaderService(db)
	searchService := searcher.NewSearcherService(esClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchService.CreateIndex(ctx, os.Getenv("CAPTIONS_INDEX"))

	return &ApplicationServices{
		Fetcher:   fetchYTService,
		Converter: converterService,
		Parser:    parserService,
		Loader:    loaderService,
		Searcher:  searchService,
	}, nil
}
