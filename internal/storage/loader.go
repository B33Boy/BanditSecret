package storage

import (
	"context"

	parser "banditsecret/internal/parser"
	fetcher "banditsecret/internal/pkg/ytdlp"
)

// ==========================================================================================================

type CaptionMetadata = fetcher.CaptionMetadata
type CaptionEntry = parser.CaptionEntry

type Loader interface {
	LoadCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error
}

type LoaderService struct {
	repo CaptionRepository
}

func NewLoaderService(repo CaptionRepository) *LoaderService {
	return &LoaderService{
		repo: repo,
	}
}

func (s *LoaderService) LoadCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error {
	return s.repo.SaveCaptions(ctx, meta, captions)
}
