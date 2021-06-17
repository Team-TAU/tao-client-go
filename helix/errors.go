package helix

// AuthorizationError represents an Unauthorized response from Twitch
type AuthorizationError struct {
	err string
}

func (a AuthorizationError) Error() string {
	return a.err
}

// BadRequestError represents bad inputs from an application trying to make an API request to
//twitch based on their documented limitations
type BadRequestError struct {
	err string
}

func (b BadRequestError) Error() string {
	return b.err
}

// GenericError represents a non-specific error, sorta a catch all.
type GenericError struct {
	err  string
	body []byte
	code int
}

func (g GenericError) Error() string {
	return g.err
}

func (g GenericError) Body() string {
	return string(g.body)
}

func (g GenericError) StatusCode() int {
	return g.code
}

// RateLimitError occurs when a twitch rate limits the request.
type RateLimitError struct {
	err string
}

func (r RateLimitError) Error() string {
	return r.err
}
