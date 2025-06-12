package storage

import (
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
	LoadCaptions(meta *CaptionMetadata, captions []CaptionEntry) error
}

type LoaderService struct {
	db *sql.DB
}

func NewLoaderService(db *sql.DB) *LoaderService {
	return &LoaderService{
		db: db,
	}
}

func (s *LoaderService) LoadCaptions(meta *CaptionMetadata, captions []CaptionEntry) error {
	tx, err := s.db.Begin()
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
	err = s.upsertVideoMetadata(tx, meta)
	if err != nil {
		return err
	}

	// 2. Clear (DELETE) existing captions for this video
	err = s.deleteExistingCaptions(tx, meta.VideoId)
	if err != nil {
		return err
	}

	// 3. Insert new captions (BATCH INSERT)
	err = s.insertNewCaptions(tx, captions)

	return nil
}

func (s *LoaderService) upsertVideoMetadata(tx *sql.Tx, meta *CaptionMetadata) error {
	upsertVideoSql := `INSERT INTO Videos (Id, Title, VideoUrl) 
							VALUES (?, ?, ?)
							ON DUPLICATE KEY UPDATE
							Title = VALUES(Title),
							VideoUrl = VALUES(VideoUrl);`

	_, err := tx.Exec(upsertVideoSql, meta.VideoId, meta.VideoTitle, meta.Url)

	if err != nil {
		return fmt.Errorf("failed to upsert video metadata for %s: %w", meta.VideoId, err)
	}
	log.Printf("Upserted video metadata for %s", meta.VideoId)
	return nil
}

func (s *LoaderService) deleteExistingCaptions(tx *sql.Tx, videoId string) error {

	deleteCaptionsSql := `DELETE FROM Captions WHERE VideoId = ?;`

	_, err := tx.Exec(deleteCaptionsSql, videoId)

	if err != nil {
		return fmt.Errorf("failed to delete existing captions for video %s: %w", videoId, err)
	}
	log.Printf("Deleted existing captions for %s", videoId)
	return nil
}

func (s *LoaderService) insertNewCaptions(tx *sql.Tx, captions []CaptionEntry) error {

	insertCaptionsSQL := `INSERT INTO Captions (VideoId, StartTime, EndTime, CaptionText)
				   		VALUES (?, ?, ?, ?);`

	st, err := tx.Prepare(insertCaptionsSQL)

	if err != nil {
		return fmt.Errorf("failed to prepare statement for captions: %w", err)
	}

	defer st.Close()

	for i, caption := range captions {
		_, err := st.Exec(caption.VideoId, caption.Start, caption.End, caption.Text)
		if err != nil {
			return fmt.Errorf("failed to insert caption %d for video %s: %w", i, caption.VideoId, err)
		}
	}
	log.Printf("Inserted %d new captions for video %s", len(captions), captions[0].VideoId)
	return nil
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

	fmt.Println("Connected to database!")

	return db, nil
}

// ==========================================================================================================

// func LoadCaptionsFromJson(filepath string) ([]CaptionEntry, error) {

// 	// TODO: Read chunk by chunk for larger files
// 	// Open JSON file and insert into DB
// 	// slurp entire json file into memory

// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return nil, fmt.Errorf("LoadCaptions failed to open file: %w", err)
// 	}
// 	defer file.Close()

// 	bytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("LoadCaptions failed to read file: %w", err)
// 	}

// 	// captions, err := parseCaptions(bytes)
// 	var captions []CaptionEntry
// 	err = json.Unmarshal(bytes, &captions)
// 	if err != nil {
// 		return nil, fmt.Errorf("LoadCaptions failed to parse JSON file: %w", err)
// 	}

// 	return captions, nil
// }

// func StoreVideoInfoToDb(db *sql.DB, metadata *CaptionMetadata) {

// 	// Populate Videos Table
// 	_, err := db.Exec("INSERT INTO Videos (Id, Title, VideoUrl) VALUES (?, ?, ?)", metadata.VideoId, metadata.VideoTitle, metadata.Url)

// 	if err != nil {
// 		log.Printf("StoreVideoInfoToDB failed: %v", err)
// 	} else {
// 		log.Printf("StoreVideoInfoToDB succeeded for Video Id: %v", metadata.VideoId)
// 	}
// }

// func StoreCaptionsToDb(db *sql.DB, captions []CaptionEntry) {

// 	idToSearch := captions[0].VideoId
// 	var count int

// 	err := db.QueryRow(`SELECT COUNT(*) FROM Captions WHERE VideoId = ?`, idToSearch).Scan(&count)
// 	if err != nil {
// 		log.Printf("StoreCaptionsToDb failed: %v", err)
// 		return
// 	}

// 	if int(count) > 0 {
// 		log.Printf("StoreCaptionsToDb, captions for video %v already exist in db", idToSearch)
// 		return
// 	}

// 	// Populate Captions Table
// 	for _, caption := range captions {
// 		err := addCaptionEntry(db, caption)
// 		if err != nil {
// 			log.Printf("StoreCaptionsToDb, caption entry failed: %v", err)
// 			continue
// 		}
// 	}

// 	log.Printf("StoreCaptionsToDb succeeded for Video Id: %v", captions[0].VideoId)
// }

// func addCaptionEntry(db *sql.DB, caption CaptionEntry) error {

// 	_, err := db.Exec("INSERT INTO Captions (VideoId, StartTime, EndTime, CaptionText) VALUES (?, ?, ?, ?)", caption.VideoId, caption.Start, caption.End, caption.Text)
// 	return err
// }
