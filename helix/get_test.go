package helix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestNewClientReturnsClient(t *testing.T) {
	client, err := NewClient("tau.example.com", 443, "abcdefg", true)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, "tau.example.com", client.hostname)
	require.Equal(t, 443, client.port)
	require.Equal(t, "abcdefg", client.token)
	require.True(t, client.hasSSL)
}

func TestClient_GetTwitchUsersReturnsError(t *testing.T) {
	client := Client{}
	users := make([]string, 101)
	ids := make([]string, 101)

	result, err := client.GetTwitchUsers(users, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, get users only supports a maximum of 100 total logins " +
		"and ids combined."})
	require.Nil(t, result)

	result, err = client.GetTwitchUsers(nil, ids)
	require.ErrorIs(t, err, BadRequestError{"invalid request, get users only supports a maximum of 100 total logins " +
		"and ids combined."})
	require.Nil(t, result)

	result, err = client.GetTwitchUsers(users, ids)
	require.ErrorIs(t, err, BadRequestError{"invalid request, get users only supports a maximum of 100 total logins " +
		"and ids combined."})
	require.Nil(t, result)
}

func TestClient_GetBitsLeaderboardReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetBitsLeaderboard(101, "", nil, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum value for count is 100 and you input 101"})
	require.Nil(t, result)

	result, err = client.GetBitsLeaderboard(-1, "", nil, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can not be negative"})
	require.Nil(t, result)

	result, err = client.GetBitsLeaderboard(0, "decade", nil, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, period can only all, day, week, month, or year, and you input decade"})
	require.Nil(t, result)
}

func TestClient_GetBitsLeaderboardReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/bits/leaderboard/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "2", r.URL.Query().Get("count"))
		require.Equal(t, "week", r.URL.Query().Get("period"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"user_id\":\"158010205\",\"user_login\":\"tundracowboy\",\"user_name\":\"TundraCowboy\",\"rank\":1,\"score\":12543},{\"user_id\":\"7168163\",\"user_login\":\"topramens\",\"user_name\":\"Topramens\",\"rank\":2,\"score\":6900}],\"date_range\":{\"started_at\":\"2018-02-05T08:00:00Z\",\"ended_at\":\"2018-02-12T08:00:00Z\"},\"total\":2}")
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

	results, err := client.GetBitsLeaderboard(2, "week", nil, "")
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 2)
	require.Equal(t, 2, results.Total)
	require.NotNil(t, results.DateRange)
	require.Equal(t, 5, results.DateRange.StartedAt.Day())
	require.Equal(t, 12, results.DateRange.EndedAt.Day())

	user := results.Data[0]
	require.Equal(t, "158010205", user.UserID)
	require.Equal(t, "tundracowboy", user.UserLogin)
	require.Equal(t, "TundraCowboy", user.UserName)
	require.Equal(t, 1, user.Rank)
	require.Equal(t, 12543, user.Score)
}

func TestClient_GetCheermotesReturns200NoBroadcaster(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/bits/cheermotes/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"prefix\":\"Cheer\",\"tiers\":[{\"min_bits\":1,\"id\":\"1\",\"color\":\"#979797\",\"images\":{\"dark\":{\"animated\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.gif\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/2.gif\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/3.gif\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/4.gif\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.5.gif\"},\"static\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.png\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/2.png\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/3.png\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/4.png\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.5.png\"}},\"light\":{\"animated\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.gif\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/2.gif\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/3.gif\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/4.gif\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.5.gif\"},\"static\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.png\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/2.png\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/3.png\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/4.png\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.5.png\"}}},\"can_cheer\":true,\"show_in_bits_card\":true}],\"type\":\"global_first_party\",\"order\":1,\"last_updated\":\"2018-05-22T00:06:04Z\",\"is_charitable\":false}]}")
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

	results, err := client.GetCheermotes("")
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 1)
	require.Equal(t, "Cheer", results.Data[0].Prefix)
	require.Equal(t, "global_first_party", results.Data[0].Type)
	require.Equal(t, 1, results.Data[0].Order)
	require.Equal(t, 2018, results.Data[0].LastUpdated.Year())
	require.False(t, results.Data[0].IsCharitable)
	require.Len(t, results.Data[0].Tiers, 1)
	tier := results.Data[0].Tiers[0]
	require.Equal(t, 1, tier.MinBits)
	require.Equal(t, "1", tier.ID)
	require.Equal(t, "#979797", tier.Color)
	require.True(t, tier.CanCheer)
	require.True(t, tier.ShowInBitsCard)
	require.NotNil(t, tier.Images)

	dark := tier.Images.Dark
	light := tier.Images.Light

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.gif", dark.Animated.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.5.gif", dark.Animated.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/2.gif", dark.Animated.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/3.gif", dark.Animated.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/4.gif", dark.Animated.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.png", dark.Static.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.5.png", dark.Static.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/2.png", dark.Static.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/3.png", dark.Static.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/4.png", dark.Static.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.gif", light.Animated.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.5.gif", light.Animated.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/2.gif", light.Animated.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/3.gif", light.Animated.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/4.gif", light.Animated.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.png", light.Static.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.5.png", light.Static.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/2.png", light.Static.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/3.png", light.Static.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/4.png", light.Static.Four)
}

func TestClient_GetCheermotesReturns200WithBroadcaster(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/bits/cheermotes/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"prefix\":\"Cheer\",\"tiers\":[{\"min_bits\":1,\"id\":\"1\",\"color\":\"#979797\",\"images\":{\"dark\":{\"animated\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.gif\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/2.gif\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/3.gif\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/4.gif\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.5.gif\"},\"static\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.png\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/2.png\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/3.png\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/4.png\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.5.png\"}},\"light\":{\"animated\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.gif\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/2.gif\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/3.gif\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/4.gif\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.5.gif\"},\"static\":{\"1\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.png\",\"2\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/2.png\",\"3\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/3.png\",\"4\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/4.png\",\"1.5\":\"https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.5.png\"}}},\"can_cheer\":true,\"show_in_bits_card\":true}],\"type\":\"global_first_party\",\"order\":1,\"last_updated\":\"2018-05-22T00:06:04Z\",\"is_charitable\":false}]}")
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

	results, err := client.GetCheermotes("141981764")
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 1)
	require.Equal(t, "Cheer", results.Data[0].Prefix)
	require.Equal(t, "global_first_party", results.Data[0].Type)
	require.Equal(t, 1, results.Data[0].Order)
	require.Equal(t, 2018, results.Data[0].LastUpdated.Year())
	require.False(t, results.Data[0].IsCharitable)
	require.Len(t, results.Data[0].Tiers, 1)
	tier := results.Data[0].Tiers[0]
	require.Equal(t, 1, tier.MinBits)
	require.Equal(t, "1", tier.ID)
	require.Equal(t, "#979797", tier.Color)
	require.True(t, tier.CanCheer)
	require.True(t, tier.ShowInBitsCard)
	require.NotNil(t, tier.Images)

	dark := tier.Images.Dark
	light := tier.Images.Light

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.gif", dark.Animated.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/1.5.gif", dark.Animated.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/2.gif", dark.Animated.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/3.gif", dark.Animated.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/animated/1/4.gif", dark.Animated.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.png", dark.Static.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/1.5.png", dark.Static.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/2.png", dark.Static.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/3.png", dark.Static.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/dark/static/1/4.png", dark.Static.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.gif", light.Animated.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/1.5.gif", light.Animated.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/2.gif", light.Animated.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/3.gif", light.Animated.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/animated/1/4.gif", light.Animated.Four)

	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.png", light.Static.One)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/1.5.png", light.Static.OnePointFive)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/2.png", light.Static.Two)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/3.png", light.Static.Three)
	require.Equal(t, "https://d3aqoihi2n8ty8.cloudfront.net/actions/cheer/light/static/1/4.png", light.Static.Four)
}

func TestClient_GetCheermotesReturns500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/bits/cheermotes/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusInternalServerError)
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

	results, err := client.GetCheermotes("141981764")
	require.Error(t, err)
	require.Nil(t, results)
}

func TestClient_GetTwitchUsersReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"141981764\",\"login\":\"twitchdev\",\"display_name\":\"TwitchDev\",\"type\":\"\",\"broadcaster_type\":\"partner\",\"description\":\"Supporting third-party developers building Twitch integrations from chatbots to game integrations.\",\"profile_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png\",\"offline_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png\",\"view_count\":5980557,\"email\":\"not-real@email.com\",\"created_at\":\"2016-12-14T20:32:28.894263Z\"}]}")
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

	users, err := client.GetTwitchUsers(nil, []string{"141981764"})
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Len(t, users.Data, 1)
	require.Equal(t, "141981764", users.Data[0].ID)
	require.Equal(t, "twitchdev", users.Data[0].Login)
	require.Equal(t, "TwitchDev", users.Data[0].DisplayName)
	require.Equal(t, "", users.Data[0].Type)
	require.Equal(t, "partner", users.Data[0].BroadcasterType)
	require.Contains(t, users.Data[0].Description, "Supporting third-party")
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png", users.Data[0].ProfileImageUrl)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png", users.Data[0].OfflineImageUrl)
	require.Equal(t, 5980557, users.Data[0].ViewCount)
	require.Equal(t, "not-real@email.com", users.Data[0].Email)
	require.Equal(t, 2016, users.Data[0].CreatedAt.Year())
}

func TestClient_GetChannelInformationReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetChannelInformation("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelInformation("   ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelInformation("	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)
}

func TestClient_GetChannelInformationReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_id\":\"141981764\",\"broadcaster_login\":\"twitchdev\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_language\":\"en\",\"game_id\":\"509670\",\"game_name\":\"Science & Technology\",\"title\":\"TwitchDev Monthly Update // May 6, 2021\",\"delay\":0}]}")
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

	info, err := client.GetChannelInformation("141981764")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Len(t, info.Data, 1)
	require.Equal(t, "141981764", info.Data[0].BroadcasterID)
	require.Equal(t, "twitchdev", info.Data[0].BroadcasterLogin)
	require.Equal(t, "TwitchDev", info.Data[0].BroadcasterName)
	require.Equal(t, "en", info.Data[0].BroadcasterLanguage)
	require.Equal(t, "509670", info.Data[0].GameID)
	require.Equal(t, "Science & Technology", info.Data[0].GameName)
	require.Contains(t, info.Data[0].Title, "TwitchDev Monthly Update")
	require.Equal(t, 0, info.Data[0].Delay)
}

func TestClient_GetChannelInformationReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
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

	info, err := client.GetChannelInformation("141981764")
	require.Error(t, err)
	require.Nil(t, info)
}

func TestClient_GetChannelInformationReturns500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusInternalServerError)
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

	info, err := client.GetChannelInformation("141981764")
	require.Error(t, err)
	require.Nil(t, info)
}

func TestClient_GetChannelEditorsReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetChannelEditors("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelEditors("   ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelEditors("	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)
}

func TestClient_GetChannelEditorsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/editors/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"user_id\":\"182891647\",\"user_name\":\"mauerbac\",\"created_at\":\"2019-02-15T21:19:50.380833Z\"},{\"user_id\":\"135093069\",\"user_name\":\"BlueLava\",\"created_at\":\"2018-03-07T16:28:29.872937Z\"}]}")
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

	editors, err := client.GetChannelEditors("141981764")
	require.NoError(t, err)
	require.NotNil(t, editors)
	require.Len(t, editors.Data, 2)
	require.Equal(t, "182891647", editors.Data[0].UserID)
	require.Equal(t, "mauerbac", editors.Data[0].UserName)
	require.Equal(t, 2019, editors.Data[0].CreatedAt.Year())
}

func TestClient_GetChannelEditorsReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/editors/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
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

	editors, err := client.GetChannelEditors("141981764")
	require.Error(t, err)
	require.Nil(t, editors)
}

func TestClient_GetChannelEditorsReturns401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
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

	info, err := client.GetChannelInformation("141981764")
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.Nil(t, info)
}

