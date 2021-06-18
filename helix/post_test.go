package helix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient_PostRequestReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	shouldBeNil, err := client.PostRequest("channels", nil, nil)
	require.ErrorIs(t, err, AuthorizationError{})
	require.Nil(t, shouldBeNil)
}

func TestClient_PostRequestReturnsGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	shouldBeNil, err := client.PostRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, GenericError{}, err)
	require.Nil(t, shouldBeNil)
}

func TestClient_PostRequestReturnsRateLimitError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		w.Header().Set("Ratelimit-Reset", "1623961625")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	shouldBeNil, err := client.PostRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, RateLimitError{}, err)
	rlErr := err.(RateLimitError)
	require.NotNil(t, rlErr)
	require.Equal(t, 2021, rlErr.ResetTime().Year())
	require.Equal(t, time.Month(6), rlErr.ResetTime().Month())
	require.Equal(t, 17, rlErr.ResetTime().Day())
	require.Nil(t, shouldBeNil)
}

func TestClient_CreateCustomRewardReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "274637212", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "POST", r.Method)

		cr := new(CustomRewardsUpdate)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &cr)
		require.NoError(t, err)
		require.NotNil(t, cr)
		require.Equal(t, "game analysis 1v1", *cr.Title)
		require.Equal(t, 50000, *cr.Cost)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"broadcaster_name\":\"torpedo09\",\"broadcaster_login\":\"torpedo09\",\"broadcaster_id\":\"274637212\",\"id\":\"afaa7e34-6b17-49f0-a19a-d1e76eaaf673\",\"image\":null,\"background_color\":\"#00E5CB\",\"is_enabled\":true,\"cost\":50000,\"title\":\"game analysis 1v1\",\"prompt\":\"\",\"is_user_input_required\":false,\"max_per_stream_setting\":{\"is_enabled\":false,\"max_per_stream\":0},\"max_per_user_per_stream_setting\":{\"is_enabled\":false,\"max_per_user_per_stream\":0},\"global_cooldown_setting\":{\"is_enabled\":false,\"global_cooldown_seconds\":0},\"is_paused\":false,\"is_in_stock\":true,\"default_image\":{\"url_1x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png\",\"url_2x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png\",\"url_4x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png\"},\"should_redemptions_skip_request_queue\":false,\"redemptions_redeemed_current_stream\":null,\"cooldown_expires_at\":null}]}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	title := "game analysis 1v1"
	cost := 50000

	customReward := CustomRewardsUpdate{
		Title: &title,
		Cost:  &cost,
	}

	reward, err := client.CreateCustomReward("274637212", &customReward)
	require.NoError(t, err)
	require.NotNil(t, reward)
	require.Len(t, reward.Data, 1)
	require.Equal(t, "torpedo09", reward.Data[0].BroadcasterName)
	require.Equal(t, "torpedo09", reward.Data[0].BroadcasterLogin)
	require.Equal(t, "274637212", reward.Data[0].BroadcasterId)
	require.Equal(t, "afaa7e34-6b17-49f0-a19a-d1e76eaaf673", reward.Data[0].ID)
	require.Nil(t, reward.Data[0].Image)
	require.Equal(t, "#00E5CB", reward.Data[0].BackgroundColor)
	require.True(t, reward.Data[0].IsEnabled)
	require.Equal(t, cost, reward.Data[0].Cost)
	require.Equal(t, title, reward.Data[0].Title)
	require.Equal(t, "", reward.Data[0].Prompt)
	require.False(t, reward.Data[0].IsUserInputRequired)
	require.False(t, reward.Data[0].MaxPerStreamSetting.IsEnabled)
	require.Equal(t, 0, reward.Data[0].MaxPerStreamSetting.MaxPerStream)
	require.False(t, reward.Data[0].MaxPerUserPerStreamSetting.IsEnabled)
	require.Equal(t, 0, reward.Data[0].MaxPerUserPerStreamSetting.MaxPerUserPerStream)
	require.False(t, reward.Data[0].GlobalCooldownSetting.IsEnabled)
	require.Equal(t, 0, reward.Data[0].GlobalCooldownSetting.GlobalCooldownSeconds)
	require.False(t, reward.Data[0].IsPaused)
	require.True(t, reward.Data[0].IsInStock)
	require.NotNil(t, reward.Data[0].DefaultImage)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png", reward.Data[0].DefaultImage.Url1X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png", reward.Data[0].DefaultImage.Url2X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png", reward.Data[0].DefaultImage.Url4X)
	require.False(t, reward.Data[0].ShouldRedemptionsSkipRequestQueue)
	require.Zero(t, reward.Data[0].RedemptionsRedeemedCurrentStream)
	require.Nil(t, reward.Data[0].CooldownExpiresAt)
}

