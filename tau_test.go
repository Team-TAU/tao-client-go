package gotau

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestGetAuthToken_ValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		require.Equal(t, req.URL.String(), "/api-token-auth/")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("{\"token\": \"baz\"}"))
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

func TestClient_SetParallelProcessing(t *testing.T) {
	client := new(Client)
	require.False(t, client.parallelProcessing)
	client.SetParallelProcessing(true)
	require.True(t, client.parallelProcessing)
	client.SetParallelProcessing(false)
	require.False(t, client.parallelProcessing)
}

func TestClient_login(t *testing.T) {
	upgrader := websocket.Upgrader{}
	wg := new(sync.WaitGroup)
	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		type token struct {
			Token string `json:"token"`
		}
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}
			myToken := new(token)
			err = json.Unmarshal(message, myToken)
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}
			if !assert.Equal(t, "foo", myToken.Token) {
				wg.Done()
				continue
			}
			called += 1
			wg.Done()
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	url := strings.TrimPrefix(server.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)
	client := new(Client)
	client.writeLock = new(sync.Mutex)
	client.hasSSL = false
	client.token = "foo"
	client.hostname = host
	client.port = portNum

	conn, err := connect(client.hostname, client.port, false)
	require.NoError(t, err)
	require.NotNil(t, conn)
	client.conn = conn

	wg.Add(1)
	err = client.login()
	wg.Wait()

	require.NoError(t, err)
	require.Equal(t, 1, called)
}

func TestClient_SendMessage(t *testing.T) {
	upgrader := websocket.Upgrader{}
	wg := new(sync.WaitGroup)
	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}
			if !assert.Equal(t, "\"message\"", strings.TrimSpace(string(message))) {
				wg.Done()
				continue
			}
			called += 1
			wg.Done()
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	url := strings.TrimPrefix(server.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)
	client := new(Client)
	client.writeLock = new(sync.Mutex)

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/ws/twitch-events/", host, portNum), nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	client.conn = conn

	wg.Add(1)
	err = client.SendMessage("message")
	wg.Wait()
	require.NoError(t, err)

	require.Equal(t, 1, called)
}

func TestClient_Reconnect(t *testing.T) {
	upgrader := websocket.Upgrader{}
	wg := new(sync.WaitGroup)
	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		type token struct {
			Token string `json:"token"`
		}
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			_, ok := err.(*websocket.CloseError)
			if ok {
				return
			}
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}
			if called == 0 {
				if !assert.Equal(t, "\"message\"", strings.TrimSpace(string(message))) {
					wg.Done()
					continue
				}
			} else if called == 1 {
				myToken := new(token)
				err = json.Unmarshal(message, myToken)
				if !assert.NoError(t, err) {
					wg.Done()
					continue
				}
				if !assert.Equal(t, "foo", myToken.Token) {
					wg.Done()
					continue
				}
			}
			called += 1
			wg.Done()
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	url := strings.TrimPrefix(server.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)
	client := new(Client)
	client.writeLock = new(sync.Mutex)
	client.hostname = host
	client.port = portNum
	client.hasSSL = false
	client.token = "foo"

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/ws/twitch-events/", host, portNum), nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	client.conn = conn

	wg.Add(1)
	err = client.SendMessage("message")
	wg.Wait()
	require.NoError(t, err)

	err = client.Reconnect()
	require.NoError(t, err)
	wg.Add(1)
	wg.Wait()
	require.NoError(t, err)

	require.Equal(t, 2, called)
}

func TestNewClient(t *testing.T) {
	upgrader := websocket.Upgrader{}
	wg := new(sync.WaitGroup)
	called := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		type token struct {
			Token string `json:"token"`
		}
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}

			myToken := new(token)
			err = json.Unmarshal(message, myToken)
			if !assert.NoError(t, err) {
				wg.Done()
				continue
			}
			if !assert.Equal(t, "foo", myToken.Token) {
				wg.Done()
				continue
			}
			called += 1
			wg.Done()
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	url := strings.TrimPrefix(server.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	wg.Add(1)
	client, err := NewClient(host, portNum, "foo", false)
	wg.Wait()

	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.conn)
	require.Equal(t, 1, called)
}
