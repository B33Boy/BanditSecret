package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	parser "banditsecret/internal/parser"
	fetcher "banditsecret/internal/pkg/ytdlp"

	"github.com/go-sql-driver/mysql"
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

func InitDb() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database!")

	return db, nil
}
