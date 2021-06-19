package gotau

type TAUStreamers struct {
	ID             string `json:"id"`
	TwitchUsername string `json:"twitch_username"`
	TwitchID       string `json:"twitch_id"`
	Streaming      bool   `json:"streaming"`
	Disabled       bool   `json:"disabled"`
	Created        Time   `json:"created"`
	Updated        Time   `json:"updated"`
}

type TAUStream struct {
	ID           string `json:"id"`
	StreamID     string `json:"stream_id"`
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserName     string `json:"user_name"`
	GameID       string `json:"game_id"`
	GameName     string `json:"game_name"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	StartedAt    Time   `json:"started_at"`
	EndedAt      Time   `json:"ended_at"`
	Language     string `json:"language"`
	ThumbnailUrl string `json:"thumbnail_url"`
	TagIDs       string `json:"tag_ids"`
	IsMature     bool   `json:"is_mature"`
}

type TAUTagID string
