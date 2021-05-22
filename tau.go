package go_tau

import (
	"fmt"
	"sync"
)
import "github.com/gorilla/websocket"

// Client represents a client connected to TAU
type Client struct {
	conn      *websocket.Conn
	hostname  string
	port      int
	writeLock *sync.Mutex

	// callback functions
	rawCallback               RawCallback
	errorCallback             ErrorCallback
	followCallback            FollowCallback
	streamUpdateCallback      StreamUpdateCallback
	cheerCallback             CheerCallback
	raidCallback              RaidCallback
	subscriptionCallback      SubscriptionCallback
	hypeTrainBeginCallback    HypeTrainBeginCallback
	hypeTrainProgressCallback HypeTrainProgressCallback
	hypeTrainEndedCallback    HypeTrainEndCallback
}

// NewClient allows you to get a new client that is connected to TAU
func NewClient(hostname string, port int, token string, hasSSL bool) (*Client, error) {
	client := &Client{
		hostname:  hostname,
		port:      port,
		writeLock: new(sync.Mutex),
	}
	prefix := "ws://"
	if hasSSL {
		prefix = "wss://"
	}
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s%s:%d/ws/twitch-events/", prefix, hostname, port), nil)
	if err != nil {
		return nil, err
	}
	client.conn = conn
	login := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	err = client.SendMessage(login)
	if err != nil {
		return nil, err
	}
	go client.readLoop()

	return client, nil
}

func (c *Client) SendMessage(msg interface{}) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.conn.WriteJSON(msg)
}
