package helix

import (
	"encoding/json"
	"fmt"
	gotau "github.com/Team-TAU/tau-client-go"
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

func TestClient_PatchRequestReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PATCH", r.Method)

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

	shouldBeFalse, shouldBeNil, err := client.PatchRequest("channels", nil, nil)
	require.ErrorIs(t, err, gotau.AuthorizationError{})
	require.False(t, shouldBeFalse)
	require.Nil(t, shouldBeNil)
}

func TestClient_PatchRequestReturnsRateLimitError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PATCH", r.Method)

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

	shouldBeFalse, shouldBeNil, err := client.PatchRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, RateLimitError{}, err)
	rlErr := err.(RateLimitError)
	require.NotNil(t, rlErr.ResetTime())
	require.Equal(t, 2021, rlErr.ResetTime().Year())
	require.Equal(t, time.Month(6), rlErr.ResetTime().Month())
	require.Equal(t, 17, rlErr.ResetTime().Day())
	require.False(t, shouldBeFalse)
	require.Nil(t, shouldBeNil)
}

func TestClient_PatchRequestReturnsGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PATCH", r.Method)

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

	shouldBeFalse, shouldBeNil, err := client.PatchRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, gotau.GenericError{}, err)
	require.False(t, shouldBeFalse)
	require.Nil(t, shouldBeNil)
}

func TestClient_ModifyChannelInformationReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "PATCH", r.Method)

		bodyMap := make(map[string]string)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &bodyMap)
		require.NoError(t, err)
		require.Len(t, bodyMap, 2)
		require.Equal(t, "33214", bodyMap["game_id"])
		require.Equal(t, "en", bodyMap["broadcaster_language"])

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

	gameID := "33214"
	language := "en"
	changed, err := client.ModifyChannelInformation("141981764", &gameID, &language, nil, nil)
	require.NoError(t, err)
	require.True(t, changed)
}

