// Package ytdlp provides tools to fetch YouTube captions and related metadata.
package ytdlp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	cmdutil "banditsecret/internal/pkg/cmdutil"
)

// CaptionMetadata holds metadata about a YouTube video and its captions.
type CaptionMetadata struct {
	VideoId     string
	VideoTitle  string
	Url         string
	CaptionPath string // Path to JSON
}

type MetadataResp struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type CaptionsReq struct {
	Url       string `json:"url"`
	OutputDir string `json:"output_dir"`
}

// Defines the interface to fetch youtube video data
type YTFetcher interface {
	GetMetadata(url, outputPath string) (*CaptionMetadata, error)
	DownloadCaptions(videoId, url, outputDir string) (string, error)
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

	// make http request to ytdlp container
	reqUrl := fmt.Sprintf("http://%s:%s/get_metadata?url=%s", os.Getenv("YTDLP_HOST"), os.Getenv("YTDLP_PORT"), url)
	resp, err := http.Get(reqUrl)

	if err != nil {
		return nil, fmt.Errorf("unable to get a valid response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var parsedResp MetadataResp
	json.Unmarshal(body, &parsedResp)

	metadata := CaptionMetadata{VideoId: parsedResp.Id, VideoTitle: parsedResp.Title, Url: url, CaptionPath: outputPath + parsedResp.Id + ".en.json"}

	return &metadata, nil
}

func (s *FetchYTService) DownloadCaptions(videoId, url, outputDir string) (string, error) {

	log.Printf("Attempting to download captions for URL: %s into %s", url, outputDir)

	vttCaptionsFile := outputDir + videoId + ".en.vtt"

	if cmdutil.FileExists(vttCaptionsFile) {
		log.Printf("Vtt file %s already exists! Skipping download\n", vttCaptionsFile)
		return vttCaptionsFile, nil
	}

	captionsReq := CaptionsReq{
		Url:       url,
		OutputDir: outputDir + videoId,
	}

	reqBytes, err := json.Marshal(captionsReq)
	if err != nil {
		return "", fmt.Errorf("unable to convert request struct to bytes: %w", err)
	}

	reqUrl := fmt.Sprintf("http://%s:%s/get_captions", os.Getenv("YTDLP_HOST"), os.Getenv("YTDLP_PORT"))
	resp, err := http.Post(reqUrl, "application/json; charset=utf-8", bytes.NewReader(reqBytes))
	if err != nil {
		return "", fmt.Errorf("unable to get a valid response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response: %w", err)
	}

	// if err != nil {
	// 	return "", fmt.Errorf("unable to download video captions using %s: %w\nOutput: %s", s.executable, err, string(body))
	// }

	log.Println(string(body))
	log.Printf("Downloaded vtt file for videoId: %s", videoId)

	return vttCaptionsFile, nil
}
