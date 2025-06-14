package app

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

	"database/sql"

	es "github.com/elastic/go-elasticsearch/v9"
)

type ApplicationServices struct {
	Fetcher   ytdlp.YTFetcher
	Converter captionconverter.Converter
	Parser    parser.Parser
	Loader    storage.Loader
	Searcher  search.Searcher
}

func InitConnections() (*sql.DB, *es.Client, error) {
	// Init db connection
	db, err := storage.InitDb()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Init esclient
	esClient, err := search.InitEsClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init Elasticsearch client: %w", err)
	}
	return db, esClient, nil
}

func NewApplicationServices(db *sql.DB, esClient *es.Client) (*ApplicationServices, error) {

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
	searchService := search.NewSearchService(esClient)

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
