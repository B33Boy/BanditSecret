package parser

type CaptionMetadata struct {
	VideoId     string
	VideoTitle  string
	Url         string
	CaptionPath string
}

type CaptionParsed struct {
	VideoId string `json:"video_id"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Text    string `json:"text"`
}
