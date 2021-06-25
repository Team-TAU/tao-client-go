package helix

import (
	"encoding/json"
	"fmt"
	gotau "github.com/Team-TAU/tau-client-go"
	"strings"
)

// PostRequest handles generic POST requests to twitch's API, leveraged internally as well as allows you
//to make raw requests in case an update to the API comes out and the library hasn't been updated yet.
func (c *Client) PostRequest(endpoint string, params map[string][]string, body []byte) ([]byte, error) {
	return c.helixRequest(endpoint, params, body, "POST")
}

// CreateCustomReward is used to create custom channel point rewards, see https://dev.twitch.tv/docs/api/reference#create-custom-rewards.
func (c *Client) CreateCustomReward(broadcasterID string, customReward *CustomRewardsUpdate) (*CustomRewards, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, gotau.BadRequestError{
			Err: "invalid request, broadcast can't be blank",
		}
	}
	if customReward == nil {
		return nil, gotau.BadRequestError{
			Err: "invalid request, custom reward can't be nil",
		}
	}
	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
	}

	body, err := json.Marshal(customReward)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.PostRequest("channel_points/custom_rewards", params, body)
	if err != nil {
		return nil, err
	}
	reward := new(CustomRewards)
	err = json.Unmarshal(responseBody, reward)
	if err != nil {
		return nil, err
	}

	return reward, nil
}

// CreateClip allows you to create clips, see https://dev.twitch.tv/docs/api/reference#create-clip
func (c *Client) CreateClip(broadcasterID string, hasDelay bool) (string, string, error) {
	type temp struct {
		Data []struct {
			ID      string `json:"id"`
			EditURL string `json:"edit_url"`
		} `json:"data"`
	}
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return "", "", gotau.BadRequestError{
			Err: "invalid request, broadcast can't be blank",
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if hasDelay {
		params["has_delay"] = []string{fmt.Sprintf("%t", hasDelay)}
	}

	response, err := c.PostRequest("clips", params, nil)
	if err != nil {
		return "", "", err
	}

	data := new(temp)
	err = json.Unmarshal(response, data)
	if err != nil {
		return "", "", err
	}

	if len(data.Data) > 0 {
		return data.Data[0].EditURL, data.Data[0].ID, err
	}

	return "", "", nil
}

// CreatePoll can be used to create a poll on your channel, see https://dev.twitch.tv/docs/api/reference#create-poll
func (c *Client) CreatePoll(poll *CreatePoll) (*Polls, error) {
	if poll == nil {
		return nil, gotau.BadRequestError{
			Err: "invalid request, poll can't be nil",
		}
	}

	body, err := json.Marshal(poll)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.PostRequest("polls", nil, body)
	if err != nil {
		return nil, err
	}

	polls := new(Polls)
	err = json.Unmarshal(responseBody, polls)
	if err != nil {
		return nil, err
	}

	return polls, nil
}

// CreatePrediction allows you to create predictions for your viewers to bet on with channel points.
// See https://dev.twitch.tv/docs/api/reference#create-prediction
func (c *Client) CreatePrediction(prediction *CreatePrediction) (*Predictions, error) {
	if prediction == nil {
		return nil, gotau.BadRequestError{
			Err: "invalid request, prediction can't be nil",
		}
	}

	body, err := json.Marshal(prediction)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.PostRequest("predictions", nil, body)
	if err != nil {
		return nil, err
	}

	predictions := new(Predictions)
	err = json.Unmarshal(responseBody, predictions)
	if err != nil {
		return nil, err
	}

	return predictions, nil
}

// CreateChannelStreamScheduleSegment can be used to create a new scheduled stream,
// see https://dev.twitch.tv/docs/api/reference#create-channel-stream-schedule-segment
func (c *Client) CreateChannelStreamScheduleSegment(broadcasterID string, segment *StreamScheduleSegmentUpdate) (*ChannelStreamSchedule, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, gotau.BadRequestError{
			Err: "invalid request, broadcast can't be blank",
		}
	}
	if segment == nil {
		return nil, gotau.BadRequestError{
			Err: "invalid request, segment can't be nil",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
	}

	body, err := json.Marshal(segment)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.PostRequest("schedule/segment", params, body)
	if err != nil {
		return nil, err
	}

	newSegment := new(ChannelStreamSchedule)
	err = json.Unmarshal(responseBody, newSegment)
	if err != nil {
		return nil, err
	}

	return newSegment, nil
}

// CreateUserFollows allows you to follow a user, see https://dev.twitch.tv/docs/api/reference#create-user-follows
func (c *Client) CreateUserFollows(fromID, toID string, allowNotifications bool) (bool, error) {
	fromID = strings.TrimSpace(fromID)
	toID = strings.TrimSpace(toID)

	if fromID == "" {
		return false, gotau.BadRequestError{
			Err: "invalid request, from id can't be blank",
		}
	}
	if toID == "" {
		return false, gotau.BadRequestError{
			Err: "invalid request, to id can't be blank",
		}
	}

	bodyMap := map[string]string{
		"from_id": fromID,
		"to_id":   toID,
	}
	if allowNotifications {
		bodyMap["allow_notifications"] = fmt.Sprintf("%t", allowNotifications)
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return false, err
	}

	_, err = c.PostRequest("users/follows", nil, body)
	if err != nil {
		return false, err
	}

	return true, nil
}

// StartCommercial allows you to start a commercial on a stream, see https://dev.twitch.tv/docs/api/reference#start-commercial
func (c *Client) StartCommercial(broadcasterID string, length int) (*Commercial, error) {
	type commercial struct {
		BroadcasterID string `json:"broadcaster_id"`
		Length        int    `json:"length"`
	}
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, gotau.BadRequestError{
			Err: "invalid request, from id can't be blank",
		}
	}

	commercialData := commercial{
		BroadcasterID: broadcasterID,
	}

	switch length {
	case 30:
		fallthrough
	case 60:
		fallthrough
	case 90:
		fallthrough
	case 120:
		fallthrough
	case 150:
		fallthrough
	case 180:
		commercialData.Length = length
	default:
		return nil, gotau.BadRequestError{
			Err: "invalid request, valid length values are 30, 60, 90, 120, 150, 180",
		}
	}

	body, err := json.Marshal(commercialData)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.PostRequest("channels/commercial", nil, body)
	if err != nil {
		return nil, err
	}

	response := new(Commercial)
	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
