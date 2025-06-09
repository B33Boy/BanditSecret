package parser

type TimeMs uint32

type CaptionMetadata struct {
	VideoId     string
	VideoTitle  string
	Url         string
	CaptionPath string
}

type CaptionEntry struct {
	VideoId string `json:"video_id"`
	Start   TimeMs `json:"start"`
	End     TimeMs `json:"end"`
	Text    string `json:"text"`
}
