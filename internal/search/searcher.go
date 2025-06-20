package searcher

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/pkg/ytdlp"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v5"
	es "github.com/elastic/go-elasticsearch/v9" // Typed Client
)

type CaptionEntry = parser.CaptionEntry
type CaptionMetadata = ytdlp.CaptionMetadata

type Searcher interface {
	CreateIndex(ctx context.Context, index string) error
	IndexCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error
	SearchCaptions(ctx context.Context, index string, query string) ([]map[string]any, error)
}

type CaptionSearchService struct {
	se CaptionSearchRepository
}

func NewSearcherService(se CaptionSearchRepository) *CaptionSearchService {
	return &CaptionSearchService{
		se: se,
	}
}

func InitEsClient() (*es.TypedClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create config, and include third party backoff to control retry delays
	retryBackoff := backoff.NewExponentialBackOff()

	log.Printf("Initializing ElasticSearch Client For instance runnning at %s, %s, \n", os.Getenv("ES_HOST"), os.Getenv("ES_PORT"))
	cfg := es.Config{
		Addresses:     []string{fmt.Sprintf("http://%s:%s", os.Getenv("ES_HOST"), os.Getenv("ES_PORT"))},
		RetryOnStatus: []int{502, 503, 504, 429}, // Retry on 429 TooManyRequests statuses
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}

	// Create client from config
	esClient, err := es.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Retrieve metadata from client
	res, err := esClient.Info().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting response from Elasticsearch: %w", err)
	}
	log.Printf("Cluster name: %s", res.ClusterName)

	return esClient, nil
}

func (s *CaptionSearchService) CreateIndex(ctx context.Context, index string) error {
	return s.se.CreateIndex(ctx, index)
}

func (s *CaptionSearchService) IndexCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error {
	return s.se.IndexCaptions(ctx, meta, captions)
}

func (s *CaptionSearchService) SearchCaptions(ctx context.Context, index string, query string) ([]map[string]any, error) {
	return s.se.SearchCaptions(ctx, index, query)
}
