package gotau

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetStreamers() ([]*TAUStreamers, error) {
	body, err := c.apiRequest("streamers", nil, nil, "GET")
	if err != nil {
		return nil, err
	}

	var streamers []*TAUStreamers
	err = json.Unmarshal(body, &streamers)
	if err != nil {
		return nil, err
	}

	return streamers, err
}

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
