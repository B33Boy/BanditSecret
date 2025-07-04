package searcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/cenkalti/backoff/v5"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type CaptionSearchRepository interface {
	CreateIndex(ctx context.Context, index string) error
	IndexCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error
	SearchCaptions(ctx context.Context, index string, query string) ([]map[string]any, error)
}

type ElasticCaptionSearchRepository struct {
	se *es.TypedClient
}

func NewElasticSearchRepository(se *es.TypedClient) *ElasticCaptionSearchRepository {
	return &ElasticCaptionSearchRepository{
		se: se,
	}
}

func (s *ElasticCaptionSearchRepository) CreateIndex(ctx context.Context, index string) error {
	_, err := s.se.Indices.
		Create(index).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to create index %v: %v", index, err)
	}
	return nil
}

func (s *ElasticCaptionSearchRepository) IndexCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error {

	log.Println("Inserting Captions into ElasticSearch")
	if s.se == nil {
		return errors.New("elasticsearch client is not initialized")
	}

	index_name := os.Getenv("CAPTIONS_INDEX")

	// Create bulk indexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index_name,
		Client:        s.se,
		NumWorkers:    runtime.NumCPU(), // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		return fmt.Errorf("error creating the indexer: %s", err)
	}
	defer bi.Close(ctx)

	// Add captions to bulk indexer
	for _, caption := range captions {
		doc := map[string]any{
			"VideoId":    meta.VideoId,
			"VideoTitle": meta.VideoTitle,
			"Url":        meta.Url,
			"Start":      caption.Start,
			"End":        caption.End,
			"Text":       caption.Text,
		}
		docJson, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal caption document to JSON: %w", err)
		}

		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: fmt.Sprintf("%s_%d", meta.VideoId, caption.Start),
			Body:       bytes.NewReader(docJson),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				// atomic.AddUint64(&countSuccessful, 1)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: Failed to index doc %s for video %s: %s", item.DocumentID, meta.VideoId, err)
				} else {
					log.Printf("ERROR: Failed to index doc %s for video %s: %s", item.DocumentID, meta.VideoId, res.Error.Reason)
				}
			},
		})

		if err != nil {
			return fmt.Errorf("failed to add caption to bulk indexer")
		}
	}

	biStats := bi.Stats()
	if biStats.NumFailed > 0 {
		return fmt.Errorf(
			"bulk indexing had %d failures out of %d actions", biStats.NumFailed, biStats.NumAdded,
		)
	}

	log.Printf("Successfully indexed %d captions for video %s into Elasticsearch index %s", len(captions), meta.VideoId, index_name)

	return nil
}

func (s *ElasticCaptionSearchRepository) SearchCaptions(ctx context.Context, index string, query string) ([]map[string]any, error) {

	res, err := s.se.
		Search().
		Index(index).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"Text": {
						Query: query,
					},
				},
			},
		}).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("search failed: %s", err)
	}

	var results []map[string]any

	for _, hit := range res.Hits.Hits {

		// var entry CaptionEntry
		var entry map[string]any
		err = json.Unmarshal(hit.Source_, &entry)

		if err != nil {
			return nil, errors.New("unmarshal to struct failed")
		}

		results = append(results, entry)
	}

	return results, nil
}

func InitEsClient() (*es.TypedClient, error) {
	timeout := 40 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create config, and include third party backoff to control retry delays
	retryBackoff := backoff.NewExponentialBackOff()

	log.Printf("Attempting to initialize ElasticSearch Client For instance runnning at %s, %s, \n", os.Getenv("ES_HOST"), os.Getenv("ES_PORT"))

	cfg := es.Config{
		Addresses:     []string{fmt.Sprintf("http://%s:%s", os.Getenv("ES_HOST"), os.Getenv("ES_PORT"))},
		RetryOnStatus: []int{502, 503, 504, 429}, // Retry on 429 TooManyRequests statuses
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			d := retryBackoff.NextBackOff()
			log.Printf("Retry #%d: sleeping for %s", i, d)
			return d
		},
		MaxRetries: 10,
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
	log.Printf("Connected to ElasticSearch cluster: %s", res.ClusterName)

	return esClient, nil
}
