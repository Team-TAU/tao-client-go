package go_tau

import (
	"encoding/json"
)

func (c *Client) readLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if c.errorCallback != nil {
				c.errorCallback(err)
			} else {
				panic(err.Error())
			}
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
	case HYPEBEGIN:
		if c.hypeTrainBeginCallback != nil {
			hypeMsg := new(HypeTrainBeginMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainBeginCallback(hypeMsg)
		}
	case HYPEPROGRESS:
		if c.hypeTrainProgressCallback != nil {
			hypeMsg := new(HypeTrainProgressMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainProgressCallback(hypeMsg)
		}
	case HYPEEND:
		if c.hypeTrainEndedCallback != nil {
			hypeMsg := new(HypeTrainEndedMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainEndedCallback(hypeMsg)
		}
	case STREAMONLINE:
		if c.streamOnlineCallback != nil {
			onlineMsg := new(StreamOnlineMsg)
			err = json.Unmarshal(msg, onlineMsg)
			if err != nil {
				return
			}
			c.streamOnlineCallback(onlineMsg)
		}
	case STREAMOFFLINE:
		if c.streamOfflineCallback != nil {
			offlineMsg := new(StreamOfflineMsg)
			err = json.Unmarshal(msg, offlineMsg)
			if err != nil {
				return
			}
			c.streamOfflineCallback(offlineMsg)
		}
	}
}