func TestClient_GetCustomRewardsReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetCustomRewards("", nil, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewards("   ", nil, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewards("	", nil, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	rewardIDs := make([]string, 51)
	result, err = client.GetCustomRewards("12345", rewardIDs, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum number of rewardIDs is 50, but you input 51"})
	require.Nil(t, result)
}

func TestClient_GetCustomRewardsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_name\":\"torpedo09\",\"broadcaster_id\":\"274637212\",\"id\":\"92af127c-7326-4483-a52b-b0da0be61c01\",\"image\":null,\"background_color\":\"#00E5CB\",\"is_enabled\":true,\"cost\":50000,\"title\":\"game analysis\",\"prompt\":\"\",\"is_user_input_required\":false,\"max_per_stream_setting\":{\"is_enabled\":false,\"max_per_stream\":0},\"max_per_user_per_stream_setting\":{\"is_enabled\":false,\"max_per_user_per_stream\":0},\"global_cooldown_setting\":{\"is_enabled\":false,\"global_cooldown_seconds\":0},\"is_paused\":false,\"is_in_stock\":true,\"default_image\":{\"url_1x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png\",\"url_2x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png\",\"url_4x\":\"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png\"},\"should_redemptions_skip_request_queue\":false,\"redemptions_redeemed_current_stream\":null,\"cooldown_expires_at\":null}]}")
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

	rewards, err := client.GetCustomRewards("141981764", nil, false)
	require.NoError(t, err)
	require.NotNil(t, rewards)
	require.Len(t, rewards.Data, 1)
	data := rewards.Data[0]
	require.Equal(t, "torpedo09", data.BroadcasterName)
	require.Equal(t, "274637212", data.BroadcasterId)
	require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", data.Id)
	require.Nil(t, data.Image)
	require.Equal(t, "#00E5CB", data.BackgroundColor)
	require.True(t, data.IsEnabled)
	require.Equal(t, 50000, data.Cost)
	require.Equal(t, "game analysis", data.Title)
	require.Equal(t, "", data.Prompt)
	require.False(t, data.IsUserInputRequired)
	require.Equal(t, 0, data.MaxPerStreamSetting.MaxPerStream)
	require.False(t, data.MaxPerStreamSetting.IsEnabled)
	require.Equal(t, 0, data.MaxPerUserPerStreamSetting.MaxPerUserPerStream)
	require.False(t, data.MaxPerUserPerStreamSetting.IsEnabled)
	require.Equal(t, 0, data.GlobalCooldownSetting.GlobalCooldownSeconds)
	require.False(t, data.GlobalCooldownSetting.IsEnabled)
	require.False(t, data.IsPaused)
	require.True(t, data.IsInStock)
	require.NotNil(t, data.DefaultImage)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png", data.DefaultImage.Url1X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png", data.DefaultImage.Url2X)
	require.Equal(t, "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png", data.DefaultImage.Url4X)
	require.False(t, data.ShouldRedemptionsSkipRequestQueue)
	require.Equal(t, 0, data.RedemptionsRedeemedCurrentStream)
	require.Nil(t, data.CooldownExpiresAt)
}

func TestClient_GetCustomRewardsReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
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

	rewards, err := client.GetCustomRewards("abcdefg", nil, false)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardsReturns401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
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

	rewards, err := client.GetCustomRewards("141981764", nil, false)
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardsReturns403(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusForbidden)
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

	rewards, err := client.GetCustomRewards("141981764", nil, false)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardsReturns404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusNotFound)
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

	rewards, err := client.GetCustomRewards("141981764", []string{"92af127c-7326-4483-a52b-b0da0be61c01"}, false)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardsReturns500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusInternalServerError)
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

	rewards, err := client.GetCustomRewards("141981764", []string{"92af127c-7326-4483-a52b-b0da0be61c01"}, false)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardRedemptionReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetCustomRewardRedemption("", "", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("    ", "", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("	", "", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "    ", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "		", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, rewardID can't be blank"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", make([]string, 51), "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum number of redemptionIDs is 50, but you input 51"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", nil, "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, if there are no redemptionIDs, status must be set"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", nil, "foo", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, status can only be UNFULFILLED, FILFILLED, or CANCELED, but you input foo"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", nil, "UNFULFILLED", "bar", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, sort can only be OLDEST or NEWEST, but you input bar"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", nil, "UNFULFILLED", "NEWEST", "", 51)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum result count is 50, you input 51"})
	require.Nil(t, result)

	result, err = client.GetCustomRewardRedemption("12345", "67890", nil, "UNFULFILLED", "NEWEST", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, result)
}

func TestClient_GetCustomRewardRedemptionReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "274637212", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_name\":\"torpedo09\",\"broadcaster_login\":\"torpedo09\",\"broadcaster_id\":\"274637212\",\"id\":\"17fa2df1-ad76-4804-bfa5-a40ef63efe63\",\"user_login\":\"torpedo09\",\"user_id\":\"274637212\",\"user_name\":\"torpedo09\",\"user_input\":\"\",\"status\":\"CANCELED\",\"redeemed_at\":\"2020-07-01T18:37:32Z\",\"reward\":{\"id\":\"92af127c-7326-4483-a52b-b0da0be61c01\",\"title\":\"game analysis\",\"prompt\":\"\",\"cost\":50000}}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6Ik1UZG1ZVEprWmpFdFlXUTNOaTAwT0RBMExXSm1ZVFV0WVRRd1pXWTJNMlZtWlRZelgxOHlNREl3TFRBM0xUQXhWREU0T2pNM09qTXlMakl6TXpFeU56RTFOMW89In19\"}}")
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

	redemptions, err := client.GetCustomRewardRedemption("274637212", "92af127c-7326-4483-a52b-b0da0be61c01", nil, "CANCELED", "", "", 0)
	require.NoError(t, err)
	require.NotNil(t, redemptions)
	require.Len(t, redemptions.Data, 1)
	data := redemptions.Data[0]
	require.Equal(t, "torpedo09", data.BroadcasterName)
	require.Equal(t, "torpedo09", data.BroadcasterLogin)
	require.Equal(t, "274637212", data.BroadcasterID)
	require.Equal(t, "17fa2df1-ad76-4804-bfa5-a40ef63efe63", data.ID)
	require.Equal(t, "torpedo09", data.UserName)
	require.Equal(t, "torpedo09", data.UserLogin)
	require.Equal(t, "274637212", data.UserID)
	require.Zero(t, data.UserInput)
	require.Equal(t, "CANCELED", data.Status)
	require.Equal(t, 2020, data.RedeemedAt.Year())
	require.NotNil(t, data.Reward)
	require.Equal(t, "92af127c-7326-4483-a52b-b0da0be61c01", data.Reward.ID)
	require.Equal(t, "game analysis", data.Reward.Title)
	require.Zero(t, data.Reward.Prompt)
	require.Equal(t, 50000, data.Reward.Cost)
	require.NotNil(t, redemptions.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6Ik1UZG1ZVEprWmpFdFlXUTNOaTAwT0RBMExXSm1ZVFV0WVRRd1pXWTJNMlZtWlRZelgxOHlNREl3TFRBM0xUQXhWREU0T2pNM09qTXlMakl6TXpFeU56RTFOMW89In19", redemptions.Pagination.Cursor)
}

func TestClient_GetCustomRewardRedemptionReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
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

	rewards, err := client.GetCustomRewardRedemption("abcdefg", "12345", nil, "CANCELED", "", "", 0)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardRedemptionReturns401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
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

	rewards, err := client.GetCustomRewardRedemption("abcdefg", "12345", nil, "CANCELED", "", "", 0)
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardRedemptionReturns403(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusForbidden)
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

	rewards, err := client.GetCustomRewardRedemption("abcdefg", "12345", nil, "CANCELED", "", "", 0)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardRedemptionReturns404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusNotFound)
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

	rewards, err := client.GetCustomRewardRedemption("abcdefg", "12345", nil, "CANCELED", "", "", 0)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetCustomRewardRedemptionReturns500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/redemptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcdefg", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusInternalServerError)
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

	rewards, err := client.GetCustomRewardRedemption("abcdefg", "12345", nil, "CANCELED", "", "", 0)
	require.Error(t, err)
	require.Nil(t, rewards)
}

func TestClient_GetChannelChatBadgesReturnsError(t *testing.T) {
	client := Client{}

	result, err := client.GetChannelChatBadges("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelChatBadges("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)

	result, err = client.GetChannelChatBadges("	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, result)
}

func TestClient_GetChannelChatBadgesReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "{\"data\":[{\"set_id\":\"bits\",\"versions\":[{\"id\":\"1\",\"image_url_1x\":\"https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/1\",\"image_url_2x\":\"https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/2\",\"image_url_4x\":\"https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/3\"}]},{\"set_id\":\"subscriber\",\"versions\":[{\"id\":\"0\",\"image_url_1x\":\"https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/1\",\"image_url_2x\":\"https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/2\",\"image_url_4x\":\"https://static-cdn.jtvnw.net/badges/v1/eb4a8a4c-eacd-4f5e-b9f2-394348310442/3\"}]}]}")
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

	badges, err := client.GetChannelChatBadges("12345")
	require.NoError(t, err)
	require.NotNil(t, badges)
	require.Len(t, badges.Data, 2)
	require.Equal(t, "bits", badges.Data[0].SetID)
	require.Len(t, badges.Data[0].Versions, 1)
	require.Equal(t, "1", badges.Data[0].Versions[0].ID)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/1", badges.Data[0].Versions[0].ImageUrl1X)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/2", badges.Data[0].Versions[0].ImageUrl2X)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/743a0f3b-84b3-450b-96a0-503d7f4a9764/3", badges.Data[0].Versions[0].ImageUrl4X)
}

func TestClient_GetChannelChatBadgesReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "abcd", r.URL.Query().Get("broadcaster_id"))
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

	badges, err := client.GetChannelChatBadges("abcd")
	require.Error(t, err)
	require.Nil(t, badges)
}

func TestClient_GetChannelChatBadgesReturns401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
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

	badges, err := client.GetChannelChatBadges("12345")
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.Nil(t, badges)
}

func TestClient_GetGlobalChatBadgesReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/global/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"set_id\":\"vip\",\"versions\":[{\"id\":\"1\",\"image_url_1x\":\"https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/1\",\"image_url_2x\":\"https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/2\",\"image_url_4x\":\"https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/3\"}]}]}")
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

	badges, err := client.GetGlobalChatBadges()
	require.NoError(t, err)
	require.NotNil(t, badges)

	require.Len(t, badges.Data, 1)
	require.Equal(t, "vip", badges.Data[0].SetID)
	require.Len(t, badges.Data[0].Versions, 1)
	require.Equal(t, "1", badges.Data[0].Versions[0].ID)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/1", badges.Data[0].Versions[0].ImageUrl1X)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/2", badges.Data[0].Versions[0].ImageUrl2X)
	require.Equal(t, "https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/3", badges.Data[0].Versions[0].ImageUrl4X)
}

func TestClient_GetGlobalChatBadgesReturns400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/global/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
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

	badges, err := client.GetGlobalChatBadges()
	require.Error(t, err)
	require.Nil(t, badges)
}

func TestClient_GetGlobalChatBadgesReturns401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/chat/badges/global/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
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

	badges, err := client.GetGlobalChatBadges()
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.Nil(t, badges)
}

