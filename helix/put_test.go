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

func TestClient_PutRequestReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PUT", r.Method)

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

	shouldBeNil, err := client.PutRequest("channels", nil, nil)
	require.ErrorIs(t, err, gotau.AuthorizationError{})
	require.Nil(t, shouldBeNil)
}

func TestClient_PutRequestReturnsGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PUT", r.Method)

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

	shouldBeNil, err := client.PutRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, gotau.GenericError{}, err)
	require.Nil(t, shouldBeNil)
}

func TestClient_PutRequestReturnsRateLimitError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channels/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PUT", r.Method)

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

	shouldBeNil, err := client.PutRequest("channels", nil, nil)
	require.Error(t, err)
	require.IsType(t, RateLimitError{}, err)
	rlErr := err.(RateLimitError)
	require.NotNil(t, rlErr)
	require.Equal(t, 2021, rlErr.ResetTime().Year())
	require.Equal(t, time.Month(6), rlErr.ResetTime().Month())
	require.Equal(t, 17, rlErr.ResetTime().Day())
	require.Nil(t, shouldBeNil)
}

func TestClient_ReplaceStreamTagsReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/streams/tags/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "257788195", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "PUT", r.Method)

		tags := make(map[string][]string)
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &tags)
		require.NoError(t, err)

		require.Len(t, tags["tag_ids"], 2)
		require.Equal(t, "621fb5bf-5498-4d8f-b4ac-db4d40d401bf", tags["tag_ids"][0])
		require.Equal(t, "79977fb9-f106-4a87-a386-f1b0f99783dd", tags["tag_ids"][1])

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

	replaced, err := client.ReplaceStreamTags("257788195", []string{"621fb5bf-5498-4d8f-b4ac-db4d40d401bf",
		"79977fb9-f106-4a87-a386-f1b0f99783dd"})
	require.NoError(t, err)
	require.True(t, replaced)
}

func TestClient_ReplaceStreamTagsReturnsError(t *testing.T) {
	client := Client{}

	replaced, err := client.ReplaceStreamTags("", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcast can't be blank"})
	require.False(t, replaced)

	replaced, err = client.ReplaceStreamTags("    ", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcast can't be blank"})
	require.False(t, replaced)

	replaced, err = client.ReplaceStreamTags("		", nil)
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcast can't be blank"})
	require.False(t, replaced)

	replaced, err = client.ReplaceStreamTags("123", make([]string, 6))
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, maximum of 5 tags can be set"})
	require.False(t, replaced)
}

func TestClient_UpdateUserReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "PUT", r.Method)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"data\":[{\"id\":\"44322889\",\"login\":\"dallas\",\"display_name\":\"dallas\",\"type\":\"staff\",\"broadcaster_type\":\"affiliate\",\"description\":\"BaldAngel\",\"profile_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/4d1f36cbf1f0072d-profile_image-300x300.png\",\"offline_image_url\":\"https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-channel_offline_image-2e82c1df2a464df7-1920x1080.jpeg\",\"view_count\":6995,\"email\":\"not-real@email.com\",\"created_at\":\"2013-06-03T19:12:02.580593Z\"}]}")
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

	description := "BaldAngel"

	user, err := client.UpdateUser(&description)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "44322889", user.Data[0].ID)
	require.Equal(t, "dallas", user.Data[0].Login)
	require.Equal(t, "dallas", user.Data[0].DisplayName)
	require.Equal(t, "staff", user.Data[0].Type)
	require.Equal(t, "affiliate", user.Data[0].BroadcasterType)
	require.Equal(t, description, user.Data[0].Description)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/4d1f36cbf1f0072d-profile_image-300x300.png", user.Data[0].ProfileImageUrl)
	require.Equal(t, "https://static-cdn.jtvnw.net/jtv_user_pictures/dallas-channel_offline_image-2e82c1df2a464df7-1920x1080.jpeg", user.Data[0].OfflineImageUrl)
	require.Equal(t, 6995, user.Data[0].ViewCount)
	require.Equal(t, "not-real@email.com", user.Data[0].Email)
	require.Equal(t, 2013, user.Data[0].CreatedAt.Year())
}
