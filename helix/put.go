package helix

import (
	"encoding/json"
	gotau "github.com/Team-TAU/tau-client-go"
	"strings"
)

// PutRequest wraps PUT requests to helix thru TAU
func (c *Client) PutRequest(endpoint string, params map[string][]string, body []byte) ([]byte, error) {
	return c.helixRequest(endpoint, params, body, "PUT")
}

// ReplaceStreamTags can be used to set/reset a streams tags, see https://dev.twitch.tv/docs/api/reference#replace-stream-tags
func (c *Client) ReplaceStreamTags(broadcasterID string, tags []string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return false, gotau.BadRequestError{
			Err: "invalid request, broadcast can't be blank",
		}
	}
	if len(tags) > 5 {
		return false, gotau.BadRequestError{
			Err: "invalid request, maximum of 5 tags can be set",
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

// UpdateUser allows you to update the description of a user.  If the description is nil it just gets the user.
// See https://dev.twitch.tv/docs/api/reference#update-user.
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