func TestClient_GetClipsXReturnsError(t *testing.T) {
	client := Client{}

	clips, err := client.GetClipsByBroadcaster("", "", "", nil, nil, 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByBroadcaster("12345", "", "", nil, nil, 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count has a maximum of 100, but you supplied 101"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByBroadcaster("12345", "", "", nil, nil, -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count cannot be negative"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByGame("", "", "", nil, nil, 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, gameID can't be blank"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByGame("12345", "", "", nil, nil, 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count has a maximum of 100, but you supplied 101"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByGame("12345", "", "", nil, nil, -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count cannot be negative"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByID(nil, "", "", nil, nil, 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, clipID can't be empty"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByID(make([]string, 101), "", "", nil, nil, 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, clipID has a maximum of 100 but you supplied 101 ids"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByID(make([]string, 42), "", "", nil, nil, 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count has a maximum of 100, but you supplied 101"})
	require.Nil(t, clips)

	clips, err = client.GetClipsByID(make([]string, 42), "", "", nil, nil, -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count cannot be negative"})
	require.Nil(t, clips)
}

func TestClient_GetClipsByBroadcasterReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/clips/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "1234", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "5", r.URL.Query().Get("first"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"RandomClip1\",\"url\":\"https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage\",\"embed_url\":\"https://clips.twitch.tv/embed?clip=RandomClip1\",\"broadcaster_id\":\"1234\",\"broadcaster_name\":\"JJ\",\"creator_id\":\"123456\",\"creator_name\":\"MrMarshall\",\"video_id\":\"1234567\",\"game_id\":\"33103\",\"language\":\"en\",\"title\":\"random1\",\"view_count\":10,\"created_at\":\"2017-11-30T22:34:18Z\",\"thumbnail_url\":\"https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg\",\"duration\":12.9}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjoiIn0\"}}")
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

	clips, err := client.GetClipsByBroadcaster("1234", "", "", nil, nil, 5)
	require.NoError(t, err)
	require.NotNil(t, clips)

	require.Len(t, clips.Data, 1)
	require.NotNil(t, clips.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjoiIn0", clips.Pagination.Cursor)

	clip := clips.Data[0]
	require.Equal(t, "RandomClip1", clip.ID)
	require.Equal(t, "https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage", clip.Url)
	require.Equal(t, "https://clips.twitch.tv/embed?clip=RandomClip1", clip.EmbedUrl)
	require.Equal(t, "1234", clip.BroadcasterID)
	require.Equal(t, "JJ", clip.BroadcasterName)
	require.Equal(t, "123456", clip.CreatorID)
	require.Equal(t, "MrMarshall", clip.CreatorName)
	require.Equal(t, "1234567", clip.VideoID)
	require.Equal(t, "33103", clip.GameID)
	require.Equal(t, "en", clip.Language)
	require.Equal(t, "random1", clip.Title)
	require.Equal(t, 10, clip.ViewCount)
	require.Equal(t, 2017, clip.CreatedAt.Year())
	require.Equal(t, "https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg", clip.ThumbnailUrl)
	require.Equal(t, 12.9, clip.Duration)
}

func TestClient_GetClipsByGameReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/clips/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "1234", r.URL.Query().Get("game_id"))
		require.Equal(t, "5", r.URL.Query().Get("first"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"RandomClip1\",\"url\":\"https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage\",\"embed_url\":\"https://clips.twitch.tv/embed?clip=RandomClip1\",\"broadcaster_id\":\"1234\",\"broadcaster_name\":\"JJ\",\"creator_id\":\"123456\",\"creator_name\":\"MrMarshall\",\"video_id\":\"1234567\",\"game_id\":\"33103\",\"language\":\"en\",\"title\":\"random1\",\"view_count\":10,\"created_at\":\"2017-11-30T22:34:18Z\",\"thumbnail_url\":\"https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg\",\"duration\":12.9}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjoiIn0\"}}")
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

	clips, err := client.GetClipsByGame("1234", "", "", nil, nil, 5)
	require.NoError(t, err)
	require.NotNil(t, clips)

	require.Len(t, clips.Data, 1)
	require.NotNil(t, clips.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjoiIn0", clips.Pagination.Cursor)

	clip := clips.Data[0]
	require.Equal(t, "RandomClip1", clip.ID)
	require.Equal(t, "https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage", clip.Url)
	require.Equal(t, "https://clips.twitch.tv/embed?clip=RandomClip1", clip.EmbedUrl)
	require.Equal(t, "1234", clip.BroadcasterID)
	require.Equal(t, "JJ", clip.BroadcasterName)
	require.Equal(t, "123456", clip.CreatorID)
	require.Equal(t, "MrMarshall", clip.CreatorName)
	require.Equal(t, "1234567", clip.VideoID)
	require.Equal(t, "33103", clip.GameID)
	require.Equal(t, "en", clip.Language)
	require.Equal(t, "random1", clip.Title)
	require.Equal(t, 10, clip.ViewCount)
	require.Equal(t, 2017, clip.CreatedAt.Year())
	require.Equal(t, "https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg", clip.ThumbnailUrl)
	require.Equal(t, 12.9, clip.Duration)
}

func TestClient_GetClipsByIDReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/clips/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "1234", r.URL.Query().Get("id"))
		require.Equal(t, "5", r.URL.Query().Get("first"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"RandomClip1\",\"url\":\"https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage\",\"embed_url\":\"https://clips.twitch.tv/embed?clip=RandomClip1\",\"broadcaster_id\":\"1234\",\"broadcaster_name\":\"JJ\",\"creator_id\":\"123456\",\"creator_name\":\"MrMarshall\",\"video_id\":\"1234567\",\"game_id\":\"33103\",\"language\":\"en\",\"title\":\"random1\",\"view_count\":10,\"created_at\":\"2017-11-30T22:34:18Z\",\"thumbnail_url\":\"https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg\",\"duration\":12.9}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjoiIn0\"}}")
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

	clips, err := client.GetClipsByID([]string{"1234"}, "", "", nil, nil, 5)
	require.NoError(t, err)
	require.NotNil(t, clips)

	require.Len(t, clips.Data, 1)
	require.NotNil(t, clips.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjoiIn0", clips.Pagination.Cursor)

	clip := clips.Data[0]
	require.Equal(t, "RandomClip1", clip.ID)
	require.Equal(t, "https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage", clip.Url)
	require.Equal(t, "https://clips.twitch.tv/embed?clip=RandomClip1", clip.EmbedUrl)
	require.Equal(t, "1234", clip.BroadcasterID)
	require.Equal(t, "JJ", clip.BroadcasterName)
	require.Equal(t, "123456", clip.CreatorID)
	require.Equal(t, "MrMarshall", clip.CreatorName)
	require.Equal(t, "1234567", clip.VideoID)
	require.Equal(t, "33103", clip.GameID)
	require.Equal(t, "en", clip.Language)
	require.Equal(t, "random1", clip.Title)
	require.Equal(t, 10, clip.ViewCount)
	require.Equal(t, 2017, clip.CreatedAt.Year())
	require.Equal(t, "https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg", clip.ThumbnailUrl)
	require.Equal(t, 12.9, clip.Duration)
}

func TestClient_GetEventSubSubscriptionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/eventsub/subscriptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"total\":2,\"data\":[{\"id\":\"26b1c993-bfcf-44d9-b876-379dacafe75a\",\"status\":\"enabled\",\"type\":\"streams.online\",\"version\":\"1\",\"condition\":{\"broadcaster_user_id\":\"1234\"},\"created_at\":\"2020-11-10T20:08:33Z\",\"transport\":{\"method\":\"webhook\",\"callback\":\"https://this-is-a-callback.com\"},\"cost\":1},{\"id\":\"35016908-41ff-33ce-7879-61b8dfc2ee16\",\"status\":\"webhook-callback-verification-pending\",\"type\":\"users.update\",\"version\":\"1\",\"condition\":{\"user_id\":\"1234\"},\"created_at\":\"2020-11-10T20:31:52Z\",\"transport\":{\"method\":\"webhook\",\"callback\":\"https://this-is-a-callback.com\"},\"cost\":0}],\"total_cost\":1,\"max_total_cost\":10000,\"pagination\":{}}")
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

	subs, err := client.GetEventSubSubscriptions("", "")
	require.NoError(t, err)
	require.NotNil(t, subs)
	require.Equal(t, 2, subs.Total)
	require.Equal(t, 1, subs.TotalCost)
	require.Equal(t, 10000, subs.MaxTotalCost)
	require.Len(t, subs.Data, 2)
	data := subs.Data[0]
	require.NotNil(t, data)
	require.Equal(t, "26b1c993-bfcf-44d9-b876-379dacafe75a", data.ID)
	require.Equal(t, "enabled", data.Status)
	require.Equal(t, "streams.online", data.Type)
	require.Equal(t, "1", data.Version)
	require.Equal(t, 2020, data.CreatedAt.Year())
	require.Equal(t, 1, data.Cost)
	require.NotNil(t, data.Condition)
	require.Equal(t, "1234", data.Condition.BroadcasterUserID)
	require.NotNil(t, data.Transport)
	require.Equal(t, "webhook", data.Transport.Method)
	require.Equal(t, "https://this-is-a-callback.com", data.Transport.Callback)

	require.Equal(t, "1234", subs.Data[1].Condition.UserID)
}

func TestClient_GetTopGamesReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/games/top/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "1", r.URL.Query().Get("first"))
		require.Equal(t, "first=1", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"493057\",\"name\":\"PLAYERUNKNOWN'S BATTLEGROUNDS\",\"box_art_url\":\"https://static-cdn.jtvnw.net/ttv-boxart/PLAYERUNKNOWN%27S%20BATTLEGROUNDS-{width}x{height}.jpg\"}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MjB9fQ==\"}}")
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

	games, err := client.GetTopGames("", "", 1)
	require.NoError(t, err)
	require.NotNil(t, games)

	require.Len(t, games.Data, 1)
	require.Equal(t, "493057", games.Data[0].ID)
	require.Equal(t, "PLAYERUNKNOWN'S BATTLEGROUNDS", games.Data[0].Name)
	require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/PLAYERUNKNOWN%27S%20BATTLEGROUNDS-{width}x{height}.jpg", games.Data[0].BoxArtUrl)

	require.NotNil(t, games.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MjB9fQ==", games.Pagination.Cursor)
}

func TestClient_GetTopGamesReturnsError(t *testing.T) {
	client := Client{}
	games, err := client.GetTopGames("", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, games)

	games, err = client.GetTopGames("", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, games)
}

func TestClient_GetGamesReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/games/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "33214", r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"box_art_url\":\"https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-52x72.jpg\",\"id\":\"33214\",\"name\":\"Fortnite\"}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN\"}}")
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

	games, err := client.GetGames([]string{"33214"}, []string{})
	require.NoError(t, err)
	require.NotNil(t, games)

	require.Len(t, games.Data, 1)
	require.Equal(t, "33214", games.Data[0].ID)
	require.Equal(t, "Fortnite", games.Data[0].Name)
	require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-52x72.jpg", games.Data[0].BoxArtUrl)

	require.NotNil(t, games.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN", games.Pagination.Cursor)
}

func TestClient_GetGamesReturnsError(t *testing.T) {
	client := Client{}
	games, err := client.GetGames(nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, either id or names is necessary"})
	require.Nil(t, games)

	games, err = client.GetGames(make([]string, 101), nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum number of games and ids is 100, you input 101"})
	require.Nil(t, games)

	games, err = client.GetGames(nil, make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum number of games and ids is 100, you input 101"})
	require.Nil(t, games)
}

func TestClient_GetHypeTrainEventsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/hypetrain/events/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"event_type\":\"hypetrain.progression\",\"event_timestamp\":\"2020-04-24T20:07:24Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"270954519\",\"cooldown_end_time\":\"2020-04-24T20:13:21.003802269Z\",\"expires_at\":\"2020-04-24T20:12:21.003802269Z\",\"goal\":1800,\"id\":\"70f0c7d8-ff60-4c50-b138-f3a352833b50\",\"last_contribution\":{\"total\":200,\"type\":\"BITS\",\"user\":\"134247454\"},\"level\":2,\"started_at\":\"2020-04-24T20:05:47.30473127Z\",\"top_contributions\":[{\"total\":600,\"type\":\"BITS\",\"user\":\"134247450\"}],\"total\":600}}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjI3MDk1NDUxOToxNTg3NzU4ODQ0OjFiMEFzYkluQ0haVzJTUUZRa0N6cU4wN0liMiJ9fQ\"}}")
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

	events, err := client.GetHypeTrainEvents("12345", 0, "", "")
	require.NoError(t, err)
	require.NotNil(t, events)
	require.Len(t, events.Data, 1)
	require.NotNil(t, events.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjI3MDk1NDUxOToxNTg3NzU4ODQ0OjFiMEFzYkluQ0haVzJTUUZRa0N6cU4wN0liMiJ9fQ", events.Pagination.Cursor)

	data := events.Data[0]
	require.Equal(t, "1b0AsbInCHZW2SQFQkCzqN07Ib2", data.ID)
	require.Equal(t, "hypetrain.progression", data.EventType)
	require.Equal(t, 2020, data.EventTimestamp.Year())
	require.Equal(t, "1.0", data.Version)
	require.NotNil(t, data.EventData)
	require.Equal(t, "270954519", data.EventData.BroadcasterId)
	require.Equal(t, 2020, data.EventData.CooldownEndTime.Year())
	require.Equal(t, 2020, data.EventData.ExpiresAt.Year())
	require.Equal(t, 1800, data.EventData.Goal)
	require.Equal(t, "70f0c7d8-ff60-4c50-b138-f3a352833b50", data.EventData.ID)
	require.Equal(t, 2, data.EventData.Level)
	require.Equal(t, 2020, data.EventData.StartedAt.Year())
	require.Equal(t, 600, data.EventData.Total)
	require.Equal(t, 200, data.EventData.LastContribution.Total)
	require.Equal(t, "BITS", data.EventData.LastContribution.Type)
	require.Equal(t, "134247454", data.EventData.LastContribution.User)
	require.Len(t, data.EventData.TopContributions, 1)
	require.Equal(t, 600, data.EventData.TopContributions[0].Total)
	require.Equal(t, "BITS", data.EventData.TopContributions[0].Type)
	require.Equal(t, "134247450", data.EventData.TopContributions[0].User)
}

func TestClient_GetHypeTrainEventsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetHypeTrainEvents("", 0, "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetHypeTrainEvents("    ", 0, "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetHypeTrainEvents("	", 0, "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetHypeTrainEvents("12345", 101, "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.GetHypeTrainEvents("12345", -1, "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)
}

func TestClient_GetBannedEventsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/moderation/banned/events/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"1IPFqAb0p0JncbPSTEPhx8JF1Sa\",\"event_type\":\"moderation.user.ban\",\"event_timestamp\":\"2019-03-13T15:55:14Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"racageneg\",\"broadcaster_name\":\"racageneg\",\"user_id\":\"424596340\",\"user_login\":\"quotrok\",\"user_name\":\"quotrok\",\"expires_at\":\"2021-06-11T01:35:43.04850111Z\"}},{\"id\":\"1IPFsDv5cs4mxfJ1s2O9Q5flf4Y\",\"event_type\":\"moderation.user.unban\",\"event_timestamp\":\"2019-03-13T15:55:30Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"racageneg\",\"broadcaster_name\":\"racageneg\",\"user_id\":\"424596340\",\"user_login\":\"quotrok\",\"user_name\":\"quotrok\",\"expires_at\":\"\"}},{\"id\":\"1IPFqmlu9W2q4mXXjULyM8zX0rb\",\"event_type\":\"moderation.user.ban\",\"event_timestamp\":\"2019-03-13T15:55:19Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"racageneg\",\"broadcaster_name\":\"racageneg\",\"user_id\":\"424596340\",\"user_login\":\"quotrok\",\"user_name\":\"quotrok\",\"expires_at\":\"\"}}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjE5OTYwNDI2MzoyMDIxMjA1MzE6MUlQRnFtbHU5VzJxNG1YWGpVTHlNOHpYMHJiIn19\"}}")
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

	events, err := client.GetBannedEvents("12345", "", "", 0)
	require.NoError(t, err)
	require.NotNil(t, events)
	require.Len(t, events.Data, 3)
	require.NotNil(t, events.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjE5OTYwNDI2MzoyMDIxMjA1MzE6MUlQRnFtbHU5VzJxNG1YWGpVTHlNOHpYMHJiIn19", events.Pagination.Cursor)

	data := events.Data[0]
	require.Equal(t, "1IPFqAb0p0JncbPSTEPhx8JF1Sa", data.ID)
	require.Equal(t, "moderation.user.ban", data.EventType)
	require.Equal(t, 2019, data.EventTimestamp.Year())
	require.Equal(t, "1.0", data.Version)
	require.NotNil(t, data.EventData)
	require.Equal(t, "198704263", data.EventData.BroadcasterID)
	require.Equal(t, "racageneg", data.EventData.BroadcasterName)
	require.Equal(t, "racageneg", data.EventData.BroadcasterLogin)
	require.Equal(t, "424596340", data.EventData.UserID)
	require.Equal(t, "quotrok", data.EventData.UserName)
	require.Equal(t, "quotrok", data.EventData.UserLogin)
	require.Equal(t, "2021-06-11T01:35:43.04850111Z", data.EventData.ExpiresAt)
}

func TestClient_GetBannedEventsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetBannedEvents("", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetBannedEvents("    ", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetBannedEvents("	", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetBannedEvents("12345", "", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.GetBannedEvents("12345", "", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)
}

func TestClient_GetModeratorsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/moderation/moderators/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"event_type\":\"hypetrain.progression\",\"event_timestamp\":\"2020-04-24T20:07:24Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"270954519\",\"cooldown_end_time\":\"2020-04-24T20:13:21.003802269Z\",\"expires_at\":\"2020-04-24T20:12:21.003802269Z\",\"goal\":1800,\"id\":\"70f0c7d8-ff60-4c50-b138-f3a352833b50\",\"last_contribution\":{\"total\":200,\"type\":\"BITS\",\"user\":\"134247454\"},\"level\":2,\"started_at\":\"2020-04-24T20:05:47.30473127Z\",\"top_contributions\":[{\"total\":600,\"type\":\"BITS\",\"user\":\"134247450\"}],\"total\":600}}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjI3MDk1NDUxOToxNTg3NzU4ODQ0OjFiMEFzYkluQ0haVzJTUUZRa0N6cU4wN0liMiJ9fQ\"}}")
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

	events, err := client.GetModerators("12345", nil, "", 0)
	require.NoError(t, err)
	require.NotNil(t, events)
}

func TestClient_GetModeratorsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetModerators("", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModerators("    ", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModerators("	", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModerators("12345", nil, "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.GetModerators("12345", nil, "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)

	event, err = client.GetModerators("12345", make([]string, 101), "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum user ids that can be supplied is 100 but you supplied 101"})
	require.Nil(t, event)
}

func TestClient_GetModeratorEventsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/moderation/moderators/events/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"1IVBTnDSUDApiBQW4UBcVTK4hPr\",\"event_type\":\"moderation.moderator.remove\",\"event_timestamp\":\"2019-03-15T18:18:14Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"aan22209\",\"broadcaster_name\":\"aan22209\",\"user_id\":\"423374343\",\"user_login\":\"glowillig\",\"user_name\":\"glowillig\"}},{\"id\":\"1IVIPQdYIEnD8nJ376qkASDzsj7\",\"event_type\":\"moderation.moderator.add\",\"event_timestamp\":\"2019-03-15T19:15:13Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"aan22209\",\"broadcaster_name\":\"aan22209\",\"user_id\":\"423374343\",\"user_login\":\"glowillig\",\"user_name\":\"glowillig\"}},{\"id\":\"1IVBTP7gG61oXLMu7fvnRhrpsro\",\"event_type\":\"moderation.moderator.remove\",\"event_timestamp\":\"2019-03-15T18:18:11Z\",\"version\":\"1.0\",\"event_data\":{\"broadcaster_id\":\"198704263\",\"broadcaster_login\":\"aan22209\",\"broadcaster_name\":\"aan22209\",\"user_id\":\"424596340\",\"user_login\":\"quotrok\",\"user_name\":\"quotrok\"}}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjEwMDQ3MzA2NDo4NjQwNjU3MToxSVZCVDFKMnY5M1BTOXh3d1E0dUdXMkJOMFcifX0\"}}")
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

	events, err := client.GetModeratorEvents("12345", nil, "", 0)
	require.NoError(t, err)
	require.NotNil(t, events)
	require.Len(t, events.Data, 3)
	require.NotNil(t, events.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IjEwMDQ3MzA2NDo4NjQwNjU3MToxSVZCVDFKMnY5M1BTOXh3d1E0dUdXMkJOMFcifX0", events.Pagination.Cursor)
	require.Equal(t, "1IVBTnDSUDApiBQW4UBcVTK4hPr", events.Data[0].ID)
	require.Equal(t, "moderation.moderator.remove", events.Data[0].EventType)
	require.Equal(t, 2019, events.Data[0].EventTimestamp.Year())
	require.Equal(t, "1.0", events.Data[0].Version)
	require.NotNil(t, events.Data[0].EventData)
	require.Equal(t, "198704263", events.Data[0].EventData.BroadcasterID)
	require.Equal(t, "aan22209", events.Data[0].EventData.BroadcasterLogin)
	require.Equal(t, "aan22209", events.Data[0].EventData.BroadcasterName)
	require.Equal(t, "423374343", events.Data[0].EventData.UserID)
	require.Equal(t, "glowillig", events.Data[0].EventData.UserLogin)
	require.Equal(t, "glowillig", events.Data[0].EventData.UserName)
}

func TestClient_GetModeratorEventsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetModeratorEvents("", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModeratorEvents("    ", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModeratorEvents("	", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetModeratorEvents("12345", nil, "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.GetModeratorEvents("12345", nil, "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)

	event, err = client.GetModeratorEvents("12345", make([]string, 101), "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum user ids that can be supplied is 100 but you supplied 101"})
	require.Nil(t, event)
}

func TestClient_GetPollsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/polls/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"ed961efd-8a3f-4cf5-a9d0-e616c590cd2a\",\"broadcaster_id\":\"55696719\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Heads or Tails?\",\"choices\":[{\"id\":\"4c123012-1351-4f33-84b7-43856e7a0f47\",\"title\":\"Heads\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0},{\"id\":\"279087e3-54a7-467e-bcd0-c1393fcea4f0\",\"title\":\"Tails\",\"votes\":0,\"channel_points_votes\":0,\"bits_votes\":0}],\"bits_voting_enabled\":false,\"bits_per_vote\":0,\"channel_points_voting_enabled\":false,\"channel_points_per_vote\":0,\"status\":\"ACTIVE\",\"duration\":1800,\"started_at\":\"2021-03-19T06:08:33.871278372Z\"}],\"pagination\":{}}")
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

	polls, err := client.GetPolls("12345", nil, "", 0)
	require.NoError(t, err)
	require.NotNil(t, polls)
	require.Len(t, polls.Data, 1)
	require.NotNil(t, polls.Pagination)
	require.Equal(t, "ed961efd-8a3f-4cf5-a9d0-e616c590cd2a", polls.Data[0].ID)
	require.Equal(t, "55696719", polls.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", polls.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", polls.Data[0].BroadcasterLogin)
	require.Equal(t, "Heads or Tails?", polls.Data[0].Title)
	require.Len(t, polls.Data[0].Choices, 2)
	require.False(t, polls.Data[0].BitsVotingEnabled)
	require.Equal(t, 0, polls.Data[0].BitsPerVote)
	require.False(t, polls.Data[0].ChannelPointsVotingEnabled)
	require.Equal(t, 0, polls.Data[0].ChannelPointsPerVote)
	require.Equal(t, "ACTIVE", polls.Data[0].Status)
	require.Equal(t, 1800, polls.Data[0].Duration)
	require.Equal(t, 2021, polls.Data[0].StartedAt.Year())

	choices := polls.Data[0].Choices
	require.Equal(t, "4c123012-1351-4f33-84b7-43856e7a0f47", choices[0].ID)
	require.Equal(t, "Heads", choices[0].Title)
	require.Equal(t, "Tails", choices[1].Title)
	require.Equal(t, 0, choices[0].Votes)
	require.Equal(t, 0, choices[0].ChannelPointsVotes)
	require.Equal(t, 0, choices[0].BitsVotes)
}

func TestClient_GetPollsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetPolls("", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPolls("    ", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPolls("	", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPolls("12345", nil, "", 21)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 20, but you supplied 21"})
	require.Nil(t, event)

	event, err = client.GetPolls("12345", nil, "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)

	event, err = client.GetPolls("12345", make([]string, 101), "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum ids that can be supplied is 100 but you supplied 101"})
	require.Nil(t, event)
}

func TestClient_GetPredictionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/predictions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"d6676d5c-c86e-44d2-bfc4-100fb48f0656\",\"broadcaster_id\":\"55696719\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"title\":\"Will there be any leaks today?\",\"winning_outcome_id\":null,\"outcomes\":[{\"id\":\"021e9234-5893-49b4-982e-cfe9a0aaddd9\",\"title\":\"Yes\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"BLUE\"},{\"id\":\"ded84c26-13cb-4b48-8cb5-5bae3ec3a66e\",\"title\":\"No\",\"users\":0,\"channel_points\":0,\"top_predictors\":null,\"color\":\"PINK\"}],\"prediction_window\":600,\"status\":\"ACTIVE\",\"created_at\":\"2021-04-28T16:03:06.320848689Z\",\"ended_at\":null,\"locked_at\":null}],\"pagination\":{}}")
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

	predictions, err := client.GetPredictions("12345", nil, "", 0)
	require.NoError(t, err)
	require.NotNil(t, predictions)
	require.Len(t, predictions.Data, 1)
	require.NotNil(t, predictions.Pagination)
	require.Equal(t, "d6676d5c-c86e-44d2-bfc4-100fb48f0656", predictions.Data[0].ID)
	require.Equal(t, "55696719", predictions.Data[0].BroadcasterID)
	require.Equal(t, "TwitchDev", predictions.Data[0].BroadcasterName)
	require.Equal(t, "twitchdev", predictions.Data[0].BroadcasterLogin)
	require.Equal(t, "Will there be any leaks today?", predictions.Data[0].Title)
	require.Zero(t, predictions.Data[0].WinningOutcomeId)
	require.Len(t, predictions.Data[0].Outcomes, 2)
	require.Equal(t, 600, predictions.Data[0].PredictionWindow)
	require.Equal(t, "ACTIVE", predictions.Data[0].Status)
	require.Equal(t, 2021, predictions.Data[0].CreatedAt.Year())
	require.Nil(t, predictions.Data[0].EndedAt)
	require.Nil(t, predictions.Data[0].LockedAt)
}

func TestClient_GetPredictionsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetPredictions("", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPredictions("    ", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPredictions("	", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetPredictions("12345", nil, "", 21)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 20, but you supplied 21"})
	require.Nil(t, event)

	event, err = client.GetPredictions("12345", nil, "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)

	event, err = client.GetPredictions("12345", make([]string, 101), "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum ids that can be supplied is 100 but you supplied 101"})
	require.Nil(t, event)
}

func TestClient_SearchCategoriesReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/search/categories/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "fort", r.URL.Query().Get("query"))
		require.Equal(t, "query=fort", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"33214\",\"name\":\"Fortnite\",\"box_art_url\":\"https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg\"}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN\"}}")
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

	results, err := client.SearchCategories("fort", "", 0)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 1)
	require.NotNil(t, results.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN", results.Pagination.Cursor)
	require.Equal(t, "33214", results.Data[0].ID)
	require.Equal(t, "Fortnite", results.Data[0].Name)
	require.Equal(t, "https://static-cdn.jtvnw.net/ttv-boxart/Fortnite-{width}x{height}.jpg", results.Data[0].BoxArtUrl)
}

func TestClient_SearchCategoriesReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.SearchCategories("", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchCategories("    ", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchCategories("	", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchCategories("12345", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.SearchCategories("12345", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)
}

func TestClient_SearchChannelsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/search/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "a_seagull", r.URL.Query().Get("query"))
		require.Equal(t, "query=a_seagull", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_language\":\"en\",\"broadcaster_login\":\"loserfruit\",\"display_name\":\"Loserfruit\",\"game_id\":\"498000\",\"game_name\":\"House Flipper\",\"id\":\"41245072\",\"is_live\":false,\"tag_ids\":[],\"thumbnail_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/fd17325a-7dc2-46c6-8617-e90ec259501c-profile_image-300x300.png\",\"title\":\"loserfruit\",\"started_at\":\"\"}],\"pagination\":{\"cursor\":\"Mg==\"}}")
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

	results, err := client.SearchChannels("a_seagull", "", 0, false)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 1)
	require.NotNil(t, results.Pagination)
	require.Equal(t, "Mg==", results.Pagination.Cursor)
	require.Equal(t, "en", results.Data[0].BroadcasterLanguage)
	require.Equal(t, "loserfruit", results.Data[0].BroadcasterLogin)
	require.Equal(t, "Loserfruit", results.Data[0].DisplayName)
	require.Equal(t, "498000", results.Data[0].GameID)
	require.Equal(t, "House Flipper", results.Data[0].GameName)
	require.Equal(t, "41245072", results.Data[0].ID)
	require.False(t, results.Data[0].IsLive)
	require.Len(t, results.Data[0].TagIDs, 0)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/fd17325a-7dc2-46c6-8617-e90ec259501c-profile_image-300x300.png", results.Data[0].ThumbnailUrl)
	require.Equal(t, "loserfruit", results.Data[0].Title)
	require.Equal(t, "", results.Data[0].StartedAt)
}

func TestClient_SearchChannelsReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.SearchChannels("", "", 0, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchChannels("    ", "", 0, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchChannels("	", "", 0, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, query can't be blank"})
	require.Nil(t, event)

	event, err = client.SearchChannels("12345", "", 101, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, event)

	event, err = client.SearchChannels("12345", "", -1, false)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, event)
}

func TestClient_GetStreamKeyReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/key/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"stream_key\":\"live_44322889_a34ub37c8ajv98a0\"}]}")
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

	results, err := client.GetStreamKey("12345")
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results.Data, 1)
	require.Equal(t, "live_44322889_a34ub37c8ajv98a0", results.Data[0].StreamKey)
}

func TestClient_GetStreamKeyReturnsError(t *testing.T) {
	client := Client{}
	event, err := client.GetStreamKey("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetStreamKey("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)

	event, err = client.GetStreamKey("		")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, event)
}

func TestClient_GetStreamsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"41375541868\",\"user_id\":\"459331509\",\"user_login\":\"auronplay\",\"user_name\":\"auronplay\",\"game_id\":\"494131\",\"game_name\":\"Little Nightmares\",\"type\":\"live\",\"title\":\"hablamos y le damos a Little Nightmares 1\",\"viewer_count\":78365,\"started_at\":\"2021-03-10T15:04:21Z\",\"language\":\"es\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/previews-ttv/live_user_auronplay-{width}x{height}.jpg\",\"tag_ids\":[\"d4bb9c58-2141-4881-bcdc-3fe0505457d1\"],\"is_mature\":false}],\"pagination\":{\"cursor\":\"eyJiIjp7IkN1cnNvciI6ImV5SnpJam8zT0RNMk5TNDBORFF4TlRjMU1UY3hOU3dpWkNJNlptRnNjMlVzSW5RaU9uUnlkV1Y5In0sImEiOnsiQ3Vyc29yIjoiZXlKeklqb3hOVGs0TkM0MU56RXhNekExTVRZNU1ESXNJbVFpT21aaGJITmxMQ0owSWpwMGNuVmxmUT09In19\"}}")
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

	streams, err := client.GetStreams("", "", 0, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, streams)
	require.NotNil(t, streams.Pagination)
	require.Equal(t, "eyJiIjp7IkN1cnNvciI6ImV5SnpJam8zT0RNMk5TNDBORFF4TlRjMU1UY3hOU3dpWkNJNlptRnNjMlVzSW5RaU9uUnlkV1Y5In0sImEiOnsiQ3Vyc29yIjoiZXlKeklqb3hOVGs0TkM0MU56RXhNekExTVRZNU1ESXNJbVFpT21aaGJITmxMQ0owSWpwMGNuVmxmUT09In19", streams.Pagination.Cursor)
	require.Len(t, streams.Data, 1)
	require.Equal(t, "41375541868", streams.Data[0].ID)
	require.Equal(t, "459331509", streams.Data[0].UserID)
	require.Equal(t, "auronplay", streams.Data[0].UserLogin)
	require.Equal(t, "auronplay", streams.Data[0].UserName)
	require.Equal(t, "494131", streams.Data[0].GameID)
	require.Equal(t, "Little Nightmares", streams.Data[0].GameName)
	require.Equal(t, "live", streams.Data[0].Type)
	require.Equal(t, "hablamos y le damos a Little Nightmares 1", streams.Data[0].Title)
	require.Equal(t, 78365, streams.Data[0].ViewerCount)
	require.Equal(t, 2021, streams.Data[0].StartedAt.Year())
	require.Equal(t, "es", streams.Data[0].Language)
	require.Equal(t, "https://static-cdn.jtvnw.net/previews-ttv/live_user_auronplay-{width}x{height}.jpg", streams.Data[0].ThumbnailUrl)
	require.Len(t, streams.Data[0].TagIDs, 1)
	require.Equal(t, "d4bb9c58-2141-4881-bcdc-3fe0505457d1", streams.Data[0].TagIDs[0])
	require.False(t, streams.Data[0].IsMature)
}

