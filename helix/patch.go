package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// PatchRequest handles generic PATCH requests to twitch's API, leveraged internally as well as allows you
//to make raw requests in case an update to the API comes out and the library hasn't been updated yet.
func (c *Client) PatchRequest(endpoint string, params map[string][]string, body []byte) (bool, []byte, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest("PATCH", endpointURL, buffer)
	if err != nil {
		return false, nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	request.Header.Add("Content-Type", "application/json")
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return false, nil, err
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
		return false, nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if response.StatusCode == 204 {
			return true, nil, nil
		} else {
			body, err := ioutil.ReadAll(response.Body)
			return true, body, err
		}
	}
	if response.StatusCode == 401 {
		return false, nil, AuthorizationError{}
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		err = GenericError{
			err:  fmt.Sprintf("response code %d: %s", response.StatusCode, body),
			body: body,
			code: response.StatusCode,
		}
		return false, nil, err
	}
}

// ModifyChannelInformation updates the channel information based on https://dev.twitch.tv/docs/api/reference#modify-channel-information.
func (c *Client) ModifyChannelInformation(broadcasterID string, gameID, language, title *string, delay *int) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
	}

	bodyObject := make(map[string]string)
	if gameID != nil {
		bodyObject["game_id"] = *gameID
	}
	if language != nil {
		bodyObject["broadcaster_language"] = *language
	}
	if title != nil {
		bodyObject["title"] = *title
	}
	if delay != nil {
		bodyObject["delay"] = fmt.Sprintf("%d", *delay)
	}

	if len(bodyObject) == 0 {
		return false, BadRequestError{
			"invalid request, at least one parameter must be provided of gameID, language, title, and delay",
		}
	}
	body, err := json.Marshal(bodyObject)
	if err != nil {
		return false, err
	}

	changed, _, err := c.PatchRequest("channels", params, body)
	return changed, err
}

// UpdateCustomReward updates a custom reward that is owned by your client id.
//For more information see https://dev.twitch.tv/docs/api/reference#update-custom-reward.
func (c *Client) UpdateCustomReward(broadcasterID string, ID string, change *CustomRewardsUpdate) (*CustomRewards, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if ID == "" {
		return nil, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}
	if change == nil {
		return nil, BadRequestError{
			"invalid request, change can't be nil",
		}
	}
	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {ID},
	}

	body, err := json.Marshal(change)
	if err != nil {
		return nil, err
	}

	_, response, err := c.PatchRequest("channel_points/custom_rewards", params, body)
	if err != nil {
		return nil, err
	}
	reward := new(CustomRewards)
	err = json.Unmarshal(response, reward)
	if err != nil {
		return nil, err
	}

	return reward, nil
}

