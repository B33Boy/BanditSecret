package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parser interface {
	ParseJSON(jsonFile string) ([]CaptionEntry, error)
}

type ParserService struct{}

func NewParserService() *ParserService {
	return &ParserService{}
}

func (s *ParserService) ParseJSON(jsonFile string) ([]CaptionEntry, error) {

	_, err := os.Stat(jsonFile)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("JSON caption file not found at %s: %w", jsonFile, err)
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file %s: %w", jsonFile, err)
	}

	var captions []CaptionEntry

	err = json.Unmarshal(data, &captions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", jsonFile, err)
	}
	return captions, nil
}

// =================================================================================================
// func fileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	return os.IsNotExist(err)
// }

// func ExtractCaptions(url string, outputDir string) (*CaptionMetadata, error) {

// 	// Run yt-dlp to download vtt file
// 	captionMetadata, err := getMetadata(url, outputDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("getMetadata failed: %w", err)
// 	}

// 	// Don't download if captions file already exists
// 	vttCaptionsDir := "tmp/captions/"
// 	vttCaptionsFile := vttCaptionsDir + captionMetadata.VideoId + ".en.vtt"

// 	if !fileExists(vttCaptionsFile) {
// 		err = getCaptions(vttCaptionsDir, captionMetadata)
// 		if err != nil {
// 			return nil, fmt.Errorf("getCaptions failed: %w", err)
// 		}

// 		log.Println("Downloaded vtt file for Id:", captionMetadata.VideoId)
// 	} else {
// 		log.Printf("Vtt file %s already exists! Skipping download\n", vttCaptionsFile)
// 	}

// 	// Run scripts/extract_captions.py to convert to a JSON format
// 	cmd := exec.Command("./venv/Scripts/python", "scripts/extract_captions.py", vttCaptionsFile)
// 	cmdOutput, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to run extract_captions.py: %w\nOutput: %s", err, cmdOutput)
// 	}

// 	log.Println("Done extracting")

// 	return captionMetadata, nil
// }

// func getCaptions(vttCaptionsDir string, metadata *CaptionMetadata) error {

// 	cmd := exec.Command("yt-dlp.exe", "--write-subs", "--write-auto-subs", "--no-warnings", "--sub-langs", "en", "--skip-download", metadata.Url, "-o", vttCaptionsDir+metadata.VideoId)

// 	cmdOutput, err := cmd.CombinedOutput()

// 	if err != nil {
// 		return fmt.Errorf("unable to download video captions: %w\nOutput: %s", err, cmdOutput)
// 	} else {
// 		log.Println(string(cmdOutput))
// 	}

// 	return nil
// }

// func getMetadata(url string, path string) (*CaptionMetadata, error) {

// 	raw, err := exec.Command("yt-dlp.exe", "--get-id", "--get-title", "--no-warnings", "--skip-download", url).Output()
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch video metadata")
// 	}

// 	parts := bytes.Split(raw, []byte("\n"))
// 	if len(parts) < 2 {
// 		return nil, fmt.Errorf("unable to parse the metadata")
// 	}

// 	// Order matters for yt-dlp as it always outputs title then id
// 	title := string(bytes.TrimSpace(parts[0]))
// 	id := string(bytes.TrimSpace(parts[1]))

// 	metadata := CaptionMetadata{VideoId: id, VideoTitle: title, Url: url, CaptionPath: path + id + ".en.json"}

// 	return &metadata, nil
// }
