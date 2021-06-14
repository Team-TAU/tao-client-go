package helix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// GetRequest runs a raw GET request against the twitch helix pass thru, and returns the []byte of the data.
func (c *Client) GetRequest(endpoint string, queryParams map[string][]string) ([]byte, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	for key, values := range queryParams {
		for _, item := range values {
			q.Add(key, item)
		}
	}
	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 401 {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, AuthorizationError{
			fmt.Sprintf("authorization error on /users, %s", string(body)),
		}
	}
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("error making request to twitch, %s", string(body))
	}
	return ioutil.ReadAll(response.Body)
}

// GetTwitchUsers makes an api call to https://dev.twitch.tv/docs/api/reference#get-users, and formats the data.
// This does some limited validation based on the API definition.
func (c *Client) GetTwitchUsers(logins []string, ids []string) (*Users, error) {
	if len(logins)+len(ids) > 100 {
		return nil, BadRequestError{
			err: "invalid request, get users only supports a maximum of 100 total logins and ids combined.",
		}
	}

	params := make(map[string][]string)
	params["login"] = logins
	params["id"] = ids
	body, err := c.GetRequest("users", params)
	if err != nil {
		return nil, err
	}
	result := new(Users)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetBitsLeaderboard makes an api call to https://dev.twitch.tv/docs/api/reference#get-bits-leaderboard and formats the data.
func (c *Client) GetBitsLeaderboard(count int, period string, startedAt *time.Time, userID string) (*BitsLeaderboard, error) {
	params := make(map[string][]string)
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum value for count is 100 and you input %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can not be negative",
		}
	}
	if count != 0 {
		params["count"] = []string{fmt.Sprintf("%d", count)}
	}
	if period != "" {
		switch period {
		case "all":
			fallthrough
		case "day":
			fallthrough
		case "week":
			fallthrough
		case "month":
			fallthrough
		case "year":
			params["period"] = []string{period}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, period can only all, day, week, month, or year, and you input %s", period),
			}
		}
	}

	if startedAt != nil {
		ts := startedAt.Format(time.RFC3339)
		params["started_at"] = []string{ts}
	}

	if userID != "" {
		params["user_id"] = []string{userID}
	}

	body, err := c.GetRequest("bits/leaderboard", params)
	if err != nil {
		return nil, err
	}

	result := new(BitsLeaderboard)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetCheermotes makes an api call to https://dev.twitch.tv/docs/api/reference#get-cheermotes and formats the data.