func TestClient_GetStreamsReturnsError(t *testing.T) {
	client := Client{}

	streams, err := client.GetStreams("", "", 101, nil, nil, nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, streams)

	streams, err = client.GetStreams("", "", -1, nil, nil, nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, streams)

	streams, err = client.GetStreams("", "", 0, make([]string, 101), nil, nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, max number of game ids is 100, but you supplied 101"})
	require.Nil(t, streams)

	streams, err = client.GetStreams("", "", 0, nil, make([]string, 101), nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, max number of languages is 100, but you supplied 101"})
	require.Nil(t, streams)

	streams, err = client.GetStreams("", "", 0, nil, nil, make([]string, 101), nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, max number of user ids is 100, but you supplied 101"})
	require.Nil(t, streams)

	streams, err = client.GetStreams("", "", 0, nil, nil, nil, make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, max number of user logins is 100, but you supplied 101"})
	require.Nil(t, streams)
}

func TestClient_GetFollowedStreamsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/followed/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("user_id"))
		require.Equal(t, "user_id=12345", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"42170724654\",\"user_id\":\"132954738\",\"user_login\":\"aws\",\"user_name\":\"AWS\",\"game_id\":\"417752\",\"game_name\":\"Talk Shows & Podcasts\",\"type\":\"live\",\"title\":\"AWS Howdy Partner! Y'all welcome ExtraHop to the show!\",\"viewer_count\":20,\"started_at\":\"2021-03-31T20:57:26Z\",\"language\":\"en\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/previews-ttv/live_user_aws-{width}x{height}.jpg\",\"tag_ids\":[\"6ea6bca4-4712-4ab9-a906-e3336a9d8039\"]}],\"pagination\":{\"cursor\":\"eyJiIjp7IkN1cnNvciI6ImV5SnpJam8zT0RNMk5TNDBORFF4TlRjMU1UY3hOU3dpWkNJNlptRnNjMlVzSW5RaU9uUnlkV1Y5In0sImEiOnsiQ3Vyc29yIjoiZXlKeklqb3hOVGs0TkM0MU56RXhNekExTVRZNU1ESXNJbVFpT21aaGJITmxMQ0owSWpwMGNuVmxmUT09In19\"}}")
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

	streams, err := client.GetFollowedStreams("12345", "", 0)
	require.NoError(t, err)
	require.NotNil(t, streams)
	require.NotNil(t, streams.Pagination)
	require.Equal(t, "eyJiIjp7IkN1cnNvciI6ImV5SnpJam8zT0RNMk5TNDBORFF4TlRjMU1UY3hOU3dpWkNJNlptRnNjMlVzSW5RaU9uUnlkV1Y5In0sImEiOnsiQ3Vyc29yIjoiZXlKeklqb3hOVGs0TkM0MU56RXhNekExTVRZNU1ESXNJbVFpT21aaGJITmxMQ0owSWpwMGNuVmxmUT09In19", streams.Pagination.Cursor)
	require.Equal(t, "42170724654", streams.Data[0].ID)
	require.Equal(t, "132954738", streams.Data[0].UserID)
	require.Equal(t, "aws", streams.Data[0].UserLogin)
	require.Equal(t, "AWS", streams.Data[0].UserName)
	require.Equal(t, "417752", streams.Data[0].GameID)
	require.Equal(t, "Talk Shows & Podcasts", streams.Data[0].GameName)
	require.Equal(t, "live", streams.Data[0].Type)
	require.Equal(t, "AWS Howdy Partner! Y'all welcome ExtraHop to the show!", streams.Data[0].Title)
	require.Equal(t, 20, streams.Data[0].ViewerCount)
	require.Equal(t, 2021, streams.Data[0].StartedAt.Year())
	require.Equal(t, "en", streams.Data[0].Language)
	require.Equal(t, "https://static-cdn.jtvnw.net/previews-ttv/live_user_aws-{width}x{height}.jpg", streams.Data[0].ThumbnailUrl)
	require.Len(t, streams.Data[0].TagIDs, 1)
	require.Equal(t, "6ea6bca4-4712-4ab9-a906-e3336a9d8039", streams.Data[0].TagIDs[0])
}

func TestClient_GetFollowedStreamsReturnsError(t *testing.T) {
	client := Client{}

	streams, err := client.GetFollowedStreams("", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, streams)

	streams, err = client.GetFollowedStreams("    ", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, streams)

	streams, err = client.GetFollowedStreams("	", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, streams)

	streams, err = client.GetFollowedStreams("12345", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, streams)

	streams, err = client.GetFollowedStreams("12345", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, streams)
}

func TestClient_GetStreamMarkersReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/markers/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("user_id"))
		require.Equal(t, "user_id=123", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"user_id\":\"123\",\"user_name\":\"TwitchName\",\"user_login\":\"twitchname\",\"videos\":[{\"video_id\":\"456\",\"markers\":[{\"id\":\"106b8d6243a4f883d25ad75e6cdffdc4\",\"created_at\":\"2018-08-20T20:10:03Z\",\"description\":\"hello, this is a marker!\",\"position_seconds\":244,\"URL\":\"https://twitch.tv/videos/456?t=0h4m06s\"}]}]}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjoiMjk1MjA0Mzk3OjI1Mzpib29rbWFyazoxMDZiOGQ1Y\"}}")
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

	markers, err := client.GetStreamMarkers("123", "", "", "", 0)
	require.NoError(t, err)
	require.NotNil(t, markers)
	require.NotNil(t, markers.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjoiMjk1MjA0Mzk3OjI1Mzpib29rbWFyazoxMDZiOGQ1Y", markers.Pagination.Cursor)
	require.Len(t, markers.Data, 1)
	require.Equal(t, "123", markers.Data[0].UserID)
	require.Equal(t, "TwitchName", markers.Data[0].UserName)
	require.Equal(t, "twitchname", markers.Data[0].UserLogin)
	require.Len(t, markers.Data[0].Videos, 1)

	video := markers.Data[0].Videos[0]
	require.Equal(t, "456", video.VideoId)
	require.Len(t, video.Markers, 1)
	require.Equal(t, "106b8d6243a4f883d25ad75e6cdffdc4", video.Markers[0].ID)
	require.Equal(t, 2018, video.Markers[0].CreatedAt.Year())
	require.Equal(t, "hello, this is a marker!", video.Markers[0].Description)
	require.Equal(t, 244, video.Markers[0].PositionSeconds)
	require.Equal(t, "https://twitch.tv/videos/456?t=0h4m06s", video.Markers[0].URL)
}

func TestClient_GetStreamMarkersReturnsError(t *testing.T) {
	client := Client{}

	markers, err := client.GetStreamMarkers("", "", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id or video id must be set"})
	require.Nil(t, markers)

	markers, err = client.GetStreamMarkers(" ", " ", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id or video id must be set"})
	require.Nil(t, markers)

	markers, err = client.GetStreamMarkers("		", "	", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id or video id must be set"})
	require.Nil(t, markers)

	markers, err = client.GetStreamMarkers("123", "456", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, only one of userID or videoID can be set"})
	require.Nil(t, markers)

	markers, err = client.GetStreamMarkers("123", "", "", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, markers)

	markers, err = client.GetStreamMarkers("123", "", "", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, markers)
}

func TestClient_GetBroadcasterSubscriptionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/subscriptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "broadcaster_id=123", r.URL.Query().Encode())
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_id\":\"141981764\",\"broadcaster_login\":\"twitchdev\",\"broadcaster_name\":\"TwitchDev\",\"gifter_id\":\"12826\",\"gifter_login\":\"twitch\",\"gifter_name\":\"Twitch\",\"is_gift\":true,\"tier\":\"1000\",\"plan_name\":\"Channel Subscription (twitchdev)\",\"user_id\":\"527115020\",\"user_name\":\"twitchgaming\",\"user_login\":\"twitchgaming\"}],\"pagination\":{\"cursor\":\"xxxx\"},\"total\":13}")
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

	subs, err := client.GetBroadcasterSubscriptions("123", nil, "", 0)
	require.NoError(t, err)
	require.NotNil(t, subs)
	require.NotNil(t, subs.Pagination)
	require.Equal(t, "xxxx", subs.Pagination.Cursor)
	require.Len(t, subs.Data, 1)
	require.Equal(t, 13, subs.Total)
	require.Equal(t, "141981764", subs.Data[0].BroadcasterID)
	require.Equal(t, "twitchdev", subs.Data[0].BroadcasterLogin)
	require.Equal(t, "TwitchDev", subs.Data[0].BroadcasterName)
	require.Equal(t, "12826", subs.Data[0].GifterID)
	require.Equal(t, "twitch", subs.Data[0].GifterLogin)
	require.Equal(t, "Twitch", subs.Data[0].GifterName)
	require.True(t, subs.Data[0].IsGift)
	require.Equal(t, "1000", subs.Data[0].Tier)
	require.Equal(t, "Channel Subscription (twitchdev)", subs.Data[0].PlanName)
	require.Equal(t, "527115020", subs.Data[0].UserID)
	require.Equal(t, "twitchgaming", subs.Data[0].UserName)
	require.Equal(t, "twitchgaming", subs.Data[0].UserLogin)
}

func TestClient_GetBroadcasterSubscriptionsReturnsError(t *testing.T) {
	client := Client{}

	subs, err := client.GetBroadcasterSubscriptions("", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, subs)

	subs, err = client.GetBroadcasterSubscriptions("    ", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, subs)

	subs, err = client.GetBroadcasterSubscriptions("	", nil, "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, subs)

	subs, err = client.GetBroadcasterSubscriptions("123", nil, "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, subs)

	subs, err = client.GetBroadcasterSubscriptions("123", nil, "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, subs)
}

func TestClient_CheckUserSubscriptionReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/subscriptions/user/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "149747285", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "141981764", r.URL.Query().Get("user_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_id\":\"149747285\",\"broadcaster_name\":\"TwitchPresents\",\"broadcaster_login\":\"twitchpresents\",\"is_gift\":false,\"tier\":\"1000\"}]}")
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

	subs, err := client.CheckUserSubscription("149747285", "141981764")
	require.NoError(t, err)
	require.NotNil(t, subs)
	require.Len(t, subs.Data, 1)
	require.Equal(t, "149747285", subs.Data[0].BroadcasterId)
	require.Equal(t, "TwitchPresents", subs.Data[0].BroadcasterName)
	require.Equal(t, "twitchpresents", subs.Data[0].BroadcasterLogin)
	require.Equal(t, "1000", subs.Data[0].Tier)
	require.False(t, subs.Data[0].IsGift)
	require.Zero(t, subs.Data[0].GifterID)
	require.Zero(t, subs.Data[0].GifterLogin)
	require.Zero(t, subs.Data[0].GifterName)

}

func TestClient_CheckUserSubscriptionReturnsError(t *testing.T) {
	client := Client{}

	subs, err := client.CheckUserSubscription("", "123")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, subs)

	subs, err = client.CheckUserSubscription("123", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user can't be blank"})
	require.Nil(t, subs)
}

func TestClient_GetAllStreamTagsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/tags/streams/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"tag_id\":\"621fb5bf-5498-4d8f-b4ac-db4d40d401bf\",\"is_auto\":false,\"localization_names\":{\"bg-bg\":\"Изчистване на 1 кредит\",\"cs-cz\":\"1 čistý kredit\",\"da-dk\":\"1 credit klaret\",\"de-de\":\"Mit 1 Leben abschließen\",\"el-gr\":\"1 μόνο πίστωση\",\"en-us\":\"1 Credit Clear\"},\"localization_descriptions\":{\"bg-bg\":\"За потоци с акцент върху завършване на аркадна игра с монети, в която не се използва продължаване\",\"cs-cz\":\"Pro vysílání s důrazem na plnění mincových arkádových her bez použití pokračování.\",\"da-dk\":\"Til streams med vægt på at gennemføre et arkadespil uden at bruge continues\",\"de-de\":\"Für Streams mit dem Ziel, ein Coin-op-Arcade-Game mit nur einem Leben abzuschließen.\",\"el-gr\":\"Για μεταδόσεις με έμφαση στην ολοκλήρωση παλαιού τύπου ηλεκτρονικών παιχνιδιών που λειτουργούν με κέρμα, χωρίς να χρησιμοποιούν συνέχειες\",\"en-us\":\"For streams with an emphasis on completing a coin-op arcade game without using any continues\"}}],\"pagination\":{\"cursor\":\"eyJiI...\"}}")
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

	tags, err := client.GetAllStreamTags("", 0, nil)
	require.NoError(t, err)
	require.NotNil(t, tags)
	require.NotNil(t, tags.Pagination)
	require.Equal(t, "eyJiI...", tags.Pagination.Cursor)
	require.Len(t, tags.Data, 1)
	require.Equal(t, "621fb5bf-5498-4d8f-b4ac-db4d40d401bf", tags.Data[0].TagID)
	require.False(t, tags.Data[0].IsAuto)
	require.NotNil(t, tags.Data[0].LocalizationNames)
	require.Len(t, tags.Data[0].LocalizationNames, 6)
	require.NotNil(t, tags.Data[0].LocalizationDescriptions)
	require.Len(t, tags.Data[0].LocalizationDescriptions, 6)

	require.Equal(t, "1 Credit Clear", tags.Data[0].LocalizationNames["en-us"])
	require.Equal(t, "For streams with an emphasis on completing a coin-op arcade game without using any continues", tags.Data[0].LocalizationDescriptions["en-us"])
}

func TestClient_GetAllStreamTagsReturnsError(t *testing.T) {
	client := Client{}

	tags, err := client.GetAllStreamTags("", 101, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, tags)

	tags, err = client.GetAllStreamTags("", -1, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, tags)

	tags, err = client.GetAllStreamTags("", 0, make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, max number of tag ids is 100, but you supplied 101"})
	require.Nil(t, tags)
}

func TestClient_GetStreamTagsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/tags/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "12345", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"tag_id\":\"6ea6bca4-4712-4ab9-a906-e3336a9d8039\",\"is_auto\":true,\"localization_names\":{\"bg-bg\":\"английски\",\"cs-cz\":\"Angličtina\",\"da-dk\":\"Engelsk\",\"de-de\":\"Englisch\",\"el-gr\":\"Αγγλικά\",\"en-us\":\"English\"},\"localization_descriptions\":{\"bg-bg\":\"За потоци с използване на английски\",\"cs-cz\":\"Pro vysílání obsahující angličtinu.\",\"da-dk\":\"Til streams, hvori der indgår engelsk\",\"de-de\":\"Für Streams auf Englisch.\",\"el-gr\":\"Για μεταδόσεις που περιλαμβάνουν τη χρήση Αγγλικών\",\"en-us\":\"For streams featuring the use of English\"}}]}")
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

	tags, err := client.GetStreamTags("12345")
	require.NoError(t, err)
	require.NotNil(t, tags)
	require.Nil(t, tags.Pagination)
	require.True(t, tags.Data[0].IsAuto)
	require.Equal(t, "6ea6bca4-4712-4ab9-a906-e3336a9d8039", tags.Data[0].TagID)
	require.NotNil(t, tags.Data[0].LocalizationNames)
	require.Len(t, tags.Data[0].LocalizationNames, 6)
	require.NotNil(t, tags.Data[0].LocalizationDescriptions)
	require.Len(t, tags.Data[0].LocalizationDescriptions, 6)
	require.Equal(t, "English", tags.Data[0].LocalizationNames["en-us"])
	require.Equal(t, "For streams featuring the use of English", tags.Data[0].LocalizationDescriptions["en-us"])
}

func TestClient_GetStreamTagsReturnsError(t *testing.T) {
	client := Client{}

	tags, err := client.GetStreamTags("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, tags)

	tags, err = client.GetStreamTags("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, tags)

	tags, err = client.GetStreamTags("		")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, tags)
}

func TestClient_GetChannelTeamsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/teams/channel/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "96909659", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"broadcaster_id\":\"96909659\",\"broadcaster_name\":\"CSharpFritz\",\"broadcaster_login\":\"csharpfritz\",\"background_image_url\":null,\"banner\":null,\"created_at\":\"2019-02-11T12:09:22Z\",\"updated_at\":\"2020-11-18T15:56:41Z\",\"info\":\"<p>An outgoing and enthusiastic group of friendly channels that write code, teach about technology, and promote the technical community.</p>\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/team-livecoders-team_logo_image-bf1d9a87ca81432687de60e24ad9593d-600x600.png\",\"team_name\":\"livecoders\",\"team_display_name\":\"Live Coders\",\"id\":\"6358\"}]}")
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

	teams, err := client.GetChannelTeams("96909659")
	require.NoError(t, err)
	require.NotNil(t, teams)
	require.Equal(t, "96909659", teams.Data[0].BroadcasterID)
	require.Equal(t, "CSharpFritz", teams.Data[0].BroadcasterName)
	require.Equal(t, "csharpfritz", teams.Data[0].BroadcasterLogin)
	require.Zero(t, teams.Data[0].BackgroundImageUrl)
	require.Zero(t, teams.Data[0].Banner)
	require.Equal(t, 2019, teams.Data[0].CreatedAt.Year())
	require.Equal(t, 2020, teams.Data[0].UpdatedAt.Year())
	require.Equal(t, "<p>An outgoing and enthusiastic group of friendly channels that write code, teach about technology, and promote the technical community.</p>", teams.Data[0].Info)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/team-livecoders-team_logo_image-bf1d9a87ca81432687de60e24ad9593d-600x600.png", teams.Data[0].ThumbnailUrl)
	require.Equal(t, "livecoders", teams.Data[0].TeamName)
	require.Equal(t, "Live Coders", teams.Data[0].TeamDisplayName)
	require.Equal(t, "6358", teams.Data[0].ID)
}

