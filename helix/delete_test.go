package helix

import (
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClient_DeleteRequestReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
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

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteRequestReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
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

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.Error(t, err)
	require.IsType(t, AuthorizationError{}, err)
	require.False(t, deleted)
}

func TestClient_DeleteRequestReturnsGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
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

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.Error(t, err)
	genericError, ok := err.(GenericError)
	require.True(t, ok)
	require.False(t, deleted)

	require.Equal(t, 404, genericError.StatusCode())
	require.Zero(t, genericError.Body())
}

func TestClient_DeleteCustomRewardReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "456", r.URL.Query().Get("id"))
		require.Equal(t, "DELETE", r.Method)
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

	deleted, err := client.DeleteCustomReward("123", "456")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteCustomRewardReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteCustomReward("", "123")
	require.ErrorIs(t, err, BadRequestError{"invalid request, broadcaster can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteCustomReward("123", "")
	require.ErrorIs(t, err, BadRequestError{"invalid request, id can't be blank"})
	require.False(t, deleted)
}