func TestClient_ModifyChannelInformationReturnsError(t *testing.T) {
	client := Client{}

	changed, err := client.ModifyChannelInformation("", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	changed, err = client.ModifyChannelInformation("    ", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	changed, err = client.ModifyChannelInformation("		", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	changed, err = client.ModifyChannelInformation("1234", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, at least one parameter must be provided of gameID, language, title, and delay"})
	require.False(t, changed)
}

func TestClient_UpdateCustomRewardReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "274637212", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", r.URL.Query().Get("id"))
		require.Equal(t, "PATCH", r.Method)

		requestData := new(CustomRewardsUpdate)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, requestData)
		require.NoError(t, err)
		require.False(t, *requestData.IsEnabled)

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"broadcaster_name\":\"torpedo09\",\"broadcaster_login\":\"torpedo09\",\"broadcaster_id\":\"274637212\",\"id\":\"92af127c-7326-4483-a52b-b0da0be61c01\",\"image\":null,\"background_color\":\"#00E5CB\",\"is_enabled\":false,\"cost\":30000,\"title\":\"game analysis 2v2\",\"prompt\":\"\",\"is_user_input_required\":false,\"max_per_stream_setting\":{\"is_enabled\":true,\"max_per_stream\":60},\"max_per_user_per_stream_setting\":{\"is_enabled\":false,\"max_per_user_per_stream\":0},\"global_cooldown_setting\":{\"is_enabled\":false,\"global_cooldown_seconds\":0},\"is_paused\":false,\"is_in_stock\":false,\"default_image\":{\"url_1x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png\",\"url_2x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png\",\"url_4x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png\"},\"should_redemptions_skip_request_queue\":true,\"redemptions_redeemed_current_stream\":60,\"cooldown_expires_at\":null}]}")
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

	changeRequest := new(CustomRewardsUpdate)
	enabled := false
	changeRequest.IsEnabled = &enabled
	reward, err := client.UpdateCustomReward("274637212", "92af127c-7326-4483-a52b-b0da0be61c01", changeRequest)
	require.NoError(t, err)
	require.NotNil(t, reward)
	require.Len(t, reward.Data, 1)
	require.Equal(t, "torpedo09", reward.Data[0].BroadcasterName)
	require.Equal(t, "torpedo09", reward.Data[0].BroadcasterLogin)
	require.Equal(t, "274637212", reward.Data[0].BroadcasterId)
	require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", reward.Data[0].ID)
	require.Nil(t, reward.Data[0].Image)
	require.Equal(t, "#00E5CB", reward.Data[0].BackgroundColor)
	require.False(t, reward.Data[0].IsEnabled)
	require.Equal(t, 30000, reward.Data[0].Cost)
	require.Equal(t, "game analysis 2v2", reward.Data[0].Title)
	require.Zero(t, reward.Data[0].Prompt)
	require.False(t, reward.Data[0].IsUserInputRequired)
	require.True(t, reward.Data[0].MaxPerStreamSetting.IsEnabled)
	require.Equal(t, 60, reward.Data[0].MaxPerStreamSetting.MaxPerStream)
	require.False(t, reward.Data[0].MaxPerUserPerStreamSetting.IsEnabled)
	require.Equal(t, 0, reward.Data[0].MaxPerUserPerStreamSetting.MaxPerUserPerStream)
	require.False(t, reward.Data[0].GlobalCooldownSetting.IsEnabled)
	require.Equal(t, 0, reward.Data[0].GlobalCooldownSetting.GlobalCooldownSeconds)
	require.False(t, reward.Data[0].IsPaused)
	require.False(t, reward.Data[0].IsInStock)
	require.NotNil(t, reward.Data[0].DefaultImage)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png", reward.Data[0].DefaultImage.Url1X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png", reward.Data[0].DefaultImage.Url2X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png", reward.Data[0].DefaultImage.Url4X)
	require.True(t, reward.Data[0].ShouldRedemptionsSkipRequestQueue)
	require.Equal(t, 60, reward.Data[0].RedemptionsRedeemedCurrentStream)
	require.Nil(t, reward.Data[0].CooldownExpiresAt)
}

func TestClient_UpdateCustomRewardReturnsError(t *testing.T) {
	client := Client{}

	reward, err := client.UpdateCustomReward("", "123", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("    ", "123", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("		", "123", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("123", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("123", "    ", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("123", "		", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, reward)

	reward, err = client.UpdateCustomReward("123", "456", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, change can't be nil"})
	require.Nil(t, reward)
}

func TestClient_UpdateRedemptionStatusReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "274637212", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", r.URL.Query().Get("reward_id"))
		require.Equal(t, "17fa2df1-ad76-4804-bfa5-a40ef63efe63", r.URL.Query().Get("id"))
		require.Equal(t, "PATCH", r.Method)

		requestData := make(map[string]string)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &requestData)
		require.NoError(t, err)
		require.Equal(t, "CANCELED", requestData["status"])

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"broadcaster_name\":\"torpedo09\",\"broadcaster_login\":\"torpedo09\",\"broadcaster_id\":\"274637212\",\"id\":\"17fa2df1-ad76-4804-bfa5-a40ef63efe63\",\"user_id\":\"274637212\",\"user_name\":\"torpedo09\",\"user_login\":\"torpedo09\",\"user_input\":\"\",\"status\":\"CANCELED\",\"redeemed_at\":\"2020-07-01T18:37:32Z\",\"reward\":{\"id\":\"92af127c-7326-4483-a52b-b0da0be61c01\",\"title\":\"game analysis\",\"prompt\":\"\",\"cost\":50000}}]}")
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

	changeRequest := new(CustomRewardsUpdate)
	enabled := false
	changeRequest.IsEnabled = &enabled
	redemption, err := client.UpdateRedemptionStatus("274637212", "92af127c-7326-4483-a52b-b0da0be61c01",
		[]string{"17fa2df1-ad76-4804-bfa5-a40ef63efe63"}, "CANCELED")
	require.NoError(t, err)
	require.NotNil(t, redemption)
	require.Nil(t, redemption.Pagination)
	require.Len(t, redemption.Data, 1)
	require.Equal(t, "torpedo09", redemption.Data[0].BroadcasterName)
	require.Equal(t, "torpedo09", redemption.Data[0].BroadcasterLogin)
	require.Equal(t, "274637212", redemption.Data[0].BroadcasterID)
	require.Equal(t, "17fa2df1-ad76-4804-bfa5-a40ef63efe63", redemption.Data[0].ID)
	require.Equal(t, "274637212", redemption.Data[0].UserID)
	require.Equal(t, "torpedo09", redemption.Data[0].UserLogin)
	require.Equal(t, "torpedo09", redemption.Data[0].UserName)
	require.Zero(t, redemption.Data[0].UserInput)
	require.Equal(t, "CANCELED", redemption.Data[0].Status)
	require.Equal(t, 2020, redemption.Data[0].RedeemedAt.Year())
	require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", redemption.Data[0].Reward.ID)
	require.Equal(t, "game analysis", redemption.Data[0].Reward.Title)
	require.Zero(t, redemption.Data[0].Reward.Prompt)
	require.Equal(t, 50000, redemption.Data[0].Reward.Cost)
}

func TestClient_UpdateRedemptionStatusReturnsError(t *testing.T) {
	client := Client{}

	redemptions, err := client.UpdateRedemptionStatus("", "", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("    ", "", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("		", "", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "   ", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "		", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "456", nil, "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, redemptionIDs can't be empty or nil"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "456", make([]string, 0), "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, redemptionIDs can't be empty or nil"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "456", make([]string, 51), "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request,  maximum of 50 redemptionIDs, but you supplied 51"})
	require.Nil(t, redemptions)

	redemptions, err = client.UpdateRedemptionStatus("123", "456", make([]string, 25), "foo")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request,  status can only be one of FULFILLED or CANCELED"})
	require.Nil(t, redemptions)
}

func TestClient_UpdateChannelStreamScheduleReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/settings/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "true", r.URL.Query().Get("is_vacation_enabled"))
		require.Equal(t, "America/New_York", r.URL.Query().Get("timezone"))
		require.Equal(t, "PATCH", r.Method)

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

	vacationMode := true
	now := time.Now()
	nextDay := now.Add(time.Hour * 24)
	tz := "America/New_York"

	redemption, err := client.UpdateChannelStreamSchedule("141981764", &vacationMode, &now, &nextDay, &tz)
	require.NoError(t, err)
	require.True(t, redemption)
}

func TestClient_UpdateChannelStreamScheduleReturnsTrueVacationDisabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/settings/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "false", r.URL.Query().Get("is_vacation_enabled"))
		require.Equal(t, "PATCH", r.Method)

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

	vacationMode := false

	redemption, err := client.UpdateChannelStreamSchedule("141981764", &vacationMode, nil, nil, nil)
	require.NoError(t, err)
	require.True(t, redemption)
}

func TestClient_UpdateChannelStreamScheduleReturnsError(t *testing.T) {
	client := Client{}

	changed, err := client.UpdateChannelStreamSchedule("", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	changed, err = client.UpdateChannelStreamSchedule("    ", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	changed, err = client.UpdateChannelStreamSchedule("	", nil, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.False(t, changed)

	vacationEnabled := true
	now := time.Now()
	changed, err = client.UpdateChannelStreamSchedule("123", &vacationEnabled, nil, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, if vacationEnabled, vacationStartTime must be specified"})
	require.False(t, changed)

	changed, err = client.UpdateChannelStreamSchedule("123", &vacationEnabled, &now, nil, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, if vacationEnabled, vacationEndTime must be specified"})
	require.False(t, changed)

	changed, err = client.UpdateChannelStreamSchedule("123", &vacationEnabled, &now, &now, nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, if vacationEnabled, timezone must be specified"})
	require.False(t, changed)
}

func TestClient_UpdateChannelStreamScheduleSegmentReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/segment/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
			r.URL.Query().Get("id"))
		require.Equal(t, "PATCH", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":{\"segments\":[{\"id\":\"eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=\",\"start_time\":\"2021-07-01T18:00:00Z\",\"end_time\":\"2021-07-01T20:00:00Z\",\"title\":\"TwitchDev Monthly Update // July 1, 2021\",\"canceled_until\":null,\"category\":{\"id\":\"509670\",\"name\":\"Science & Technology\"},\"is_recurring\":false}],\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"vacation\":null}}")
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

	duration := "120"
	update := StreamScheduleSegmentUpdate{
		Duration: &duration,
	}

	segment, err := client.UpdateChannelStreamScheduleSegment("141981764",
		"eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=",
		&update)
	require.NoError(t, err)
	require.NotNil(t, segment)
	require.Nil(t, segment.Pagination)
	require.Len(t, segment.Data.Segments, 1)
	require.Equal(t, "141981764", segment.Data.BroadcasterId)
	require.Equal(t, "TwitchDev", segment.Data.BroadcasterName)
	require.Equal(t, "twitchdev", segment.Data.BroadcasterLogin)
	require.Nil(t, segment.Data.Vacation)

	segmentData := segment.Data.Segments[0]
	require.Equal(t, "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=", segmentData.ID)
	require.Equal(t, 2021, segmentData.StartTime.Year())
	require.Equal(t, 2021, segmentData.EndTime.Year())
	require.Equal(t, "TwitchDev Monthly Update // July 1, 2021", segmentData.Title)
	require.Nil(t, segmentData.CanceledUntil)
	require.False(t, segmentData.IsRecurring)
	require.Equal(t, "509670", segmentData.Category.Id)
	require.Equal(t, "Science & Technology", segmentData.Category.Name)
}

func TestClient_UpdateChannelStreamScheduleSegmentReturnsError(t *testing.T) {
	client := Client{}

	segment, err := client.UpdateChannelStreamScheduleSegment("", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("    ", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("		", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("1234", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, segment id can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("1234", "    ", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, segment id can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("1234", "	", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, segment id can't be blank"})
	require.Nil(t, segment)

	segment, err = client.UpdateChannelStreamScheduleSegment("1234", "5678", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, update can't be nil"})
	require.Nil(t, segment)
}

func TestClient_EndPollReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/polls/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PATCH", r.Method)

		bodyMap := make(map[string]string)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &bodyMap)
		require.NoError(t, err)
		require.Equal(t, "141981764", bodyMap["broadcaster_id"])
		require.Equal(t, "ed961efd-8a3f-4cf5-a9d0-e616c590cd2a", bodyMap["id"])
		require.Equal(t, "TERMINATED", bodyMap["status"])
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"id\":\"ed961efd-8a3f-4cf5-a9d0-e616c590cd2a\",\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Heads or Tails?\",\"choices\":[{\"id\":\"4c123012-1351-4f33-84b7-43856e7a0f47\",\"title\":\"Heads\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0},{\"id\":\"279087e3-54a7-467e-bcd0-c1393fcea4f0\",\"title\":\"Tails\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0}],\"bits_voting_enabled\":false,\"bits_per_vote\":0,\"channel_points_voting_enabled\":true,\"channel_points_per_vote\":100,\"status\":\"TERMINATED\",\"duration\":1800,\"started_at\":\"2021-03-19T06:08:33.871278372Z\",\"ended_at\":\"2021-03-19T06:11:26.746889614Z\"}]}")
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

	poll, err := client.EndPoll("141981764", "ed961efd-8a3f-4cf5-a9d0-e616c590cd2a", "TERMINATED")
	require.NoError(t, err)
	require.NotNil(t, poll)
	require.Nil(t, poll.Pagination)
	require.Len(t, poll.Data, 1)
	require.Equal(t, "ed961efd-8a3f-4cf5-a9d0-e616c590cd2a", poll.Data[0].ID)
	require.Equal(t, "141981764", poll.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", poll.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", poll.Data[0].BroadcasterLogin)
	require.Equal(t, "Heads or Tails?", poll.Data[0].Title)
	require.Len(t, poll.Data[0].Choices, 2)
	require.Equal(t, "4c123012-1351-4f33-84b7-43856e7a0f47", poll.Data[0].Choices[0].ID)
	require.Equal(t, "Heads", poll.Data[0].Choices[0].Title)
	require.Equal(t, "279087e3-54a7-467e-bcd0-c1393fcea4f0", poll.Data[0].Choices[1].ID)
	require.Equal(t, "Tails", poll.Data[0].Choices[1].Title)
	require.False(t, poll.Data[0].BitsVotingEnabled)
	require.Equal(t, 0, poll.Data[0].BitsPerVote)
	require.True(t, poll.Data[0].ChannelPointsVotingEnabled)
	require.Equal(t, 100, poll.Data[0].ChannelPointsPerVote)
	require.Equal(t, "TERMINATED", poll.Data[0].Status)
	require.Equal(t, 1800, poll.Data[0].Duration)
	require.Equal(t, 2021, poll.Data[0].StartedAt.Year())
	require.Equal(t, 2021, poll.Data[0].EndedAt.Year())
}

func TestClient_EndPollReturnsError(t *testing.T) {
	client := Client{}

	poll, err := client.EndPoll("", "", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("    ", "", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("		", "", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, poll id can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "    ", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, poll id can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "	", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, poll id can't be blank"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "5678", "")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be TERMINATED or ARCHIVED"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "5678", "    ")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be TERMINATED or ARCHIVED"})
	require.Nil(t, poll)

	poll, err = client.EndPoll("1234", "5678", "		")
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be TERMINATED or ARCHIVED"})
	require.Nil(t, poll)
}

func TestClient_EndPredictionReturns200(t *testing.T) {
	winner := "73085848-a94d-4040-9d21-2cb7a89374b7"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/predictions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PATCH", r.Method)

		bodyMap := make(map[string]string)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &bodyMap)
		require.NoError(t, err)
		require.Equal(t, "141981764", bodyMap["broadcaster_id"])
		require.Equal(t, "bc637af0-7766-4525-9308-4112f4cbf178", bodyMap["id"])
		require.Equal(t, "RESOLVED", bodyMap["status"])
		require.Equal(t, winner, bodyMap["winning_outcome_id"])
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, "{\"data\":[{\"id\":\"bc637af0-7766-4525-9308-4112f4cbf178\",\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Will we win all the games?\",\"winning_outcome_id\":\"73085848-a94d-4040-9d21-2cb7a89374b7\",\"outcomes\":[{\"id\":\"73085848-a94d-4040-9d21-2cb7a89374b7\",\"title\":\"yes\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"BLUE\"},{\"id\":\"86010b2e-9764-4136-9359-fd1c9c5a8033\",\"title\":\"no\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"PINK\"}],\"prediction_window\":120,\"status\":\"RESOLVED\",\"created_at\":\"2021-04-28T21:48:19.480371331Z\",\"ended_at\":\"2021-04-28T21:54:24.026833954Z\",\"locked_at\":\"2021-04-28T21:48:34.636685705Z\"}]}")
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

	prediction, err := client.EndPrediction("141981764", "bc637af0-7766-4525-9308-4112f4cbf178",
		"RESOLVED", &winner)
	require.NoError(t, err)
	require.NotNil(t, prediction)
	require.Nil(t, prediction.Pagination)
	require.Len(t, prediction.Data, 1)
	require.Equal(t, "bc637af0-7766-4525-9308-4112f4cbf178", prediction.Data[0].ID)
	require.Equal(t, "141981764", prediction.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", prediction.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", prediction.Data[0].BroadcasterLogin)
	require.Equal(t, "Will we win all the games?", prediction.Data[0].Title)
	require.Equal(t, winner, prediction.Data[0].WinningOutcomeId)
	require.Len(t, prediction.Data[0].Outcomes, 2)
	require.Equal(t, 120, prediction.Data[0].PredictionWindow)
	require.Equal(t, "RESOLVED", prediction.Data[0].Status)
	require.Equal(t, 2021, prediction.Data[0].CreatedAt.Year())
	require.Equal(t, 2021, prediction.Data[0].EndedAt.Year())
	require.Equal(t, 2021, prediction.Data[0].LockedAt.Year())

	outcomeYes := prediction.Data[0].Outcomes[0]
	outcomeNo := prediction.Data[0].Outcomes[1]

	require.Equal(t, "73085848-a94d-4040-9d21-2cb7a89374b7", outcomeYes.ID)
	require.Equal(t, "yes", outcomeYes.Title)
	require.Equal(t, 0, outcomeYes.Users)
	require.Equal(t, 0, outcomeYes.ChannelPoints)
	require.Nil(t, outcomeYes.TopPredictors)
	require.Equal(t, "BLUE", outcomeYes.Color)

	require.Equal(t, "86010b2e-9764-4136-9359-fd1c9c5a8033", outcomeNo.ID)
	require.Equal(t, "no", outcomeNo.Title)
	require.Equal(t, 0, outcomeNo.Users)
	require.Equal(t, 0, outcomeNo.ChannelPoints)
	require.Nil(t, outcomeNo.TopPredictors)
	require.Equal(t, "PINK", outcomeNo.Color)
}

func TestClient_EndPredictionReturnsError(t *testing.T) {
	client := Client{}

	prediction, err := client.EndPrediction("", "", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("    ", "", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("		", "", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, prediction id can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "    ", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, prediction id can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "	", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, prediction id can't be blank"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "5678", "", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be RESOLVED, CANCELED, or LOCKED"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "5678", "    ", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be RESOLVED, CANCELED, or LOCKED"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "5678", "		", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, status must either be RESOLVED, CANCELED, or LOCKED"})
	require.Nil(t, prediction)

	prediction, err = client.EndPrediction("1234", "5678", "RESOLVED", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{"invalid request, if status RESOLVED, winning outcome must be set"})
	require.Nil(t, prediction)
}