func TestClient_GetChannelTeamsReturnsError(t *testing.T) {
	client := Client{}

	teams, err := client.GetChannelTeams("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, teams)

	teams, err = client.GetChannelTeams("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, teams)

	teams, err = client.GetChannelTeams("		")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, teams)
}

func TestClient_GetTeamReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/teams/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "6358", r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"users\":[{\"user_id\":\"278217731\",\"user_name\":\"mastermndio\",\"user_login\":\"mastermndio\"},{\"user_id\":\"41284990\",\"user_name\":\"jenninexus\",\"user_login\":\"jenninexus\"}],\"background_image_url\":null,\"banner\":null,\"created_at\":\"2019-02-11T12:09:22Z\",\"updated_at\":\"2020-11-18T15:56:41Z\",\"info\":\"<p>An outgoing and enthusiastic group of friendly channels that write code, teach about technology, and promote the technical community.</p>\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/team-livecoders-team_logo_image-bf1d9a87ca81432687de60e24ad9593d-600x600.png\",\"team_name\":\"livecoders\",\"team_display_name\":\"Live Coders\",\"id\":\"6358\"}]}")
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

	teams, err := client.GetTeam("", "6358")
	require.NoError(t, err)
	require.NotNil(t, teams)
	require.Len(t, teams.Data, 1)
	require.Len(t, teams.Data[0].Users, 2)
	require.Equal(t, "278217731", teams.Data[0].Users[0].UserID)
	require.Equal(t, "mastermndio", teams.Data[0].Users[0].UserName)
	require.Equal(t, "mastermndio", teams.Data[0].Users[0].UserLogin)
	require.Zero(t, teams.Data[0].BackgroundImageUrl)
	require.Zero(t, teams.Data[0].Banner)
	require.Equal(t, 2019, teams.Data[0].CreatedAt.Year())
	require.Equal(t, 2020, teams.Data[0].UpdatedAt.Year())
	require.Equal(t, "<p>An outgoing and enthusiastic group of friendly channels that write code, teach about technology, and promote the technical community.</p>", teams.Data[0].Info)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/team-livecoders-team_logo_image-bf1d9a87ca81432687de60e24ad9593d-600x600.png", teams.Data[0].ThumbnailUrl)
	require.Equal(t, "livecoders", teams.Data[0].TeamName)
	require.Equal(t, "Live Coders", teams.Data[0].TeamDisplayName)
	require.Equal(t, "6358", teams.Data[0].ID)
}

