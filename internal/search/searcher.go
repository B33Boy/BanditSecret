package searcher

import (
	"banditsecret/internal/parser"
	"banditsecret/internal/pkg/ytdlp"
	"context"
	// Typed Client
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

func (s *CaptionSearchService) CreateIndex(ctx context.Context, index string) error {
	return s.se.CreateIndex(ctx, index)
}

func (s *CaptionSearchService) IndexCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error {
	return s.se.IndexCaptions(ctx, meta, captions)
}

func (s *CaptionSearchService) SearchCaptions(ctx context.Context, index string, query string) ([]map[string]any, error) {
	return s.se.SearchCaptions(ctx, index, query)
}
