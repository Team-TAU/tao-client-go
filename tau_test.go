package go_tau

import (
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetAuthToken_ValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		require.Equal(t, req.URL.String(), "/api-token-auth/")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("{\"token\": \"baz\"}"))
	}))
	defer server.Close()
	url := server.URL
	url = strings.TrimPrefix(url, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)
	token, err := GetAuthToken("foo", "bar", host, portNum, false)
	require.NoError(t, err)
	require.Equal(t, "baz", token)
}

func TestGetAuthToken_InvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		require.Equal(t, req.URL.String(), "/api-token-auth/")
		// Send response to be tested
		rw.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	url := server.URL
	url = strings.TrimPrefix(url, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)
	token, err := GetAuthToken("foo", "bar", host, portNum, false)
	require.Error(t, err)
	require.Equal(t, "", token)
}
