// Package ytdlp provides tools to fetch YouTube captions and related metadata.
package ytdlp

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	cmdutil "banditsecret/internal/pkg/cmdutil"
)

// CaptionMetadata holds metadata about a YouTube video and its captions.
type CaptionMetadata struct {
	VideoId     string
	VideoTitle  string
	Url         string
	CaptionPath string // Path to JSON
}

// Defines the interface to fetch youtube video data
type YTFetcher interface {
	GetMetadata(url, outputPath string) (*CaptionMetadata, error)
	DownloadCaptions(videoId, url, outputDir string) error
}

// Concrete implementation of YTFetcher
type FetchYTService struct {
	executable string
	cmdRunner  cmdutil.CmdRunner
}

// Factory to return a new FetchYTService Service
func NewFetchYTService(executable string, cmdRunner cmdutil.CmdRunner) (*FetchYTService, error) {

	if executable == "" {
		return nil, errors.New("executable cannot be empty")
	}

	return &FetchYTService{
		executable: executable,
		cmdRunner:  cmdRunner,
	}, nil
}

// GetMetadata fetches the video ID and title from a YouTube URL using yt-dlp
func (s *FetchYTService) GetMetadata(url, outputPath string) (*CaptionMetadata, error) {

	if url == "" || outputPath == "" {
		return nil, errors.New("GetMetadata requires a valid url and outputPath")
	}

	fmt.Println(url, outputPath)

	raw, err := s.cmdRunner.Output(
		s.executable,
		"--get-id",
		"--get-title",
		"--no-warnings",
		"--skip-download",
		url)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch video metadata: %w", err)
	}

	parts := bytes.Split(raw, []byte("\n"))
	if len(parts) < 2 {
		return nil, errors.New("unable to parse the metadata")
	}

	// Order matters for yt-dlp as it always outputs title then id
	title := string(bytes.TrimSpace(parts[0]))
	id := string(bytes.TrimSpace(parts[1]))

	metadata := CaptionMetadata{VideoId: id, VideoTitle: title, Url: url, CaptionPath: outputPath + id + ".en.json"}

	return &metadata, nil
}

func (s *FetchYTService) DownloadCaptions(videoId, url, outputDir string) error {

	log.Printf("Attempting to download captions for ID: %s, URL: %s into %s", videoId, url, outputDir)
	cmdOutput, err := s.cmdRunner.CombinedOutput(
		s.executable,
		"--write-subs",
		"--write-auto-subs",
		"--no-warnings",
		"--sub-langs", "en",
		"--skip-download",
		url,
		"-o", outputDir+videoId,
	)

	if err != nil {
		return fmt.Errorf("unable to download video captions using %s: %w\nOutput: %s", s.executable, err, cmdOutput)
	} else {
		log.Println(string(cmdOutput))
	}

	return nil
}