func TestClient_GetTeamReturnsError(t *testing.T) {
	client := Client{}

	teams, err := client.GetTeam("", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, name or id must be specified"})
	require.Nil(t, teams)

	teams, err = client.GetTeam("  ", "  ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, name or id must be specified"})
	require.Nil(t, teams)

	teams, err = client.GetTeam("		", "	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, name or id must be specified"})
	require.Nil(t, teams)

	teams, err = client.GetTeam("foo", "bar")
	require.ErrorIs(t, err, BadRequestError{"invalid request, only one of name or id may be specified"})
	require.Nil(t, teams)
}

func TestClient_GetUsersReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"141981764\",\"login\":\"twitchdev\",\"display_name\":\"TwitchDev\",\"type\":\"\",\"broadcaster_type\":\"partner\",\"description\":\"Supporting third-party developers building Twitch integrations from chatbots to game integrations.\",\"profile_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png\",\"offline_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png\",\"view_count\":5980557,\"email\":\"not-real@email.com\",\"created_at\":\"2016-12-14T20:32:28.894263Z\"}]}")
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

	users, err := client.GetUsers([]string{"141981764"}, nil)
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Len(t, users.Data, 1)
	require.Equal(t, "141981764", users.Data[0].ID)
	require.Equal(t, "twitchdev", users.Data[0].Login)
	require.Equal(t, "TwitchDev", users.Data[0].DisplayName)
	require.Equal(t, "", users.Data[0].Type)
	require.Equal(t, "partner", users.Data[0].BroadcasterType)
	require.Equal(t, "Supporting third-party developers building Twitch integrations from chatbots to game integrations.", users.Data[0].Description)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/8a6381c7-d0c0-4576-b179-38bd5ce1d6af-profile_image-300x300.png", users.Data[0].ProfileImageUrl)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/3f13ab61-ec78-4fe6-8481-8682cb3b0ac2-channel_offline_image-1920x1080.png", users.Data[0].OfflineImageUrl)
	require.Equal(t, 5980557, users.Data[0].ViewCount)
	require.Equal(t, "not-real@email.com", users.Data[0].Email)
	require.Equal(t, 2016, users.Data[0].CreatedAt.Year())
}

