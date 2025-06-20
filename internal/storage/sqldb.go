package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/go-sql-driver/mysql"
)

type CaptionRepository interface {
	SaveCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error
}

type SQLCaptionRepository struct {
	db *sql.DB
}

func NewSQLCaptionRepository(db *sql.DB) *SQLCaptionRepository {
	return &SQLCaptionRepository{
		db: db,
	}
}

func (s *SQLCaptionRepository) SaveCaptions(ctx context.Context, meta *CaptionMetadata, captions []CaptionEntry) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction %w", err)
	}

	defer func() {
		// Catch any unexpected panics, rollback our transaction, and re-throw panic
		r := recover()
		if r != nil {
			log.Printf("Recovered from panic during transaction for video %s: %v. Rolling back.", meta.VideoId, r)
			tx.Rollback()
			panic(r) // throw panic again
		} else if err != nil {
			log.Printf("Error during transaction for video %s: %v. Rolling back.", meta.VideoId, err)
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Failed to commit transaction for video %s: %v", meta.VideoId, err)
			}
		}
	}()

	// 1. Store/Update (UPSERT) video metadata
	err = s.upsertVideoMetadata(ctx, tx, meta)
	if err != nil {
		return err
	}

	// 2. Clear (DELETE) existing captions for this video
	err = s.deleteExistingCaptions(ctx, tx, meta.VideoId)
	if err != nil {
		return err
	}

	// 3. Insert new captions (BATCH INSERT)
	err = s.insertNewCaptions(ctx, tx, captions)

	return nil
}

func (s *SQLCaptionRepository) upsertVideoMetadata(ctx context.Context, tx *sql.Tx, meta *CaptionMetadata) error {
	upsertVideoSql := `INSERT INTO Videos (Id, Title, VideoUrl) 
							VALUES (?, ?, ?)
							ON DUPLICATE KEY UPDATE
							Title = VALUES(Title),
							VideoUrl = VALUES(VideoUrl);`

	_, err := tx.ExecContext(ctx, upsertVideoSql, meta.VideoId, meta.VideoTitle, meta.Url)

	if err != nil {
		return fmt.Errorf("failed to upsert video metadata for %s: %w", meta.VideoId, err)
	}
	log.Printf("Upserted video metadata for %s", meta.VideoId)
	return nil
}

func (s *SQLCaptionRepository) deleteExistingCaptions(ctx context.Context, tx *sql.Tx, videoId string) error {

	deleteCaptionsSql := `DELETE FROM Captions WHERE VideoId = ?;`

	_, err := tx.ExecContext(ctx, deleteCaptionsSql, videoId)

	if err != nil {
		return fmt.Errorf("failed to delete existing captions for video %s: %w", videoId, err)
	}
	log.Printf("Deleted existing captions for %s", videoId)
	return nil
}

func (s *SQLCaptionRepository) insertNewCaptions(ctx context.Context, tx *sql.Tx, captions []CaptionEntry) error {

	insertCaptionsSQL := `INSERT INTO Captions (VideoId, StartTime, EndTime, CaptionText)
				   		VALUES (?, ?, ?, ?);`

	st, err := tx.PrepareContext(ctx, insertCaptionsSQL)

	if err != nil {
		return fmt.Errorf("failed to prepare statement for captions: %w", err)
	}

	defer st.Close()

	for i, caption := range captions {
		_, err := st.ExecContext(ctx, caption.VideoId, caption.Start, caption.End, caption.Text)
		if err != nil {
			return fmt.Errorf("failed to insert caption %d for video %s: %w", i, caption.VideoId, err)
		}
	}
	log.Printf("Inserted %d new captions for video %s", len(captions), captions[0].VideoId)
	return nil
}

func InitDb() (*sql.DB, error) {

	timeout := 40 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cfg := getDbConfig()

	notify := makeNotify()

	db, err := backoff.Retry(ctx, func() (*sql.DB, error) {
		return connectToDb(cfg)
	},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxElapsedTime(timeout),
		backoff.WithNotify(notify))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to sql database! %v", err)
	}

	log.Println("Connected to database!")
	return db, nil
}

func makeNotify() func(error, time.Duration) {
	retryAttempt := 0
	return func(err error, d time.Duration) {
		retryAttempt++
		log.Printf("Error: %v\nRetry #%d: sleeping for %.1fs", err, retryAttempt, d.Seconds())
	}
}

func getDbConfig() *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.ParseTime = true
	return cfg
}

func connectToDb(cfg *mysql.Config) (*sql.DB, error) {

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

	return db, nil
}
