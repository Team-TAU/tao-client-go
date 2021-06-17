package helix

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DeleteRequest runs a raw GET request against the twitch helix pass thru, and returns true if the delete worked.
func (c *Client) DeleteRequest(endpoint string, params map[string][]string) (bool, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	request, err := http.NewRequest("DELETE", endpointURL, nil)
	if err != nil {
		return false, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return false, err
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
		return false, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return true, nil
	}
	if response.StatusCode == 401 {
		return false, AuthorizationError{}
	} else if response.StatusCode == 429 {
		return false, RateLimitError{
			"rate limited: received http 429",
		}
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		err = GenericError{
			err:  fmt.Sprintf("response code %d: %s", response.StatusCode, body),
			body: body,
			code: response.StatusCode,
		}
		return false, err
	}
}

// DeleteCustomReward makes an api call to https://dev.twitch.tv/docs/api/reference#delete-custom-reward, and formats the data.
func (c *Client) DeleteCustomReward(broadcasterID, ID string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcaster can't be blank",
		}
	}
	if ID == "" {
		return false, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {ID},
	}

	return c.DeleteRequest("channel_points/custom_rewards", params)
}

// DeleteEventSubSubscription makes an api call to https://dev.twitch.tv/docs/api/reference#delete-eventsub-subscription, and formats the data.
func (c *Client) DeleteEventSubSubscription(ID string) (bool, error) {
	ID = strings.TrimSpace(ID)
	if ID == "" {
		return false, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	params := map[string][]string{
		"id": {ID},
	}

	return c.DeleteRequest("eventsub/subscriptions", params)
}

// DeleteUserFollows makes an api call to https://dev.twitch.tv/docs/api/reference#delete-user-follows, and formats the data.
func (c *Client) DeleteUserFollows(fromID, toID string) (bool, error) {
	fromID = strings.TrimSpace(fromID)
	toID = strings.TrimSpace(toID)

	if fromID == "" {
		return false, BadRequestError{
			"invalid request, fromID can't be blank",
		}
	}
	if toID == "" {
		return false, BadRequestError{
			"invalid request, toID can't be blank",
		}
	}

	params := map[string][]string{
		"from_id": {fromID},
		"to_id":   {toID},
	}

	return c.DeleteRequest("users/follows", params)
}

// UnblockUser makes an api call to https://dev.twitch.tv/docs/api/reference#unblock-user, and formats the data.
func (c *Client) UnblockUser(userID string) (bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false, BadRequestError{
			"invalid request, userID can't be blank",
		}
	}

	params := map[string][]string{
		"user_id": {userID},
	}

	return c.DeleteRequest("users/blocks", params)
}

// DeleteVideos makes an api call to https://dev.twitch.tv/docs/api/reference#delete-videos, and formats the data.
func (c *Client) DeleteVideos(IDs []string) (bool, error) {
	if len(IDs) == 0 {
		return false, BadRequestError{
			"invalid request, IDs can't be empty",
		}
	}
	if len(IDs) > 5 {
		return false, BadRequestError{
			fmt.Sprintf("invalid request, maximum number of IDs is 5 but you supplied %d", len(IDs)),
		}
	}

	params := map[string][]string{
		"id": IDs,
	}

	return c.DeleteRequest("videos", params)
}

// DeleteChannelStreamScheduleSegment makes an api call to https://dev.twitch.tv/docs/api/reference#delete-channel-stream-schedule-segment, and formats the data.
func (c *Client) DeleteChannelStreamScheduleSegment(broadcasterID, ID string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)

	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcaster can't be blank",
		}
	}
	if ID == "" {
		return false, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {ID},
	}

	return c.DeleteRequest("schedule/segment", params)
}