func TestClient_GetUsersReturnsError(t *testing.T) {
	client := Client{}

	users, err := client.GetUsers(nil, nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, login or id must be specified"})
	require.Nil(t, users)

	users, err = client.GetUsers(make([]string, 101), nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum of 100 ids and logins, but you supplied 101"})
	require.Nil(t, users)

	users, err = client.GetUsers(nil, make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum of 100 ids and logins, but you supplied 101"})
	require.Nil(t, users)

	users, err = client.GetUsers(make([]string, 101), make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum of 100 ids and logins, but you supplied 202"})
	require.Nil(t, users)
}

func TestClient_GetUsersFollowsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/follows/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "1234", r.URL.Query().Get("from_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"total\":12345,\"data\":[{\"from_id\":\"171003792\",\"from_login\":\"iiisutha067iii\",\"from_name\":\"IIIsutha067III\",\"to_id\":\"23161357\",\"to_name\":\"LIRIK\",\"followed_at\":\"2017-08-22T22:55:24Z\"},{\"from_id\":\"113627897\",\"from_login\":\"birdman616\",\"from_name\":\"Birdman616\",\"to_id\":\"23161357\",\"to_name\":\"LIRIK\",\"followed_at\":\"2017-08-22T22:55:04Z\"}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjoiMTUwMzQ0MTc3NjQyNDQyMjAwMCJ9\"}}")
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

	followers, err := client.GetUsersFollows("1234", "", "", 0)
	require.NoError(t, err)
	require.NotNil(t, followers)
	require.NotNil(t, followers.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjoiMTUwMzQ0MTc3NjQyNDQyMjAwMCJ9", followers.Pagination.Cursor)
	require.Equal(t, 12345, followers.Total)
	require.Len(t, followers.Data, 2)
	require.Equal(t, "171003792", followers.Data[0].FromId)
	require.Equal(t, "iiisutha067iii", followers.Data[0].FromLogin)
	require.Equal(t, "IIIsutha067III", followers.Data[0].FromName)
	require.Equal(t, "23161357", followers.Data[0].ToId)
	require.Equal(t, "LIRIK", followers.Data[0].ToName)
	require.Equal(t, 2017, followers.Data[0].FollowedAt.Year())
}

func TestClient_GetUsersFollowsReturnsError(t *testing.T) {
	client := Client{}

	followers, err := client.GetUsersFollows("", "", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, either fromID or toID must be specified"})
	require.Nil(t, followers)

	followers, err = client.GetUsersFollows("foo", "bar", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, fromID or toID must be specified, but not both"})
	require.Nil(t, followers)

	followers, err = client.GetUsersFollows("1234", "", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, followers)

	followers, err = client.GetUsersFollows("1234", "", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, followers)
}

func TestClient_GetUsersBlockListReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/blocks/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"user_id\":\"135093069\",\"user_login\":\"bluelava\",\"display_name\":\"BlueLava\"},{\"user_id\":\"27419011\",\"user_login\":\"travistyoj\",\"display_name\":\"TravistyOJ\"}]}")
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

	blockList, err := client.GetUsersBlockList("141981764", "", 0)
	require.NoError(t, err)
	require.NotNil(t, blockList)
	require.Nil(t, blockList.Pagination)
	require.Len(t, blockList.Data, 2)
	require.Equal(t, "135093069", blockList.Data[0].UserID)
	require.Equal(t, "bluelava", blockList.Data[0].UserLogin)
	require.Equal(t, "BlueLava", blockList.Data[0].DisplayName)
}

func TestClient_GetUsersBlockListReturnsError(t *testing.T) {
	client := Client{}

	blocklist, err := client.GetUsersBlockList("", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, blocklist)

	blocklist, err = client.GetUsersBlockList("    ", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, blocklist)

	blocklist, err = client.GetUsersBlockList("	", "", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, blocklist)

	blocklist, err = client.GetUsersBlockList("1234", "", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, blocklist)

	blocklist, err = client.GetUsersBlockList("1234", "", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, blocklist)
}

func TestClient_GetUserExtensionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/extensions/list/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"wi08ebtatdc7oj83wtl9uxwz807l8b\",\"version\":\"1.1.8\",\"name\":\"Streamlabs Leaderboard\",\"can_activate\":true,\"type\":[\"panel\"]},{\"id\":\"d4uvtfdr04uq6raoenvj7m86gdk16v\",\"version\":\"2.0.2\",\"name\":\"Prime Subscription and Loot Reminder\",\"can_activate\":true,\"type\":[\"overlay\"]},{\"id\":\"rh6jq1q334hqc2rr1qlzqbvwlfl3x0\",\"version\":\"1.1.0\",\"name\":\"TopClip\",\"can_activate\":true,\"type\":[\"mobile\",\"panel\"]},{\"id\":\"zfh2irvx2jb4s60f02jq0ajm8vwgka\",\"version\":\"1.0.19\",\"name\":\"Streamlabs\",\"can_activate\":true,\"type\":[\"mobile\",\"overlay\"]},{\"id\":\"lqnf3zxk0rv0g7gq92mtmnirjz2cjj\",\"version\":\"0.0.1\",\"name\":\"Dev Experience Test\",\"can_activate\":true,\"type\":[\"component\",\"mobile\",\"panel\",\"overlay\"]}]}")
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

	extensions, err := client.GetUserExtensions()
	require.NoError(t, err)
	require.NotNil(t, extensions)
	require.Len(t, extensions.Data, 5)
	require.Equal(t, "wi08ebtatdc7oj83wtl9uxwz807l8b", extensions.Data[0].ID)
	require.Equal(t, "1.1.8", extensions.Data[0].Version)
	require.Equal(t, "Streamlabs Leaderboard", extensions.Data[0].Name)
	require.True(t, extensions.Data[0].CanActivate)
	require.Len(t, extensions.Data[0].Type, 1)
	require.Equal(t, "panel", extensions.Data[0].Type[0])
}

func TestClient_GetUserActiveExtensionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/extensions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":{\"panel\":{\"1\":{\"active\":true,\"id\":\"rh6jq1q334hqc2rr1qlzqbvwlfl3x0\",\"version\":\"1.1.0\",\"name\":\"TopClip\"},\"2\":{\"active\":true,\"id\":\"wi08ebtatdc7oj83wtl9uxwz807l8b\",\"version\":\"1.1.8\",\"name\":\"Streamlabs Leaderboard\"},\"3\":{\"active\":true,\"id\":\"naty2zwfp7vecaivuve8ef1hohh6bo\",\"version\":\"1.0.9\",\"name\":\"Streamlabs Stream Schedule & Countdown\"}},\"overlay\":{\"1\":{\"active\":true,\"id\":\"zfh2irvx2jb4s60f02jq0ajm8vwgka\",\"version\":\"1.0.19\",\"name\":\"Streamlabs\"}},\"component\":{\"1\":{\"active\":true,\"id\":\"lqnf3zxk0rv0g7gq92mtmnirjz2cjj\",\"version\":\"0.0.1\",\"name\":\"Dev Experience Test\",\"x\":0,\"y\":0},\"2\":{\"active\":false}}}}")
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

	extensions, err := client.GetUserActiveExtensions("")
	require.NoError(t, err)
	require.NotNil(t, extensions)
	require.NotNil(t, extensions.Data)
	require.NotNil(t, extensions.Data.Component)
	require.NotNil(t, extensions.Data.Panel)
	require.NotNil(t, extensions.Data.Overlay)

	require.True(t, extensions.Data.Panel.Field1.Active)
	require.Equal(t, "rh6jq1q334hqc2rr1qlzqbvwlfl3x0", extensions.Data.Panel.Field1.ID)
	require.Equal(t, "1.1.0", extensions.Data.Panel.Field1.Version)
	require.Equal(t, "TopClip", extensions.Data.Panel.Field1.Name)

	require.True(t, extensions.Data.Overlay.Field1.Active)
	require.Equal(t, "zfh2irvx2jb4s60f02jq0ajm8vwgka", extensions.Data.Overlay.Field1.ID)
	require.Equal(t, "1.0.19", extensions.Data.Overlay.Field1.Version)
	require.Equal(t, "Streamlabs", extensions.Data.Overlay.Field1.Name)

	require.True(t, extensions.Data.Component.Field1.Active)
	require.Equal(t, "lqnf3zxk0rv0g7gq92mtmnirjz2cjj", extensions.Data.Component.Field1.ID)
	require.Equal(t, "0.0.1", extensions.Data.Component.Field1.Version)
	require.Equal(t, "Dev Experience Test", extensions.Data.Component.Field1.Name)
	require.Equal(t, 0, extensions.Data.Component.Field1.X)
	require.Equal(t, 0, extensions.Data.Component.Field1.Y)

	require.False(t, extensions.Data.Component.Field2.Active)
}

func TestClient_GetVideosByIDReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/videos/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "335921245", r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"335921245\",\"stream_id\":null,\"user_id\":\"141981764\",\"user_login\":\"twitchdev\",\"user_name\":\"TwitchDev\",\"title\":\"Twitch Developers 101\",\"description\":\"Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.\",\"created_at\":\"2018-11-14T21:30:18Z\",\"published_at\":\"2018-11-14T22:04:30Z\",\"url\":\"https://www.twitch.tv/videos/335921245\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg\",\"viewable\":\"public\",\"view_count\":1863062,\"language\":\"en\",\"type\":\"upload\",\"duration\":\"3m21s\",\"muted_segments\":[{\"duration\":30,\"offset\":120}]}],\"pagination\":{}}")
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

	videos, err := client.GetVideosByID([]string{"335921245"})
	require.NoError(t, err)
	require.NotNil(t, videos)
	require.NotNil(t, videos.Pagination)
	require.Len(t, videos.Data, 1)
	require.Equal(t, "335921245", videos.Data[0].ID)
	require.Zero(t, videos.Data[0].StreamID)
	require.Equal(t, "141981764", videos.Data[0].UserID)
	require.Equal(t, "twitchdev", videos.Data[0].UserLogin)
	require.Equal(t, "TwitchDev", videos.Data[0].UserName)
	require.Equal(t, "Twitch Developers 101", videos.Data[0].Title)
	require.Equal(t, "Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.", videos.Data[0].Description)
	require.Equal(t, 2018, videos.Data[0].CreatedAt.Year())
	require.Equal(t, 2018, videos.Data[0].PublishedAt.Year())
	require.Equal(t, "https://www.twitch.tv/videos/335921245", videos.Data[0].Url)
	require.Equal(t, "https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg", videos.Data[0].ThumbnailUrl)
	require.Equal(t, "public", videos.Data[0].Viewable)
	require.Equal(t, 1863062, videos.Data[0].ViewCount)
	require.Equal(t, "en", videos.Data[0].Language)
	require.Equal(t, "upload", videos.Data[0].Type)
	require.Equal(t, "3m21s", videos.Data[0].Duration)
	require.Len(t, videos.Data[0].MutedSegments, 1)
	require.Equal(t, 30, videos.Data[0].MutedSegments[0].Duration)
	require.Equal(t, 120, videos.Data[0].MutedSegments[0].Offset)
}

