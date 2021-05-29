package gotau

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)
import "github.com/gorilla/websocket"

// Client represents a client connected to TAU
type Client struct {
	conn               *websocket.Conn
	hostname           string
	port               int
	hasSSL             bool
	token              string
	writeLock          *sync.Mutex
	parallelProcessing bool

	// callback functions
	rawCallback               RawCallback
	errorCallback             ErrorCallback
	streamOnlineCallback      StreamOnlineCallback
	streamOfflineCallback     StreamOfflineCallback
	followCallback            FollowCallback
	streamUpdateCallback      StreamUpdateCallback
	cheerCallback             CheerCallback
	raidCallback              RaidCallback
	subscriptionCallback      SubscriptionCallback
	pointsRedemptionCallback  PointsRedemptionCallback
	hypeTrainBeginCallback    HypeTrainBeginCallback
	hypeTrainProgressCallback HypeTrainProgressCallback
	hypeTrainEndedCallback    HypeTrainEndCallback
}

// NewClient allows you to get a new client that is connected to TAU
func NewClient(hostname string, port int, token string, hasSSL bool) (*Client, error) {
	client := &Client{
		hostname:           hostname,
		port:               port,
		hasSSL:             hasSSL,
		token:              token,
		writeLock:          new(sync.Mutex),
		parallelProcessing: false,
	}
	conn, err := connect(hostname, port, hasSSL)
	client.conn = conn
	err = client.login()
	if err != nil {
		return nil, err
	}
	go client.readLoop()

	return client, nil
}

// SendMessage Allows you to send json message to the server.
func (c *Client) SendMessage(msg interface{}) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.conn.WriteJSON(msg)
}

// SetParallelProcessing Allows you to enable processing events in parallel.  By default this is false, and most people
//probably would want it to be false, but there could be cases where processing in parallel would be useful/desirable.
func (c *Client) SetParallelProcessing(parallel bool) {
	c.parallelProcessing = parallel
}

// Reconnect can be used to reconnect if a connection error comes in via the ErrorCallback.
func (c *Client) Reconnect() error {
	if c.conn != nil {
		_ = c.conn.Close()
	}
	conn, err := connect(c.hostname, c.port, c.hasSSL)
	if err != nil {
		return err
	}
	c.conn = conn
	return c.login()
}

func (c *Client) login() error {
	login := struct {
		Token string `json:"token"`
	}{
		Token: c.token,
	}
	return c.SendMessage(login)
}

func connect(hostname string, port int, hasSSL bool) (*websocket.Conn, error) {
	prefix := "ws://"
	if hasSSL {
		prefix = "wss://"
	}
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s%s:%d/ws/twitch-events/", prefix, hostname, port), nil)
	return conn, err
}

// GetAuthToken is used to get the auth token for a user to interact with TAU given a username and password.
//Ideally this would be gathered from the UI and potentially stored in a config of some sort, but this option exists
//in case that is not an option.
func GetAuthToken(username, password, hostname string, port int, hasSSL bool) (string, error) {
	protocol := "http"
	if hasSSL {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s:%d/api-token-auth/", protocol, hostname, port)
	body := fmt.Sprintf("{\"username\": \"%s\",\"password\": \"%s\"}", username, password)
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("expected 200 status code but got %d, response body %s", resp.StatusCode, string(data))
	}

	type tmp struct {
		Token string `json:"token"`
	}
	data := new(tmp)
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(rawData, data)
	if err != nil {
		return "", err
	}

	return data.Token, nil
}
