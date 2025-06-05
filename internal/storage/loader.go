package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type ParsedCaption struct {
	Id    string `json:"video_id"`
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

func LoadCaptionsFromJson(filepath string) ([]ParsedCaption, error) {

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

	var captions []ParsedCaption
	if err := json.Unmarshal(bytes, &captions); err != nil {
		return nil, fmt.Errorf("LoadCaptions failed to parse JSON file: %w", err)
	}

	// for _, cap := range captions {
	// 	fmt.Printf("ID: %s, Start: %s, End: %s, Text: %s, \n", cap.Id, cap.Start, cap.End, cap.Text)
	// }

	return captions, nil
}

func StoreCaptionsToDb(captions []ParsedCaption) error {
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("StoreCaptionsToDb failed to connect to DB: %w", err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		return fmt.Errorf("StoreCaptionsToDb failed to ping DB: %w", pingErr)
	}

	fmt.Println("Connected")

	return nil
}
