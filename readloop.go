package gotau

import (
	"encoding/json"
)

func (c *Client) readLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if c.errorCallback != nil {
				c.errorCallback(err)
				return
			}
			panic(err.Error())
		}

		if c.parallelProcessing {
			go c.handleMessage(message)
		} else {
			c.handleMessage(message)
		}
	}
}

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
	case follow:
		if c.followCallback != nil {
			followMsg := new(FollowMsg)
			err = json.Unmarshal(msg, followMsg)
			if err != nil {
				return
			}
			c.followCallback(followMsg)
		}
	case update:
		if c.streamUpdateCallback != nil {
			updateMsg := new(StreamUpdateMsg)
			err = json.Unmarshal(msg, updateMsg)
			if err != nil {
				return
			}
			c.streamUpdateCallback(updateMsg)
		}
	case cheer:
		if c.cheerCallback != nil {
			cheerMsg := new(CheerMsg)
			err = json.Unmarshal(msg, cheerMsg)
			if err != nil {
				return
			}
			c.cheerCallback(cheerMsg)
		}
	case raid:
		if c.raidCallback != nil {
			raidMsg := new(RaidMsg)
			err = json.Unmarshal(msg, raidMsg)
			if err != nil {
				return
			}
			c.raidCallback(raidMsg)
		}
	case subscription:
		if c.subscriptionCallback != nil {
			subMsg := new(SubscriptionMsg)
			err = json.Unmarshal(msg, subMsg)
			if err != nil {
				return
			}
			c.subscriptionCallback(subMsg)
		}
	case pointsRedemption:
		if c.pointsRedemptionCallback != nil {
			pointsMsg := new(PointsRedemptionMsg)
			err = json.Unmarshal(msg, pointsMsg)
			if err != nil {
				return
			}
			c.pointsRedemptionCallback(pointsMsg)
		}
	case hypeBegin:
		if c.hypeTrainBeginCallback != nil {
			hypeMsg := new(HypeTrainBeginMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainBeginCallback(hypeMsg)
		}
	case hypeProgress:
		if c.hypeTrainProgressCallback != nil {
			hypeMsg := new(HypeTrainProgressMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainProgressCallback(hypeMsg)
		}
	case hypeEnd:
		if c.hypeTrainEndedCallback != nil {
			hypeMsg := new(HypeTrainEndedMsg)
			err = json.Unmarshal(msg, hypeMsg)
			if err != nil {
				return
			}
			c.hypeTrainEndedCallback(hypeMsg)
		}
	case streamOnline:
		if c.streamOnlineCallback != nil {
			onlineMsg := new(StreamOnlineMsg)
			err = json.Unmarshal(msg, onlineMsg)
			if err != nil {
				return
			}
			c.streamOnlineCallback(onlineMsg)
		}
	case streamOffline:
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
