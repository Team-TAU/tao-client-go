package gotau

// AuthorizationError represents an Unauthorized response from Twitch
type AuthorizationError struct {
	// Err is the error given to the user
	Err string
}

func (a AuthorizationError) Error() string {
	return a.Err
}

// BadRequestError represents bad inputs from an application trying to make an API request to
//twitch based on their documented limitations
type BadRequestError struct {
	// Err is the error given to the user
	Err string
}

func (b BadRequestError) Error() string {
	return b.Err
}

// GenericError represents a non-specific error, sorta a catch all.
type GenericError struct {
	// Err is the error given to the user
	Err string
	// Body the body of the web response
	Body []byte
	// Code the http status code from the response
	Code int
}

func (g GenericError) Error() string {
	return g.Err
}
