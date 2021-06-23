package gotau

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetStreamers can be used to get a list of all the streamers that TAU is listening for going live alerts.
func (c *Client) GetStreamers() ([]*TAUStreamer, error) {
	body, err := c.apiRequest("streamers", nil, nil, "GET")
	if err != nil {
		return nil, err
	}

	var streamers []*TAUStreamer
	err = json.Unmarshal(body, &streamers)
	if err != nil {
		return nil, err
	}

	return streamers, err
}

// GetLatestStreamForStreamer gets the latest stream for a given streamer
func (c *Client) GetLatestStreamForStreamer(ID string) (*TAUStream, error) {
	ID = strings.TrimSpace(ID)
	if ID == "" {
		return nil, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	body, err := c.apiRequest(fmt.Sprintf("streamers/%s/streams/latest", ID), nil, nil, "GET")
	if err != nil {
		return nil, err
	}

	stream := new(TAUStream)
	err = json.Unmarshal(body, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// FollowStreamerOnTau follows the users and subscribes for notifications when they go live
func (c *Client) FollowStreamerOnTau(username string) (*TAUStreamer, error) {
	type tmp struct {
		Username  string `json:"twitch_username"`
		Streaming bool   `json:"streaming"`
		Disabled  bool   `json:"disabled"`
	}

	username = strings.TrimSpace(username)
	if username == "" {
		return nil, BadRequestError{
			err: "invalid request, username can't be blank",
		}
	}

	data := tmp{
		Username: username,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.apiRequest("streamers", nil, body, "POST")
	if err != nil {
		return nil, err
	}
	streamer := new(TAUStreamer)
	err = json.Unmarshal(responseBody, streamer)
	if err != nil {
		return nil, err
	}

	return streamer, nil
}