// UpdateRedemptionStatus updates the status of one or more redemptions for a reward owned by your client id.
// For more information see https://dev.twitch.tv/docs/api/reference#update-redemption-status.
func (c *Client) UpdateRedemptionStatus(broadcasterID, rewardID string, redemptionIDs []string, status string) (*CustomRewardRedemptions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	rewardID = strings.TrimSpace(rewardID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if rewardID == "" {
		return nil, BadRequestError{
			"invalid request, rewardID can't be blank",
		}
	}
	if len(redemptionIDs) == 0 {
		return nil, BadRequestError{
			"invalid request, redemptionIDs can't be empty or nil",
		}
	} else if len(redemptionIDs) > 50 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request,  maximum of 50 redemptionIDs, but you supplied %d", len(redemptionIDs)),
		}
	}
	if status != "FULFILLED" && status != "CANCELED" {
		return nil, BadRequestError{
			"invalid request,  status can only be one of FULFILLED or CANCELED",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"reward_id":      {rewardID},
		"id":             redemptionIDs,
	}

	bodyMap := map[string]string{
		"status": status,
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	_, response, err := c.PatchRequest("channel_points/custom_rewards", params, body)
	if err != nil {
		return nil, err
	}

	redemptions := new(CustomRewardRedemptions)
	err = json.Unmarshal(response, redemptions)
	if err != nil {
		return nil, err
	}

	return redemptions, nil
}

// EndPoll allows you to end/archive a poll, see https://dev.twitch.tv/docs/api/reference#end-poll
func (c *Client) EndPoll(broadcasterID, pollID, status string) (*Polls, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	pollID = strings.TrimSpace(pollID)
	status = strings.TrimSpace(status)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if pollID == "" {
		return nil, BadRequestError{
			"invalid request, poll id can't be blank",
		}
	}
	if status != "TERMINATED" && status != "ARCHIVED" {
		return nil, BadRequestError{
			"invalid request, status must either be TERMINATED or ARCHIVED",
		}
	}
	bodyMap := map[string]string{
		"broadcaster_id": broadcasterID,
		"id":             pollID,
		"status":         status,
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	_, responseBody, err := c.PatchRequest("polls", nil, body)
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

// EndPrediction allows you to lock/payout/cancel a prediction, see https://dev.twitch.tv/docs/api/reference#end-prediction
func (c *Client) EndPrediction(broadcasterID, predictionID, status string, winningOutcome *string) (*Predictions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	predictionID = strings.TrimSpace(predictionID)
	status = strings.TrimSpace(status)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if predictionID == "" {
		return nil, BadRequestError{
			"invalid request, prediction id can't be blank",
		}
	}
	if status != "RESOLVED" && status != "CANCELED" && status != "LOCKED" {
		return nil, BadRequestError{
			"invalid request, status must either be RESOLVED, CANCELED, or LOCKED",
		}
	}
	if status == "RESOLVED" && winningOutcome == nil {
		return nil, BadRequestError{
			"invalid request, if status RESOLVED, winning outcome must be set",
		}
	}
	bodyMap := map[string]string{
		"broadcaster_id": broadcasterID,
		"id":             predictionID,
		"status":         status,
	}
	if winningOutcome != nil {
		bodyMap["winning_outcome_id"] = *winningOutcome
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	_, responseBody, err := c.PatchRequest("predictions", nil, body)
	if err != nil {
		return nil, err
	}
	prediction := new(Predictions)
	err = json.Unmarshal(responseBody, prediction)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}

// UpdateChannelStreamSchedule updates the channel stream schedule per https://dev.twitch.tv/docs/api/reference#update-channel-stream-schedule.
func (c *Client) UpdateChannelStreamSchedule(broadcasterID string, vacationEnabled *bool,
	vacationStartTime, vacationEndTime *time.Time, timezone *string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if vacationEnabled != nil && *vacationEnabled {
		if vacationStartTime == nil {
			return false, BadRequestError{
				"invalid request, if vacationEnabled, vacationStartTime must be specified",
			}
		}
		if vacationEndTime == nil {
			return false, BadRequestError{
				"invalid request, if vacationEnabled, vacationEndTime must be specified",
			}
		}
		if timezone == nil {
			return false, BadRequestError{
				"invalid request, if vacationEnabled, timezone must be specified",
			}
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if vacationEnabled != nil && *vacationEnabled {
		params["is_vacation_enabled"] = []string{"true"}
		params["vacation_start_time"] = []string{vacationStartTime.Format(time.RFC3339)}
		params["vacation_end_time"] = []string{vacationEndTime.Format(time.RFC3339)}
		params["timezone"] = []string{*timezone}
	} else if vacationEnabled != nil && !*vacationEnabled {
		params["is_vacation_enabled"] = []string{"false"}
	}

	changed, _, err := c.PatchRequest("schedule/settings", params, nil)
	if err != nil {
		return false, err
	}

	return changed, nil
}

// UpdateChannelStreamScheduleSegment updates a segment of your schedule per
// https://dev.twitch.tv/docs/api/reference#update-channel-stream-schedule-segment returning the resulting schedule segment.
func (c *Client) UpdateChannelStreamScheduleSegment(broadcasterID, segmentID string, update *StreamScheduleSegmentUpdate) (*ChannelStreamSchedule, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	segmentID = strings.TrimSpace(segmentID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if segmentID == "" {
		return nil, BadRequestError{
			"invalid request, segment id can't be blank",
		}
	}
	if update == nil {
		return nil, BadRequestError{
			"invalid request, update can't be nil",
		}
	}
	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {segmentID},
	}

	body, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}
	_, responseBody, err := c.PatchRequest("schedule/segment", params, body)
	if err != nil {
		return nil, err
	}

	schedule := new(ChannelStreamSchedule)
	err = json.Unmarshal(responseBody, schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
