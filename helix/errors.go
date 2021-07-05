package helix

import "time"

// RateLimitError occurs when a twitch rate limits the request.
type RateLimitError struct {
	err   string
	reset *time.Time
}

func (r RateLimitError) Error() string {
	return r.err
}

// ResetTime can be used to check when the point bucket used by twitch will refill.
// See https://dev.twitch.tv/docs/api/guide#rate-limits
func (r RateLimitError) ResetTime() *time.Time {
	return r.reset
}
