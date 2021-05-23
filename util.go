package gotau

import (
	"strings"
	"time"
)

// Time wraps the standard time.Time to allow for custom parsing from json
type Time struct {
	time.Time
}

// UnmarshalJSON implemented to allow for parsing this time object from TAU
func (t *Time) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02T15:04:05.999999999-07:00"
	timeAsString := strings.TrimSpace(string(b))
	timeAsString = strings.Trim(timeAsString, "\"")
	timestamp, err := time.Parse(layout, timeAsString)
	if err != nil {
		return err
	}
	t.Time = timestamp
	return nil
}
