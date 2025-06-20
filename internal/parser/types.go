package parser

import (
	"strings"
	"time"
)

type TimeMs uint32
type CaptionEntry struct {
	VideoId string `json:"video_id"`
	Start   TimeMs `json:"start"`
	End     TimeMs `json:"end"`
	Text    string `json:"text"`
}

func (t *TimeMs) UnmarshalJSON(data []byte) error {

	s := strings.Trim(string(data), `"`)

	parsedTime, err := time.Parse("15:04:05.000", s)
	if err != nil {
		return err
	}

	hoursToMsecs := parsedTime.Hour() * 60 * 60 * 1000
	minToMsecs := parsedTime.Minute() * 60 * 1000
	secToMsecs := parsedTime.Second() * 1000
	msecs := parsedTime.Nanosecond() / 1e6

	*t = TimeMs(hoursToMsecs + minToMsecs + secToMsecs + msecs)

	return nil
}
