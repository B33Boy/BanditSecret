package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	parser "banditsecret/internal/parser"

	"github.com/go-sql-driver/mysql"
)

type CaptionMetadata = parser.CaptionMetadata
type CaptionParsed = parser.CaptionParsed

func LoadCaptionsFromJson(filepath string) ([]CaptionParsed, error) {

	// TODO: Read chunk by chunk for larger files
	// Open JSON file and insert into DB
	// slurp entire json file into memory

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("LoadCaptions failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("LoadCaptions failed to read file: %w", err)
	}

	// captions, err := parseCaptions(bytes)
	var captions []CaptionParsed
	err = json.Unmarshal(bytes, &captions)
	if err != nil {
		return nil, fmt.Errorf("LoadCaptions failed to parse JSON file: %w", err)
	}

	return captions, nil
}

func InitDb() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	fmt.Println("Connected to DB!")

	return db, nil
}

func StoreVideoInfoToDb(db *sql.DB, metadata *CaptionMetadata) {

	// Populate Videos Table
	_, err := db.Exec("INSERT INTO Videos (Id, Title, VideoUrl) VALUES (?, ?, ?)", metadata.VideoId, metadata.VideoTitle, metadata.Url)

	if err != nil {
		log.Printf("StoreVideoInfoToDB failed: %v", err)
	} else {
		log.Printf("StoreVideoInfoToDB succeeded for Video Id: %v", metadata.VideoId)
	}
}

func StoreCaptionsToDb(db *sql.DB, captions []CaptionParsed) {

	// Populate Captions Table
	for _, caption := range captions {
		err := addCaptionEntry(db, caption)
		if err != nil {
			log.Printf("StoreCaptionsToDb, caption entry failed: %v", err)
			continue
		}
	}

	log.Printf("StoreCaptionsToDb succeeded for Video Id: %v", captions[0].VideoId)
}

func addCaptionEntry(db *sql.DB, caption CaptionParsed) error {

	_, err := db.Exec("INSERT INTO Captions (VideoId, StartTime, EndTime, CaptionText) VALUES (?, ?, ?, ?)", caption.VideoId, caption.Start, caption.End, caption.Text)
	return err
}