func TestClient_CreateCustomRewardReturnsError(t *testing.T) {
	client := Client{}

	reward, err := client.CreateCustomReward("", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.CreateCustomReward("    ", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.CreateCustomReward("		", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.CreateCustomReward("1234", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, custom reward can't be nil"})
	require.Nil(t, reward)
}

func TestClient_CreateClipReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/clips/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "274637212", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"FiveWordsForClipSlug\",\"edit_url\":\"http://clips.twitch.tv/FiveWordsForClipSlug/edit\"}]}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	editURL, id, err := client.CreateClip("274637212", false)
	require.NoError(t, err)
	require.Equal(t, "http://clips.twitch.tv/FiveWordsForClipSlug/edit", editURL)
	require.Equal(t, "FiveWordsForClipSlug", id)
}

func TestClient_CreateClipReturnsError(t *testing.T) {
	client := Client{}

	editURL, id, err := client.CreateClip("", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Zero(t, editURL)
	require.Zero(t, id)

	editURL, id, err = client.CreateClip("    ", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Zero(t, editURL)
	require.Zero(t, id)

	editURL, id, err = client.CreateClip("		", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Zero(t, editURL)
	require.Zero(t, id)
}

func TestClient_CreatePollReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/polls/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		poll := new(CreatePoll)
		err = json.Unmarshal(body, poll)
		require.NoError(t, err)

		require.Equal(t, "141981764", poll.BroadcasterID)
		require.Equal(t, "Heads or Tails?", poll.Title)
		require.Equal(t, "Heads", poll.Choices[0].Title)
		require.Equal(t, "Tails", poll.Choices[1].Title)
		require.Equal(t, 1800, poll.Duration)
		require.Equal(t, 100, *poll.ChannelPointsPerVote)
		require.True(t, *poll.ChannelPointsVotingEnabled)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"id\":\"ed961efd-8a3f-4cf5-a9d0-e616c590cd2a\",\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Heads or Tails?\",\"choices\":[{\"id\":\"4c123012-1351-4f33-84b7-43856e7a0f47\",\"title\":\"Heads\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0},{\"id\":\"279087e3-54a7-467e-bcd0-c1393fcea4f0\",\"title\":\"Tails\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0}],\"bits_voting_enabled\":false,\"bits_per_vote\":0,\"channel_points_voting_enabled\":true,\"channel_points_per_vote\":100,\"status\":\"ACTIVE\",\"duration\":1800,\"started_at\":\"2021-03-19T06:08:33.871278372Z\"}]}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	channelPointEnabled := true
	channelPointPerVote := 100

	newPoll := CreatePoll{
		BroadcasterID: "141981764",
		Title:         "Heads or Tails?",
		Choices: []Outcome{
			{
				Title: "Heads",
			},
			{
				Title: "Tails",
			},
		},
		ChannelPointsVotingEnabled: &channelPointEnabled,
		ChannelPointsPerVote:       &channelPointPerVote,
		Duration:                   1800,
	}

	poll, err := client.CreatePoll(&newPoll)
	require.NoError(t, err)
	require.NotNil(t, poll)
	require.Nil(t, poll.Pagination)
	require.Equal(t, "ed961efd-8a3f-4cf5-a9d0-e616c590cd2a", poll.Data[0].ID)
	require.Equal(t, "141981764", poll.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", poll.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", poll.Data[0].BroadcasterLogin)
	require.Equal(t, "Heads or Tails?", poll.Data[0].Title)
	require.Len(t, poll.Data[0].Choices, 2)
	require.False(t, poll.Data[0].BitsVotingEnabled)
	require.Zero(t, poll.Data[0].BitsPerVote)
	require.True(t, poll.Data[0].ChannelPointsVotingEnabled)
	require.Equal(t, 100, poll.Data[0].ChannelPointsPerVote)
	require.Equal(t, "ACTIVE", poll.Data[0].Status)
	require.Equal(t, 1800, poll.Data[0].Duration)
	require.Equal(t, 2021, poll.Data[0].StartedAt.Year())
}

func TestClient_CreatePollReturnsError(t *testing.T) {
	client := Client{}

	poll, err := client.CreatePoll(nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, poll can't be nil"})
	require.Nil(t, poll)
}

func TestClient_CreatePredictionReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/predictions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		prediction := new(CreatePrediction)
		err = json.Unmarshal(body, prediction)
		require.NoError(t, err)
		require.Equal(t, "141981764", prediction.BroadcasterId)
		require.Equal(t, "Any leeks in the stream?", prediction.Title)
		require.Len(t, prediction.Choices, 2)
		require.Equal(t, "Yes, give it time.", prediction.Choices[0].Title)
		require.Equal(t, "Definitely not.", prediction.Choices[1].Title)
		require.Equal(t, 120, prediction.PredictionWindow)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"id\":\"bc637af0-7766-4525-9308-4112f4cbf178\",\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Any leeks in the stream?\",\"winning_outcome_id\":null,\"outcomes\":[{\"id\":\"73085848-a94d-4040-9d21-2cb7a89374b7\",\"title\":\"Yes, give it time.\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"BLUE\"},{\"id\":\"906b70ba-1f12-47ea-9e95-e5f93d20e9cc\",\"title\":\"Definitely not.\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"PINK\"}],\"prediction_window\":120,\"status\":\"ACTIVE\",\"created_at\":\"2021-04-28T17:11:22.595914172Z\",\"ended_at\":null,\"locked_at\":null}]}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	newPredictions := CreatePrediction{
		BroadcasterId: "141981764",
		Title:         "Any leeks in the stream?",
		Choices: []Outcome{
			{
				Title: "Yes, give it time.",
			},
			{
				Title: "Definitely not.",
			},
		},
		PredictionWindow: 120,
	}

	prediction, err := client.CreatePrediction(&newPredictions)
	require.NoError(t, err)
	require.NotNil(t, prediction)
	require.Nil(t, prediction.Pagination)
	require.Len(t, prediction.Data, 1)
	require.Equal(t, "bc637af0-7766-4525-9308-4112f4cbf178", prediction.Data[0].ID)
	require.Equal(t, "141981764", prediction.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", prediction.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", prediction.Data[0].BroadcasterLogin)
	require.Equal(t, "Any leeks in the stream?", prediction.Data[0].Title)
	require.Zero(t, prediction.Data[0].WinningOutcomeId)
	require.Len(t, prediction.Data[0].Outcomes, 2)
	require.Equal(t, 120, prediction.Data[0].PredictionWindow)
	require.Equal(t, "ACTIVE", prediction.Data[0].Status)
	require.Equal(t, 2021, prediction.Data[0].CreatedAt.Year())
	require.Nil(t, prediction.Data[0].EndedAt)
	require.Nil(t, prediction.Data[0].LockedAt)

	blueOutcome := prediction.Data[0].Outcomes[0]
	pinkOutcome := prediction.Data[0].Outcomes[1]

	require.Equal(t, "73085848-a94d-4040-9d21-2cb7a89374b7", blueOutcome.ID)
	require.Equal(t, "906b70ba-1f12-47ea-9e95-e5f93d20e9cc", pinkOutcome.ID)
	require.Equal(t, "Yes, give it time.", blueOutcome.Title)
	require.Equal(t, "Definitely not.", pinkOutcome.Title)
	require.Equal(t, 0, blueOutcome.Users)
	require.Equal(t, 0, pinkOutcome.Users)
	require.Equal(t, 0, blueOutcome.ChannelPoints)
	require.Equal(t, 0, pinkOutcome.ChannelPoints)
	require.Nil(t, blueOutcome.TopPredictors)
	require.Nil(t, pinkOutcome.TopPredictors)
	require.Equal(t, "BLUE", blueOutcome.Color)
	require.Equal(t, "PINK", pinkOutcome.Color)
}

func TestClient_CreatePredictionReturnsError(t *testing.T) {
	client := Client{}

	predictions, err := client.CreatePrediction(nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, prediction can't be nil"})
	require.Nil(t, predictions)
}

func TestClient_CreateChannelStreamScheduleSegmentReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/segment/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		prediction := new(StreamScheduleSegmentUpdate)
		err = json.Unmarshal(body, prediction)
		require.NoError(t, err)
		require.False(t, *prediction.IsRecurring)
		require.Equal(t, "America/New_York", *prediction.Timezone)
		require.Equal(t, "60", *prediction.Duration)
		require.Equal(t, "509670", *prediction.CategoryID)
		require.Equal(t, "TwitchDev Monthly Update // July 1, 2021", *prediction.Title)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":{\"segments\":[{\"id\":\"eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=\",\"start_time\":\"2021-07-01T18:00:00Z\",\"end_time\":\"2021-07-01T19:00:00Z\",\"title\":\"TwitchDev Monthly Update // July 1, 2021\",\"canceled_until\":null,\"category\":{\"id\":\"509670\",\"name\":\"Science & Technology\"},\"is_recurring\":false}],\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"vacation\":null}}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	startTime, err := time.Parse(time.RFC3339, "2021-07-01T18:00:00Z")
	require.NoError(t, err)
	tz := "America/New_York"
	recurring := false
	duration := "60"
	category := "509670"
	title := "TwitchDev Monthly Update // July 1, 2021"
	segment := StreamScheduleSegmentUpdate{
		StartTime:   &startTime,
		Timezone:    &tz,
		IsRecurring: &recurring,
		Duration:    &duration,
		CategoryID:  &category,
		Title:       &title,
	}

	schedule, err := client.CreateChannelStreamScheduleSegment("141981764", &segment)
	require.NoError(t, err)
	require.NotNil(t, schedule)
	require.Nil(t, schedule.Pagination)
	require.Len(t, schedule.Data.Segments, 1)
	require.Equal(t, "141981764", schedule.Data.BroadcasterId)
	require.Equal(t, "TwitchDev", schedule.Data.BroadcasterName)
	require.Equal(t, "twitchdev", schedule.Data.BroadcasterLogin)
	require.Nil(t, schedule.Data.Vacation)

	seg := schedule.Data.Segments[0]
	require.Equal(t, "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=", seg.ID)
	require.Equal(t, startTime, seg.StartTime)
	require.Equal(t, startTime.Add(time.Hour), seg.EndTime)
	require.Equal(t, "TwitchDev Monthly Update // July 1, 2021", seg.Title)
	require.Nil(t, seg.CanceledUntil)
	require.Equal(t, "509670", seg.Category.Id)
	require.Equal(t, "Science & Technology", seg.Category.Name)
	require.False(t, seg.IsRecurring)
}

func TestClient_CreateChannelStreamScheduleSegmentReturnsError(t *testing.T) {
	client := Client{}

	schedule, err := client.CreateChannelStreamScheduleSegment("", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.CreateChannelStreamScheduleSegment("    ", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.CreateChannelStreamScheduleSegment("	", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.CreateChannelStreamScheduleSegment("1234", nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, segment can't be nil"})
	require.Nil(t, schedule)
}

func TestClient_CreateUserFollowsReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/follows/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		bodyMap := make(map[string]string)
		err = json.Unmarshal(body, &bodyMap)
		require.NoError(t, err)
		require.Equal(t, "57059344", bodyMap["from_id"])
		require.Equal(t, "41245072", bodyMap["to_id"])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	followed, err := client.CreateUserFollows("57059344", "41245072", false)
	require.NoError(t, err)
	require.True(t, followed)
}

func TestClient_CreateUserFollowsReturnsError(t *testing.T) {
	client := Client{}

	followed, err := client.CreateUserFollows("", "123", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, from id can't be blank"})
	require.False(t, followed)

	followed, err = client.CreateUserFollows("    ", "123", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, from id can't be blank"})
	require.False(t, followed)

	followed, err = client.CreateUserFollows("	", "123", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, from id can't be blank"})
	require.False(t, followed)

	followed, err = client.CreateUserFollows("123", "", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, to id can't be blank"})
	require.False(t, followed)

	followed, err = client.CreateUserFollows("123", "    ", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, to id can't be blank"})
	require.False(t, followed)

	followed, err = client.CreateUserFollows("123", "	", false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, to id can't be blank"})
	require.False(t, followed)
}
