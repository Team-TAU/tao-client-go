package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *Client) PutRequest(endpoint string, params map[string][]string, body []byte) ([]byte, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest("PUT", endpointURL, buffer)
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

// ReplaceStreamTags can be used to set/reset a streams tags, see https://dev.twitch.tv/docs/api/reference#replace-stream-tags
func (c *Client) ReplaceStreamTags(broadcasterID string, tags []string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if len(tags) > 5 {
		return false, BadRequestError{
			"invalid request, maximum of 5 tags can be set",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
	}

	bodyMap := make(map[string][]string)
	if len(tags) > 0 {
		bodyMap["tag_ids"] = tags
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return false, err
	}

	_, err = c.PutRequest("streams/tags", params, body)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) UpdateUser(description *string) (*Users, error) {
	params := make(map[string][]string)
	if description != nil {
		params["description"] = []string{*description}
	}

	body, err := c.PutRequest("users", params, nil)
	if err != nil {
		return nil, err
	}

	user := new(Users)
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
