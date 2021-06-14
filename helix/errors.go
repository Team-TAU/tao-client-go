package helix

type AuthorizationError struct {
	err string
}

func (a AuthorizationError) Error() string {
	return a.err
}

type BadRequestError struct {
	err string
}

func (b BadRequestError) Error() string {
	return b.err
}
