// Package helix handles interactions with twitch via the TAU pass through.
package helix

// Client for interacting with the TAU Helix pass thru, storing the important credential information.
type Client struct {
	hostname string
	port     int
	hasSSL   bool
	token    string
}

// NewClient will generate a new client for interacting with the twitch helix api using the TAU pass thru.  Currently
// never returns an error but could in the future so including for changes to be less likely to be breaking.
func NewClient(hostname string, port int, token string, hasSSL bool) (*Client, error) {
	client := &Client{
		hostname: hostname,
		port:     port,
		hasSSL:   hasSSL,
		token:    token,
	}
	return client, nil
}