func (c *Client) GetCheermotes(broadcasterID string) (*CheermotesList, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	params := make(map[string][]string)
	if broadcasterID != "" {
		params["broadcaster_id"] = []string{broadcasterID}
	}

	body, err := c.GetRequest("bits/cheermotes", params)
	if err != nil {
		return nil, err
	}
	result := new(CheermotesList)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetChannelInformation makes an api call to https://dev.twitch.tv/docs/api/reference#get-channel-information and formats the data.
func (c *Client) GetChannelInformation(broadcasterID string) (*ChannelInformation, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	param := make(map[string][]string)
	param["broadcaster_id"] = []string{broadcasterID}
	body, err := c.GetRequest("channels", param)
	if err != nil {
		return nil, err
	}

	result := new(ChannelInformation)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetChannelEditors makes an api call to https://dev.twitch.tv/docs/api/reference#get-channel-editors and formats the data.
func (c *Client) GetChannelEditors(broadcasterID string) (*ChannelEditors, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}

	param := make(map[string][]string)
	param["broadcaster_id"] = []string{broadcasterID}
	body, err := c.GetRequest("channels/editors", param)
	if err != nil {
		return nil, err
	}

	result := new(ChannelEditors)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetCustomRewards makes an api call to https://dev.twitch.tv/docs/api/reference#get-custom-reward and formats the data.
func (c *Client) GetCustomRewards(broadcasterID string, rewardID []string, onlyManageableRewards bool) (*CustomRewards, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if len(rewardID) > 50 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum number of rewardIDs is 50, but you input %d", len(rewardID)),
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if len(rewardID) > 0 {
		params["id"] = rewardID
	}

	if onlyManageableRewards {
		params["only_manageable_rewards"] = []string{"true"}
	}

	body, err := c.GetRequest("channel_points/custom_rewards", params)
	if err != nil {
		return nil, err
	}

	result := new(CustomRewards)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetCustomRewardRedemption makes an api call to https://dev.twitch.tv/docs/api/reference#get-custom-reward-redemption and formats the data.
func (c *Client) GetCustomRewardRedemption(broadcasterID string, rewardID string, redemptionID []string,
	status string, sort string, cursor string, resultCount int) (*CustomRewardRedemptions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	rewardID = strings.TrimSpace(rewardID)
	if rewardID == "" {
		return nil, BadRequestError{
			"invalid request, rewardID can't be blank",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	params["reward_id"] = []string{rewardID}

	if len(redemptionID) > 50 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum number of redemptionIDs is 50, but you input %d", len(redemptionID)),
		}
	} else if len(redemptionID) > 0 {
		params["id"] = redemptionID
	}

	status = strings.TrimSpace(status)
	if len(redemptionID) == 0 && status == "" {
		return nil, BadRequestError{
			"invalid request, if there are no redemptionIDs, status must be set",
		}
	}

	if status != "" {
		switch status {
		case "UNFULFILLED":
			fallthrough
		case "FULFILLED":
			fallthrough
		case "CANCELED":
			params["status"] = []string{status}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, status can only be UNFULFILLED, FILFILLED, or CANCELED, but you input %s", status),
			}
		}
	}

	sort = strings.TrimSpace(sort)
	if sort != "" {
		switch sort {
		case "OLDEST":
			fallthrough
		case "NEWEST":
			params["sort"] = []string{sort}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, sort can only be OLDEST or NEWEST, but you input %s", sort),
			}
		}
	}

	if resultCount > 50 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum result count is 50, you input %d", resultCount),
		}
	} else if resultCount < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if resultCount != 0 {
		params["first"] = []string{fmt.Sprintf("%d", resultCount)}
	}

	cursor = strings.TrimSpace(cursor)
	if cursor != "" {
		params["after"] = []string{cursor}
	}

	body, err := c.GetRequest("channel_points/custom_rewards/redemptions", params)
	if err != nil {
		return nil, err
	}

	result := new(CustomRewardRedemptions)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetChannelChatBadges makes an api call to https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges and formats the data.
