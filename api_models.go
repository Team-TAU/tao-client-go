package gotau

import (
	"encoding/json"
	"strings"
)

// TAUStreamer represents a streamer as TAU returns them from it's database.
type TAUStreamer struct {
	ID             string `json:"id"`
	TwitchUsername string `json:"twitch_username"`
	TwitchID       string `json:"twitch_id"`
	Streaming      bool   `json:"streaming"`
	Disabled       bool   `json:"disabled"`
	Created        Time   `json:"created"`
	Updated        Time   `json:"updated"`
}

// TAUStream represents a stream as TAU returns them from it's database.
type TAUStream struct {
	ID           string  `json:"id"`
	StreamID     string  `json:"stream_id"`
	UserID       string  `json:"user_id"`
	UserLogin    string  `json:"user_login"`
	UserName     string  `json:"user_name"`
	GameID       string  `json:"game_id"`
	GameName     string  `json:"game_name"`
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	ViewerCount  int     `json:"viewer_count"`
	StartedAt    Time    `json:"started_at"`
	EndedAt      Time    `json:"ended_at"`
	Language     string  `json:"language"`
	ThumbnailUrl string  `json:"thumbnail_url"`
	TagIDs       TAUTags `json:"tag_ids"`
	IsMature     bool    `json:"is_mature"`
}

//TAUTags is a list of strings containing tags from a stream
type TAUTags []string

// UnmarshalJSON handles unmarshalling the json data from a less than useful form.
func (t *TAUTags) UnmarshalJSON(b []byte) error {
	var stringVal string
	err := json.Unmarshal(b, &stringVal)
	if err != nil {
		return err
	}
	stringVal = strings.TrimPrefix(stringVal, "[")
	stringVal = strings.TrimSuffix(stringVal, "]")
	stringVal = strings.Replace(stringVal, "'", "", -1)
	stringVal = strings.Replace(stringVal, " ", "", -1)
	*t = strings.Split(stringVal, ",")
	return nil
}
