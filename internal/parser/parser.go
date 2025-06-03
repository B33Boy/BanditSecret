package parser

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExtractCaptions(url string) error {

	// Run yt-dlp to download vtt file
	captionsFileName, err := downloadCaptions(url)

	fmt.Println("Downloaded file:", captionsFileName)

	if err != nil {
		return fmt.Errorf("ExtractCaptions failed: %w", err)
	}

	// Run scripts/extract_captions to convert to a JSON format
	cmd := exec.Command("python3", "scripts/extract_captions.py", captionsFileName)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run extract_captions.py: %w\nOutput: %s", err, cmdOutput)
	} else {
		fmt.Println("Python output:\n", string(cmdOutput))
	}

	fmt.Println("Done extracting")

	return nil
}

func downloadCaptions(url string) (string, error) {

	videoId, err := getId(url)
	if err != nil {
		return "", err
	}

	captionsFileName := "tmp/captions/" + videoId

	cmd := exec.Command("yt-dlp.exe", "--write-subs", "--write-auto-subs", "--no-warnings", "--sub-langs", "en", "--skip-download", url, "-o", captionsFileName)

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("unable to download video captions: %w\nOutput: %s", err, cmdOutput)
	} else {
		fmt.Println(string(cmdOutput))
	}

	return captionsFileName + ".en.vtt", nil
}

func getId(url string) (string, error) {

	idBytes, err := exec.Command("yt-dlp.exe", "--get-id", "--no-warnings", "--skip-download", url).Output()
	if err != nil {
		return "", fmt.Errorf("unable to fetch video Id")
	}

	return strings.TrimSpace(string(idBytes)), nil
}