func (c *Client) GetChannelChatBadges(broadcasterID string) (*ChannelChatBadges, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	param := make(map[string][]string)
	param["broadcaster_id"] = []string{broadcasterID}

	body, err := c.GetRequest("chat/badges", param)
	if err != nil {
		return nil, err
	}

	result := new(ChannelChatBadges)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetGlobalChatBadges makes an api call to https://dev.twitch.tv/docs/api/reference#get-global-chat-badges and formats the data.
func (c *Client) GetGlobalChatBadges() (*ChannelChatBadges, error) {
	body, err := c.GetRequest("chat/badges/global", nil)
	if err != nil {
		return nil, err
	}

	result := new(ChannelChatBadges)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetClipsByBroadcaster makes an api call to https://dev.twitch.tv/docs/api/reference#get-clips based on broadcaster, and formats the data.
func (c *Client) GetClipsByBroadcaster(broadcasterID, after, before string, startedAt,
	endedAt *time.Time, count int) (*Clips, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count has a maximum of 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count cannot be negative",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if after != "" {
		params["after"] = []string{after}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if startedAt != nil {
		params["started_at"] = []string{startedAt.Format(time.RFC3339)}
	}
	if endedAt != nil {
		params["ended_at"] = []string{endedAt.Format(time.RFC3339)}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}

	body, err := c.GetRequest("clips", params)
	if err != nil {
		return nil, err
	}
	clips := new(Clips)

	err = json.Unmarshal(body, clips)
	if err != nil {
		return nil, err
	}

	return clips, nil
}

// GetClipsByGame makes an api call to https://dev.twitch.tv/docs/api/reference#get-clips based on game, and formats the data.
func (c *Client) GetClipsByGame(gameID, after, before string, startedAt,
	endedAt *time.Time, count int) (*Clips, error) {
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return nil, BadRequestError{
			"invalid request, gameID can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count has a maximum of 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count cannot be negative",
		}
	}
	params := make(map[string][]string)
	params["game_id"] = []string{gameID}
	if after != "" {
		params["after"] = []string{after}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if startedAt != nil {
		params["started_at"] = []string{startedAt.Format(time.RFC3339)}
	}
	if endedAt != nil {
		params["ended_at"] = []string{endedAt.Format(time.RFC3339)}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}

	body, err := c.GetRequest("clips", params)
	if err != nil {
		return nil, err
	}
	clips := new(Clips)

	err = json.Unmarshal(body, clips)
	if err != nil {
		return nil, err
	}

	return clips, nil
}

// GetClipsByID makes an api call to https://dev.twitch.tv/docs/api/reference#get-clips based on clip id, and formats the data.
func (c *Client) GetClipsByID(clipID []string, after, before string, startedAt,
	endedAt *time.Time, count int) (*Clips, error) {
	if len(clipID) == 0 {
		return nil, BadRequestError{
			"invalid request, clipID can't be empty",
		}
	} else if len(clipID) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, clipID has a maximum of 100 but you supplied %d ids", len(clipID)),
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count has a maximum of 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count cannot be negative",
		}
	}
	params := make(map[string][]string)
	params["id"] = clipID

	if after != "" {
		params["after"] = []string{after}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if startedAt != nil {
		params["started_at"] = []string{startedAt.Format(time.RFC3339)}
	}
	if endedAt != nil {
		params["ended_at"] = []string{endedAt.Format(time.RFC3339)}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}

	body, err := c.GetRequest("clips", params)
	if err != nil {
		return nil, err
	}
	clips := new(Clips)

	err = json.Unmarshal(body, clips)
	if err != nil {
		return nil, err
	}

	return clips, nil
}

// GetEventSubSubscriptions makes an api call to https://dev.twitch.tv/docs/api/reference#get-eventsub-subscriptions, and formats the data.
func (c *Client) GetEventSubSubscriptions(status, eventType string) (*EventSubSubscriptions, error) {
	params := make(map[string][]string)
	if status != "" {
		params["status"] = []string{status}
	}
	if eventType != "" {
		params["type"] = []string{eventType}
	}

	body, err := c.GetRequest("eventsub/subscriptions", params)
	if err != nil {
		return nil, err
	}

	subs := new(EventSubSubscriptions)
	err = json.Unmarshal(body, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// GetTopGames makes an api call to https://dev.twitch.tv/docs/api/reference#get-top-games, and formats the data.
func (c *Client) GetTopGames(before, after string, count int) (*Games, error) {
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	params := make(map[string][]string)
	if before != "" {
		params["before"] = []string{before}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}

	body, err := c.GetRequest("games/top", params)
	if err != nil {
		return nil, err
	}

	games := new(Games)
	err = json.Unmarshal(body, games)
	if err != nil {
		return nil, err
	}
	return games, nil
}

// GetGames makes an api call to https://dev.twitch.tv/docs/api/reference#get-games, and formats the data.
func (c *Client) GetGames(ids, names []string) (*Games, error) {
	if len(ids) == 0 && len(names) == 0 {
		return nil, BadRequestError{
			"invalid request, either id or names is necessary",
		}
	}
	if len(ids)+len(names) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum number of games and ids is 100, you input %d", len(ids)+len(names)),
		}
	}

	params := make(map[string][]string)
	params["id"] = ids
	params["name"] = names

	body, err := c.GetRequest("games", params)
	if err != nil {
		return nil, err
	}

	games := new(Games)
	err = json.Unmarshal(body, games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

// GetHypeTrainEvents makes an api call to https://dev.twitch.tv/docs/api/reference#get-hype-train-events, and formats the data.
func (c *Client) GetHypeTrainEvents(broadcasterID string, count int, id, cursor string) (*HypeTrainEvents, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if id != "" {
		params["id"] = []string{id}
	}
	if cursor != "" {
		params["cursor"] = []string{cursor}
	}

	body, err := c.GetRequest("hypetrain/events", params)
	if err != nil {
		return nil, err
	}
	events := new(HypeTrainEvents)
	err = json.Unmarshal(body, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetBannedEvents makes an api call to https://dev.twitch.tv/docs/api/reference#get-banned-events, and formats the data.
func (c *Client) GetBannedEvents(broadcasterID, userID, after string, count int) (*BannedEvents, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if userID != "" {
		params["user_id"] = []string{userID}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("moderation/banned/events", params)
	if err != nil {
		return nil, err
	}

	events := new(BannedEvents)
	err = json.Unmarshal(body, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetModerators makes an api call to https://dev.twitch.tv/docs/api/reference#get-moderators, and formats the data.
func (c *Client) GetModerators(broadcasterID string, userIDs []string, after string, count int) (*Moderators, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(userIDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum user ids that can be supplied is 100 but you supplied %d", len(userIDs)),
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	params["user_id"] = userIDs

	body, err := c.GetRequest("moderation/moderators", params)
	if err != nil {
		return nil, err
	}

	mods := new(Moderators)
	err = json.Unmarshal(body, mods)
	if err != nil {
		return nil, err
	}

	return mods, nil
}

// GetModeratorEvents makes an api call to https://dev.twitch.tv/docs/api/reference#get-moderator-events, and formats the data.
func (c *Client) GetModeratorEvents(broadcasterID string, userIDs []string, after string, count int) (*ModeratorEvents, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(userIDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum user ids that can be supplied is 100 but you supplied %d", len(userIDs)),
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	params["user_id"] = userIDs

	body, err := c.GetRequest("moderation/moderators/events", params)
	if err != nil {
		return nil, err
	}

	events := new(ModeratorEvents)
	err = json.Unmarshal(body, events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetPolls makes an api call to https://dev.twitch.tv/docs/api/reference#get-polls, and formats the data.
func (c *Client) GetPolls(broadcasterID string, IDs []string, after string, count int) (*Polls, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 20 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 20, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(IDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum ids that can be supplied is 100 but you supplied %d", len(IDs)),
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	params["id"] = IDs

	body, err := c.GetRequest("polls", params)
	if err != nil {
		return nil, err
	}

	polls := new(Polls)
	err = json.Unmarshal(body, polls)
	if err != nil {
		return nil, err
	}
	return polls, nil
}

// GetPredictions makes an api call to https://dev.twitch.tv/docs/api/reference#get-predictions, and formats the data.
func (c *Client) GetPredictions(broadcasterID string, IDs []string, after string, count int) (*Predictions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 20 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 20, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(IDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum ids that can be supplied is 100 but you supplied %d", len(IDs)),
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	params["id"] = IDs

	body, err := c.GetRequest("predictions", params)
	if err != nil {
		return nil, err
	}

	predictions := new(Predictions)
	err = json.Unmarshal(body, predictions)
	if err != nil {
		return nil, err
	}
	return predictions, nil
}

// SearchCategories makes an api call to https://dev.twitch.tv/docs/api/reference#search-categories, and formats the data.
func (c *Client) SearchCategories(query, after string, count int) (*Games, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, BadRequestError{
			"invalid request, query can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["query"] = []string{query}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("search/categories", params)
	if err != nil {
		return nil, err
	}

	results := new(Games)
	err = json.Unmarshal(body, results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// SearchChannels makes an api call to https://dev.twitch.tv/docs/api/reference#search-channels, and formats the data.
func (c *Client) SearchChannels(query, after string, count int, liveOnly bool) (*ChannelSearchResults, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, BadRequestError{
			"invalid request, query can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["query"] = []string{query}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if liveOnly {
		params["live_only"] = []string{"true"}
	}

	body, err := c.GetRequest("search/channels", params)
	if err != nil {
		return nil, err
	}

	results := new(ChannelSearchResults)
	err = json.Unmarshal(body, results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetStreamKey makes an api call to https://dev.twitch.tv/docs/api/reference#get-stream-key, and formats the data.
func (c *Client) GetStreamKey(broadcasterID string) (*StreamKey, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}

	body, err := c.GetRequest("streams/key", params)
	if err != nil {
		return nil, err
	}

	key := new(StreamKey)
	err = json.Unmarshal(body, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetStreams makes an api call to https://dev.twitch.tv/docs/api/reference#get-streams, and formats the data.
func (c *Client) GetStreams(before, after string, count int, gameIDs []string, languages []string,
	userIDs []string, userLogins []string) (*Streams, error) {
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(gameIDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, max number of game ids is 100, but you supplied %d", len(gameIDs)),
		}
	}
	if len(languages) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, max number of languages is 100, but you supplied %d", len(languages)),
		}
	}
	if len(userIDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, max number of user ids is 100, but you supplied %d", len(userIDs)),
		}
	}
	if len(userLogins) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, max number of user logins is 100, but you supplied %d", len(userLogins)),
		}
	}

	params := make(map[string][]string)
	if before != "" {
		params["before"] = []string{before}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if len(gameIDs) > 0 {
		params["game_id"] = gameIDs
	}
	if len(languages) > 0 {
		params["language"] = languages
	}
	if len(userIDs) > 0 {
		params["user_id"] = userIDs
	}
	if len(userLogins) > 0 {
		params["user_login"] = userLogins
	}

	body, err := c.GetRequest("streams", params)
	if err != nil {
		return nil, err
	}

	streams := new(Streams)
	err = json.Unmarshal(body, streams)
	if err != nil {
		return nil, err
	}
	return streams, nil
}

// GetFollowedStreams makes an api call to https://dev.twitch.tv/docs/api/reference#get-followed-streams, and formats the data.
func (c *Client) GetFollowedStreams(userID, after string, count int) (*Streams, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	params := make(map[string][]string)
	params["user_id"] = []string{userID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("streams/followed", params)
	if err != nil {
		return nil, err
	}
	streams := new(Streams)
	err = json.Unmarshal(body, streams)
	if err != nil {
		return nil, err
	}

	return streams, nil
}

// GetStreamMarkers makes an api call to https://dev.twitch.tv/docs/api/reference#get-stream-markers, and formats the data.
func (c *Client) GetStreamMarkers(userID, videoID, before, after string, count int) (*StreamMarkers, error) {
	userID = strings.TrimSpace(userID)
	videoID = strings.TrimSpace(videoID)
	if userID == "" && videoID == "" {
		return nil, BadRequestError{
			"invalid request, user id or video id must be set",
		}
	}
	if userID != "" && videoID != "" {
		return nil, BadRequestError{
			"invalid request, only one of userID or videoID can be set",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	if userID != "" {
		params["user_id"] = []string{userID}
	} else {
		params["video_id"] = []string{videoID}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("streams/markers", params)
	if err != nil {
		return nil, err
	}
	markers := new(StreamMarkers)
	err = json.Unmarshal(body, markers)
	if err != nil {
		return nil, err
	}

	return markers, nil
}

// GetBroadcasterSubscriptions makes an api call to https://dev.twitch.tv/docs/api/reference#get-stream-markers, and formats the data.
func (c *Client) GetBroadcasterSubscriptions(broadcasterID string, userIDs []string, after string, count int) (*Subscriptions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if len(userIDs) > 0 {
		params["user_id"] = userIDs
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if count > 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}

	body, err := c.GetRequest("subscriptions", params)
	if err != nil {
		return nil, err
	}

	subs := new(Subscriptions)
	err = json.Unmarshal(body, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// CheckUserSubscription makes an api call to https://dev.twitch.tv/docs/api/reference#check-user-subscription, and formats the data.
func (c *Client) CheckUserSubscription(broadcasterID, userID string) (*UserSubscriptions, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	userID = strings.TrimSpace(userID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if userID == "" {
		return nil, BadRequestError{
			"invalid request, user can't be blank",
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	params["user_id"] = []string{userID}

	body, err := c.GetRequest("subscriptions/user", params)
	if err != nil {
		return nil, err
	}

	subs := new(UserSubscriptions)
	err = json.Unmarshal(body, subs)

	if err != nil {
		return nil, err
	}
	return subs, nil
}

// GetAllStreamTags makes an api call to https://dev.twitch.tv/docs/api/reference#get-all-stream-tags, and formats the data.
func (c *Client) GetAllStreamTags(after string, count int, tagIDs []string) (*StreamTags, error) {
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}
	if len(tagIDs) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, max number of tag ids is 100, but you supplied %d", len(tagIDs)),
		}
	}

	params := make(map[string][]string)
	if count > 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if len(tagIDs) > 0 {
		params["tag_id"] = tagIDs
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("tags/streams", params)
	if err != nil {
		return nil, err
	}

	tags := new(StreamTags)
	err = json.Unmarshal(body, tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetStreamTags makes an api call to https://dev.twitch.tv/docs/api/reference#get-stream-tags, and formats the data.
func (c *Client) GetStreamTags(broadcasterID string) (*StreamTags, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}

	body, err := c.GetRequest("streams/tags", params)
	if err != nil {
		return nil, err
	}

	tags := new(StreamTags)
	err = json.Unmarshal(body, tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// GetChannelTeams makes an api call to https://dev.twitch.tv/docs/api/reference#get-channel-teams, and formats the data.
func (c *Client) GetChannelTeams(broadcasterID string) (*ChannelTeams, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}

	body, err := c.GetRequest("teams/channel", params)
	if err != nil {
		return nil, err
	}

	teams := new(ChannelTeams)
	err = json.Unmarshal(body, teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// GetTeam makes an api call to https://dev.twitch.tv/docs/api/reference#get-teams, and formats the data.
func (c *Client) GetTeam(name, id string) (*Teams, error) {
	name = strings.TrimSpace(name)
	id = strings.TrimSpace(id)

	if name == "" && id == "" {
		return nil, BadRequestError{
			"invalid request, name or id must be specified",
		}
	}
	if name != "" && id != "" {
		return nil, BadRequestError{
			"invalid request, only one of name or id may be specified",
		}
	}
	params := make(map[string][]string)
	if name != "" {
		params["name"] = []string{name}
	} else {
		params["id"] = []string{id}
	}

	body, err := c.GetRequest("teams", params)
	if err != nil {
		return nil, err
	}

	teams := new(Teams)
	err = json.Unmarshal(body, teams)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// GetUsers makes an api call to https://dev.twitch.tv/docs/api/reference#get-users, and formats the data.
func (c *Client) GetUsers(IDs, logins []string) (*Users, error) {
	if len(IDs) == 0 && len(logins) == 0 {
		return nil, BadRequestError{
			"invalid request, login or id must be specified",
		}
	}
	if len(IDs)+len(logins) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum of 100 ids and logins, but you supplied %d", len(IDs)+len(logins)),
		}
	}

	params := make(map[string][]string)
	if len(IDs) > 0 {
		params["id"] = IDs
	}
	if len(logins) > 0 {
		params["login"] = logins
	}

	body, err := c.GetRequest("users", params)
	if err != nil {
		return nil, err
	}

	users := new(Users)
	err = json.Unmarshal(body, users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersFollows makes an api call to https://dev.twitch.tv/docs/api/reference#get-users-follows, and formats the data.
func (c *Client) GetUsersFollows(fromID, toID, after string, count int) (*UserFollows, error) {
	fromID = strings.TrimSpace(fromID)
	toID = strings.TrimSpace(toID)
	if fromID == "" && toID == "" {
		return nil, BadRequestError{"invalid request, either fromID or toID must be specified"}
	}
	if fromID != "" && toID != "" {
		return nil, BadRequestError{"invalid request, fromID or toID must be specified, but not both"}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	if toID != "" {
		params["to_id"] = []string{toID}
	} else {
		params["from_id"] = []string{fromID}
	}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("users/follows", params)
	if err != nil {
		return nil, err
	}

	follows := new(UserFollows)
	err = json.Unmarshal(body, follows)
	if err != nil {
		return nil, err
	}

	return follows, nil
}

// GetUsersBlockList makes an api call to https://dev.twitch.tv/docs/api/reference#get-user-block-list, and formats the data.
func (c *Client) GetUsersBlockList(broadcasterID, after string, count int) (*UserBlockList, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	if broadcasterID == "" {
		return nil, BadRequestError{
			"invalid request, broadcast can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["broadcaster_id"] = []string{broadcasterID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("users/blocks", params)
	if err != nil {
		return nil, err
	}

	blockList := new(UserBlockList)
	err = json.Unmarshal(body, blockList)
	if err != nil {
		return nil, err
	}

	return blockList, nil
}

// GetUserExtensions makes an api call to https://dev.twitch.tv/docs/api/reference#get-user-extensions, and formats the data.
func (c *Client) GetUserExtensions() (*UserExtensions, error) {
	body, err := c.GetRequest("users/extensions/list", nil)
	if err != nil {
		return nil, err
	}

	extensions := new(UserExtensions)
	err = json.Unmarshal(body, extensions)
	if err != nil {
		return nil, err
	}

	return extensions, nil
}

// GetUserActiveExtensions makes an api call to https://dev.twitch.tv/docs/api/reference#get-user-active-extensions, and formats the data.
func (c *Client) GetUserActiveExtensions(userID string) (*UserActiveExtensions, error) {
	userID = strings.TrimSpace(userID)
	params := make(map[string][]string)
	if userID != "" {
		params["user_id"] = []string{userID}
	}

	body, err := c.GetRequest("users/extensions", params)
	if err != nil {
		return nil, err
	}

	extensions := new(UserActiveExtensions)
	err = json.Unmarshal(body, extensions)
	if err != nil {
		return nil, err
	}

	return extensions, nil
}

// GetVideosByID makes an api call to https://dev.twitch.tv/docs/api/reference#get-videos, and formats the data.
func (c *Client) GetVideosByID(ids []string) (*Video, error) {
	if len(ids) == 0 {
		return nil, BadRequestError{
			"invalid request, at least one video id is required",
		}
	}
	if len(ids) > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, maximum of 100 ids, but you supplied %d", len(ids)),
		}
	}

	params := make(map[string][]string)
	params["id"] = ids

	body, err := c.GetRequest("videos", params)
	if err != nil {
		return nil, err
	}

	video := new(Video)
	err = json.Unmarshal(body, video)
	if err != nil {
		return nil, err
	}

	return video, nil
}

// GetVideosByUser makes an api call to https://dev.twitch.tv/docs/api/reference#get-videos, and formats the data.
func (c *Client) GetVideosByUser(userID, before, after, language string, count int,
	period, sort, _type string) (*Video, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, BadRequestError{
			"invalid request, user id can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["user_id"] = []string{userID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if language != "" {
		params["language"] = []string{language}
	}

	if period != "" {
		switch period {
		case "all":
			fallthrough
		case "day":
			fallthrough
		case "week":
			fallthrough
		case "month":
			params["period"] = []string{period}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, period can only all, day, week, month, or year, and you input %s", period),
			}
		}
	}

	if sort != "" {
		switch sort {
		case "time":
			fallthrough
		case "trending":
			fallthrough
		case "views":
			params["sort"] = []string{sort}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, sort can only be all, time, trending, or views, but you input %s", sort),
			}
		}
	}

	if _type != "" {
		switch _type {
		case "all":
			fallthrough
		case "upload":
			fallthrough
		case "archive":
			fallthrough
		case "highlight":
			params["type"] = []string{_type}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, type can only be all, upload, archive, or highlight, but you input %s", _type),
			}
		}
	}

	body, err := c.GetRequest("videos", params)
	if err != nil {
		return nil, err
	}

	videos := new(Video)
	err = json.Unmarshal(body, videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

// GetVideosByGame makes an api call to https://dev.twitch.tv/docs/api/reference#get-videos, and formats the data.
func (c *Client) GetVideosByGame(gameID, before, after, language string, count int,
	period, sort, _type string) (*Video, error) {
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return nil, BadRequestError{
			"invalid request, user id can't be blank",
		}
	}
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	params["game_id"] = []string{gameID}
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if before != "" {
		params["before"] = []string{before}
	}
	if after != "" {
		params["after"] = []string{after}
	}
	if language != "" {
		params["language"] = []string{language}
	}

	if period != "" {
		switch period {
		case "all":
			fallthrough
		case "day":
			fallthrough
		case "week":
			fallthrough
		case "month":
			params["period"] = []string{period}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, period can only all, day, week, month, or year, and you input %s", period),
			}
		}
	}

	if sort != "" {
		switch sort {
		case "time":
			fallthrough
		case "trending":
			fallthrough
		case "views":
			params["sort"] = []string{sort}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, sort can only be all, time, trending, or views, but you input %s", sort),
			}
		}
	}

	if _type != "" {
		switch _type {
		case "all":
			fallthrough
		case "upload":
			fallthrough
		case "archive":
			fallthrough
		case "highlight":
			params["type"] = []string{_type}
		default:
			return nil, BadRequestError{
				fmt.Sprintf("invalid request, type can only be all, upload, archive, or highlight, but you input %s", _type),
			}
		}
	}

	body, err := c.GetRequest("videos", params)
	if err != nil {
		return nil, err
	}

	videos := new(Video)
	err = json.Unmarshal(body, videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (c *Client) GetWebhookSubscriptions(after string, count int) (*WebhookSubscriptions, error) {
	if count > 100 {
		return nil, BadRequestError{
			fmt.Sprintf("invalid request, count maximum value is 100, but you supplied %d", count),
		}
	} else if count < 0 {
		return nil, BadRequestError{
			"invalid request, count can't be negative",
		}
	}

	params := make(map[string][]string)
	if count != 0 {
		params["first"] = []string{fmt.Sprintf("%d", count)}
	}
	if after != "" {
		params["after"] = []string{after}
	}

	body, err := c.GetRequest("webhooks/subscriptions", params)
	if err != nil {
		return nil, err
	}

	subscription := new(WebhookSubscriptions)
	err = json.Unmarshal(body, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
