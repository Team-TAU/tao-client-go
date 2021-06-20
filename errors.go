package gotau

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

// Body returns the body of a failed request
func (g GenericError) Body() string {
	return string(g.body)
}

// StatusCode returns the status code returned by twitch
func (g GenericError) StatusCode() int {
	return g.code
}