func TestClient_GetVideosByIDReturnsError(t *testing.T) {
	client := Client{}

	videos, err := client.GetVideosByID(nil)
	require.ErrorIs(t, err, BadRequestError{"invalid request, at least one video id is required"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByID(make([]string, 0))
	require.ErrorIs(t, err, BadRequestError{"invalid request, at least one video id is required"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByID(make([]string, 101))
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum of 100 ids, but you supplied 101"})
	require.Nil(t, videos)
}

func TestClient_GetVideosByUserReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/videos/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("user_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"335921245\",\"stream_id\":null,\"user_id\":\"141981764\",\"user_login\":\"twitchdev\",\"user_name\":\"TwitchDev\",\"title\":\"Twitch Developers 101\",\"description\":\"Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.\",\"created_at\":\"2018-11-14T21:30:18Z\",\"published_at\":\"2018-11-14T22:04:30Z\",\"url\":\"https://www.twitch.tv/videos/335921245\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg\",\"viewable\":\"public\",\"view_count\":1863062,\"language\":\"en\",\"type\":\"upload\",\"duration\":\"3m21s\",\"muted_segments\":[{\"duration\":30,\"offset\":120}]}],\"pagination\":{}}")
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

	videos, err := client.GetVideosByUser("141981764", "", "", "", 0, "", "", "")
	require.NoError(t, err)
	require.NotNil(t, videos)
	require.NotNil(t, videos.Pagination)
	require.Len(t, videos.Data, 1)
	require.Equal(t, "335921245", videos.Data[0].ID)
	require.Zero(t, videos.Data[0].StreamID)
	require.Equal(t, "141981764", videos.Data[0].UserID)
	require.Equal(t, "twitchdev", videos.Data[0].UserLogin)
	require.Equal(t, "TwitchDev", videos.Data[0].UserName)
	require.Equal(t, "Twitch Developers 101", videos.Data[0].Title)
	require.Equal(t, "Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.", videos.Data[0].Description)
	require.Equal(t, 2018, videos.Data[0].CreatedAt.Year())
	require.Equal(t, 2018, videos.Data[0].PublishedAt.Year())
	require.Equal(t, "https://www.twitch.tv/videos/335921245", videos.Data[0].Url)
	require.Equal(t, "https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg", videos.Data[0].ThumbnailUrl)
	require.Equal(t, "public", videos.Data[0].Viewable)
	require.Equal(t, 1863062, videos.Data[0].ViewCount)
	require.Equal(t, "en", videos.Data[0].Language)
	require.Equal(t, "upload", videos.Data[0].Type)
	require.Equal(t, "3m21s", videos.Data[0].Duration)
	require.Len(t, videos.Data[0].MutedSegments, 1)
	require.Equal(t, 30, videos.Data[0].MutedSegments[0].Duration)
	require.Equal(t, 120, videos.Data[0].MutedSegments[0].Offset)
}

func TestClient_GetVideosByUserReturnsError(t *testing.T) {
	client := Client{}

	videos, err := client.GetVideosByUser("", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("    ", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("	", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("1234", "", "", "", 101, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("1234", "", "", "", -1, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("1234", "", "", "", 0, "decade", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, period can only all, day, week, month, or year, and you input decade"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("1234", "", "", "", 0, "", "foo", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, sort can only be all, time, trending, or views, but you input foo"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByUser("1234", "", "", "", 0, "", "", "foo")
	require.ErrorIs(t, err, BadRequestError{"invalid request, type can only be all, upload, archive, or highlight, but you input foo"})
	require.Nil(t, videos)
}

func TestClient_GetVideosByGameReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/videos/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("game_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"335921245\",\"stream_id\":null,\"user_id\":\"141981764\",\"user_login\":\"twitchdev\",\"user_name\":\"TwitchDev\",\"title\":\"Twitch Developers 101\",\"description\":\"Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.\",\"created_at\":\"2018-11-14T21:30:18Z\",\"published_at\":\"2018-11-14T22:04:30Z\",\"url\":\"https://www.twitch.tv/videos/335921245\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg\",\"viewable\":\"public\",\"view_count\":1863062,\"language\":\"en\",\"type\":\"upload\",\"duration\":\"3m21s\",\"muted_segments\":[{\"duration\":30,\"offset\":120}]}],\"pagination\":{}}")
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

	videos, err := client.GetVideosByGame("141981764", "", "", "", 0, "", "", "")
	require.NoError(t, err)
	require.NotNil(t, videos)
	require.NotNil(t, videos.Pagination)
	require.Len(t, videos.Data, 1)
	require.Equal(t, "335921245", videos.Data[0].ID)
	require.Zero(t, videos.Data[0].StreamID)
	require.Equal(t, "141981764", videos.Data[0].UserID)
	require.Equal(t, "twitchdev", videos.Data[0].UserLogin)
	require.Equal(t, "TwitchDev", videos.Data[0].UserName)
	require.Equal(t, "Twitch Developers 101", videos.Data[0].Title)
	require.Equal(t, "Welcome to Twitch development! Here is a quick overview of our products and information to help you get started.", videos.Data[0].Description)
	require.Equal(t, 2018, videos.Data[0].CreatedAt.Year())
	require.Equal(t, 2018, videos.Data[0].PublishedAt.Year())
	require.Equal(t, "https://www.twitch.tv/videos/335921245", videos.Data[0].Url)
	require.Equal(t, "https://static-cdn.jtvnw.net/cf_vods/d2nvs31859zcd8/twitchdev/335921245/ce0f3a7f-57a3-4152-bc06-0c6610189fb3/thumb/index-0000000000-%{width}x%{height}.jpg", videos.Data[0].ThumbnailUrl)
	require.Equal(t, "public", videos.Data[0].Viewable)
	require.Equal(t, 1863062, videos.Data[0].ViewCount)
	require.Equal(t, "en", videos.Data[0].Language)
	require.Equal(t, "upload", videos.Data[0].Type)
	require.Equal(t, "3m21s", videos.Data[0].Duration)
	require.Len(t, videos.Data[0].MutedSegments, 1)
	require.Equal(t, 30, videos.Data[0].MutedSegments[0].Duration)
	require.Equal(t, 120, videos.Data[0].MutedSegments[0].Offset)
}

func TestClient_GetVideosByGameReturnsError(t *testing.T) {
	client := Client{}

	videos, err := client.GetVideosByGame("", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("    ", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("	", "", "", "", 0, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, user id can't be blank"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("1234", "", "", "", 101, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("1234", "", "", "", -1, "", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("1234", "", "", "", 0, "decade", "", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, period can only all, day, week, month, or year, and you input decade"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("1234", "", "", "", 0, "", "foo", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, sort can only be all, time, trending, or views, but you input foo"})
	require.Nil(t, videos)

	videos, err = client.GetVideosByGame("1234", "", "", "", 0, "", "", "foo")
	require.ErrorIs(t, err, BadRequestError{"invalid request, type can only be all, upload, archive, or highlight, but you input foo"})
	require.Nil(t, videos)
}

func TestClient_GetWebhookSubscriptionsReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/webhooks/subscriptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"total\":12,\"data\":[{\"topic\":\"https://api.twitch.tv/helix/streams?user_id=123\",\"callback\":\"http://example.com/your_callback\",\"expires_at\":\"2018-07-30T20:00:00Z\"},{\"topic\":\"https://api.twitch.tv/helix/streams?user_id=345\",\"callback\":\"http://example.com/your_callback\",\"expires_at\":\"2018-07-30T20:03:00Z\"}],\"pagination\":{\"cursor\":\"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IkFYc2laU0k2TVN3aWFTSTZNWDAifX0\"}}")
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

	subs, err := client.GetWebhookSubscriptions("", 0)
	require.NoError(t, err)
	require.NotNil(t, subs)
	require.NotNil(t, subs.Pagination)
	require.Equal(t, "eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6IkFYc2laU0k2TVN3aWFTSTZNWDAifX0", subs.Pagination.Cursor)
	require.Equal(t, 12, subs.Total)
	require.Len(t, subs.Data, 2)
	require.Equal(t, "https://api.twitch.tv/helix/streams?user_id=123", subs.Data[0].Topic)
	require.Equal(t, "http://example.com/your_callback", subs.Data[0].Callback)
	require.Equal(t, 2018, subs.Data[0].ExpiresAt.Year())
}

func TestClient_GetWebhookSubscriptionsReturnsError(t *testing.T) {
	client := Client{}

	subs, err := client.GetWebhookSubscriptions("", 101)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 100, but you supplied 101"})
	require.Nil(t, subs)

	subs, err = client.GetWebhookSubscriptions("", -1)
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, subs)
}

func TestClient_GetChannelStreamScheduleReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":{\"segments\":[{\"id\":\"eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=\",\"start_time\":\"2021-07-01T18:00:00Z\",\"end_time\":\"2021-07-01T19:00:00Z\",\"title\":\"TwitchDev Monthly Update // July 1, 2021\",\"canceled_until\":null,\"category\":{\"id\":\"509670\",\"name\":\"Science & Technology\"},\"is_recurring\":false}],\"broadcaster_id\":\"141981764\",\"broadcaster_name\":\"TwitchDev\",\"broadcaster_login\":\"twitchdev\",\"vacation\":null},\"pagination\":{}}")
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

	schedule, err := client.GetChannelStreamSchedule("141981764", nil, nil, 0, 0, "")
	require.NoError(t, err)
	require.NotNil(t, schedule)
	require.NotNil(t, schedule.Pagination)
	require.Zero(t, schedule.Pagination.Cursor)
	require.Equal(t, "141981764", schedule.Data.BroadcasterId)
	require.Equal(t, "TwitchDev", schedule.Data.BroadcasterName)
	require.Equal(t, "twitchdev", schedule.Data.BroadcasterLogin)
	require.Nil(t, schedule.Data.Vacation)
	require.Len(t, schedule.Data.Segments, 1)
	require.Nil(t, schedule.Data.Segments[0].CanceledUntil)
	require.Equal(t, "eyJzZWdtZW50SUQiOiJlNGFjYzcyNC0zNzFmLTQwMmMtODFjYS0yM2FkYTc5NzU5ZDQiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyNn0=", schedule.Data.Segments[0].Id)
	require.Equal(t, 2021, schedule.Data.Segments[0].StartTime.Year())
	require.Equal(t, 2021, schedule.Data.Segments[0].EndTime.Year())
	require.Equal(t, "TwitchDev Monthly Update // July 1, 2021", schedule.Data.Segments[0].Title)
	require.Equal(t, "509670", schedule.Data.Segments[0].Category.Id)
	require.Equal(t, "Science & Technology", schedule.Data.Segments[0].Category.Name)
	require.False(t, schedule.Data.Segments[0].IsRecurring)
}

func TestClient_GetChannelStreamScheduleReturnsError(t *testing.T) {
	client := Client{}

	schedule, err := client.GetChannelStreamSchedule("", nil, nil, 0, 0, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamSchedule("    ", nil, nil, 0, 0, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamSchedule("		", nil, nil, 0, 0, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamSchedule("123", make([]string, 101), nil, 0, 0, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, maximum of 100 ids, but you supplied 101"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamSchedule("123", nil, nil, 0, 26, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count maximum value is 25, but you supplied 26"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamSchedule("123", nil, nil, 0, -1, "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, count can't be negative"})
	require.Nil(t, schedule)
}

func TestClient_GetChannelStreamScheduleAsICalReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/icalendar/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "BEGIN:VCALENDAR\nPRODID:-//twitch.tv//StreamSchedule//1.0\nVERSION:2.0\nCALSCALE:GREGORIAN\nREFRESH-INTERVAL;VALUE=DURATION:PT1H\nNAME:TwitchDev\nBEGIN:VEVENT\nUID:e4acc724-371f-402c-81ca-23ada79759d4\nDTSTAMP:20210323T040131Z\nDTSTART;TZID=/America/New_York:20210701T140000\nDTEND;TZID=/America/New_York:20210701T150000\nSUMMARY:TwitchDev Monthly Update // July 1, 2021\nDESCRIPTION:Science & Technology.\nCATEGORIES:Science & Technology\nEND:VEVENT\nEND:VCALENDAR%")
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

	schedule, err := client.GetChannelStreamScheduleAsICal("141981764")
	require.NoError(t, err)
	require.NotNil(t, schedule)
	require.Contains(t, string(schedule), "e4acc724-371f-402c-81ca-23ada79759d4")
}

func TestClient_GetChannelStreamScheduleAsICalReturnsError(t *testing.T) {
	client := Client{}

	schedule, err := client.GetChannelStreamScheduleAsICal("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamScheduleAsICal("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)

	schedule, err = client.GetChannelStreamScheduleAsICal("		")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcast can't be blank"})
	require.Nil(t, schedule)
}
