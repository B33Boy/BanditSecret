package parser

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func ExtractCaptions(url string, outputDir string) (*CaptionMetadata, error) {

	// Run yt-dlp to download vtt file
	captionMetadata, err := getMetadata(url, outputDir)

	err = getCaptions(captionMetadata)
	if err != nil {
		return nil, fmt.Errorf("getCaptions failed: %w", err)
	}

	log.Println("Downloaded vtt file for Id:", captionMetadata.VideoId)

	// Run scripts/extract_captions.py to convert to a JSON format
	cmd := exec.Command("./venv/Scripts/python", "scripts/extract_captions.py", getVttFileNameFromId(captionMetadata.VideoId))
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run extract_captions.py: %w\nOutput: %s", err, cmdOutput)
	}

	log.Println("Done extracting")

	return captionMetadata, nil
}

func getCaptions(metadata *CaptionMetadata) error {

	cmd := exec.Command("yt-dlp.exe", "--write-subs", "--write-auto-subs", "--no-warnings", "--sub-langs", "en", "--skip-download", metadata.Url, "-o", metadata.CaptionPath)

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to download video captions: %w\nOutput: %s", err, cmdOutput)
	} else {
		log.Println(string(cmdOutput))
	}

	return nil
}

func getMetadata(url string, path string) (*CaptionMetadata, error) {

	raw, err := exec.Command("yt-dlp.exe", "--get-id", "--get-title", "--no-warnings", "--skip-download", url).Output()
	if err != nil {
		return "", "", fmt.Errorf("unable to fetch video Id")
	}

	parts := bytes.Split(raw, []byte("\n"), 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("unable to download video captions: %w\nOutput: %s", err, cmdOutput)
	}

	id := string(bytes.TrimSpace(parts[0]))
	title := string(bytes.TrimSpace(parts[1]))

	metadata := CaptionMetadata{VideoId: id, VideoTitle: title, Url: url, CaptionPath: path + id + "en.json"}

	return &metadata, nil
}

func getVttFileNameFromId(id string) string {
	return "tmp/captions/" + id + ".en.vtt"
}
