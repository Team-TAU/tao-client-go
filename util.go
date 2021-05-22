package go_tau

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

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

//func (e *Event) CreatedAsTime() (time.Time, error) {
//	// 2021-05-22T05:20:06.120452+00:00
//	return time.Parse("2006-01-02T15:04:05.999999999-07:00", e.Created)
//}
