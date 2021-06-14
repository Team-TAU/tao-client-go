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
