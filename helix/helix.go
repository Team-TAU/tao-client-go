// Package helix handles interactions with twitch via the TAU pass through.
package helix

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

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

func (c *Client) helixRequest(endpoint string, params map[string][]string, body []byte, method string) ([]byte, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest(method, endpointURL, buffer)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	request.Header.Add("Content-Type", "application/json")
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	for key, values := range params {
		for _, item := range values {
			q.Add(key, item)
		}
	}
	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		body, err := ioutil.ReadAll(response.Body)
		return body, err
	}
	if response.StatusCode == 401 {
		return nil, AuthorizationError{}
	} else if response.StatusCode == 429 {
		resetEpoch := response.Header.Get("Ratelimit-Reset")
		rlErr := RateLimitError{
			err: "rate limited: received http 429",
		}
		if resetEpoch != "" {
			epoch, err := strconv.ParseInt(resetEpoch, 10, 64)
			if err != nil {
				return nil, err
			}
			reset := time.Unix(epoch, 0)
			rlErr.reset = &reset
		}
		return nil, rlErr
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		err = GenericError{
			err:  fmt.Sprintf("response code %d: %s", response.StatusCode, body),
			body: body,
			code: response.StatusCode,
		}
		return nil, err
	}
}
