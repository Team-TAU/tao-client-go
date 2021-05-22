package go_tau

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"os"
)

func (c *Client) readLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			closeErr, ok := err.(*websocket.CloseError)
			if ok {
				if closeErr.Code == websocket.CloseNormalClosure || closeErr.Code == websocket.CloseGoingAway {
					os.Exit(0)
				}
			}
			panic(err.Error())
		}

		c.handleMessage(message)
	}
}

//TODO: Log parsing errors
func (c *Client) handleMessage(msg []byte) {
	if c.rawCallback != nil {
		c.rawCallback(msg)
	}

	event := new(Event)
	err := json.Unmarshal(msg, event)
	if err != nil {
		// Not much we can do here, skip
		return
	}

	switch event.EventType {
	case FOLLOWEVENT:
		if c.followCallback != nil {
			followMsg := new(FollowMsg)
			err = json.Unmarshal(msg, followMsg)
			if err != nil {
				return
			}
			c.followCallback(followMsg)
		}
	case UPDATEEVENT:
		if c.streamUpdateCallback != nil {
			updateMsg := new(StreamUpdateMsg)
			err = json.Unmarshal(msg, updateMsg)
			if err != nil {
				return
			}
			c.streamUpdateCallback(updateMsg)
		}
	case CHEEREVENT:
		if c.cheerCallback != nil {
			cheerMsg := new(CheerMsg)
			err = json.Unmarshal(msg, cheerMsg)
			if err != nil {
				return
			}
			c.cheerCallback(cheerMsg)
		}
	case RAIDEVENT:
		if c.raidCallback != nil {
			raidMsg := new(RaidMsg)
			err = json.Unmarshal(msg, raidMsg)
			if err != nil {
				return
			}
			c.raidCallback(raidMsg)
		}
	case SUBSCRIPTIONEVENT:
		if c.subscriptionCallback != nil {
			subMsg := new(SubscriptionMsg)
			err = json.Unmarshal(msg, subMsg)
			if err != nil {
				return
			}
			c.subscriptionCallback(subMsg)
		}
	}
}
