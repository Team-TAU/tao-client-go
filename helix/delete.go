package helix

import (
	"fmt"
	gotau "github.com/Team-TAU/tau-client-go"
	"strings"
)

// DeleteRequest runs a raw GET request against the twitch helix pass thru, and returns true if the delete worked.
func (c *Client) DeleteRequest(endpoint string, params map[string][]string) (bool, error) {
	_, err := c.helixRequest(endpoint, params, nil, "DELETE")
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteCustomReward makes an api call to https://dev.twitch.tv/docs/api/reference#delete-custom-reward, and formats the data.
func (c *Client) DeleteCustomReward(broadcasterID, ID string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)
	if broadcasterID == "" {
		return false, gotau.BadRequestError{
			"invalid request, broadcaster can't be blank",
		}
	}
	if ID == "" {
		return false, gotau.BadRequestError{
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
		return false, gotau.BadRequestError{
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
		return false, gotau.BadRequestError{
			"invalid request, fromID can't be blank",
		}
	}
	if toID == "" {
		return false, gotau.BadRequestError{
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
		return false, gotau.BadRequestError{
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
		return false, gotau.BadRequestError{
			"invalid request, IDs can't be empty",
		}
	}
	if len(IDs) > 5 {
		return false, gotau.BadRequestError{
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
		return false, gotau.BadRequestError{
			"invalid request, broadcaster can't be blank",
		}
	}
	if ID == "" {
		return false, gotau.BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {ID},
	}

	return c.DeleteRequest("schedule/segment", params)
}
