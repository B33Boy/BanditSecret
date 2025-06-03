package parser

import (
	"fmt"
	"os/exec"
	"strings"
)

type CaptionResult struct {
	VideoId     string
	CaptionPath string
}

func ExtractCaptions(url string) (*CaptionResult, error) {

	// Run yt-dlp to download vtt file
	videoId, err := downloadCaptions(url)

	fmt.Println("Downloaded vtt file for Id:", videoId)

	if err != nil {
		return nil, fmt.Errorf("ExtractCaptions failed: %w", err)
	}

	// Run scripts/extract_captions to convert to a JSON format
	cmd := exec.Command("./venv/Scripts/python", "scripts/extract_captions.py", getFileNameFromId(videoId))
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run extract_captions.py: %w\nOutput: %s", err, cmdOutput)
	}

	captionResult := CaptionResult{VideoId: videoId, CaptionPath: "tmp/captions_parsed/" + videoId}
	fmt.Println("Done extracting")

	return &captionResult, nil
}

func downloadCaptions(url string) (string, error) {

	videoId, err := getId(url)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("yt-dlp.exe", "--write-subs", "--write-auto-subs", "--no-warnings", "--sub-langs", "en", "--skip-download", url, "-o", "tmp/captions/"+videoId)

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("unable to download video captions: %w\nOutput: %s", err, cmdOutput)
	} else {
		fmt.Println(string(cmdOutput))
	}

	return videoId, nil
}

func getId(url string) (string, error) {

	idBytes, err := exec.Command("yt-dlp.exe", "--get-id", "--no-warnings", "--skip-download", url).Output()
	if err != nil {
		return "", fmt.Errorf("unable to fetch video Id")
	}

	return strings.TrimSpace(string(idBytes)), nil
}

func getFileNameFromId(id string) string {
	return "tmp/captions/" + id + ".en.vtt"
}
