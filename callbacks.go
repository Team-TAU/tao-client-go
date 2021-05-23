package gotau

import "fmt"

// RawCallback is a callback that will be called for every message that is received.  Setting this callback does
//not prevent other callbacks from being called, though this one will be called first.
type RawCallback func(msg []byte)

// ErrorCallback is a callback to handle websocket specific errors where the websocket likely needs to be re-established
//and a new client is needed.
type ErrorCallback func(err error)

// StreamOnlineCallback is a callback to handle the stream coming online event
type StreamOnlineCallback func(msg *StreamOnlineMsg)

// StreamOfflineCallback is a callback to handle the stream going offline event
type StreamOfflineCallback func(msg *StreamOfflineMsg)

// FollowCallback is a callback to handle follow events
type FollowCallback func(msg *FollowMsg)

// StreamUpdateCallback is a callback to handle updates to the stream information (like title)
type StreamUpdateCallback func(msg *StreamUpdateMsg)

// CheerCallback is a callback to handle cheer events
type CheerCallback func(msg *CheerMsg)

// RaidCallback is a callback to handle raid events
type RaidCallback func(msg *RaidMsg)

// SubscriptionCallback is a callback to handle subscription events
type SubscriptionCallback func(msg *SubscriptionMsg)

// PointsRedemptionCallback is a callback to handle points redemption events
type PointsRedemptionCallback func(msg *PointsRedemptionMsg)

// HypeTrainBeginCallback is a callback to handle hype train begin events
type HypeTrainBeginCallback func(msg *HypeTrainBeginMsg)

// HypeTrainProgressCallback is a callback to handle hype train progress events
type HypeTrainProgressCallback func(msg *HypeTrainProgressMsg)

// HypeTrainEndCallback is a callback to handle hype train end events
type HypeTrainEndCallback func(msg *HypeTrainEndedMsg)

// SetRawCallback sets a callback to be called on all received messages.
func (c *Client) SetRawCallback(callback RawCallback) {
	c.rawCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetErrorCallback sets a callback for the user to handle errors that might arise (like connection closed)
func (c *Client) SetErrorCallback(callback ErrorCallback) {
	c.errorCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.errorCallback)
}

// SetStreamOnlineCallback sets a callback to be called when a stream online event is received.
func (c *Client) SetStreamOnlineCallback(callback StreamOnlineCallback) {
	c.streamOnlineCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.streamOnlineCallback)
}

// SetStreamOfflineCallback sets a callback to be called when a stream offline event is received.
func (c *Client) SetStreamOfflineCallback(callback StreamOfflineCallback) {
	c.streamOfflineCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.streamOfflineCallback)
}

// SetFollowCallback sets a callback to be called on a follow event received.
func (c *Client) SetFollowCallback(callback FollowCallback) {
	c.followCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.followCallback)
}

// SetStreamUpdateCallback sets a callback to be called on stream update events received
func (c *Client) SetStreamUpdateCallback(callback StreamUpdateCallback) {
	c.streamUpdateCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.streamUpdateCallback)
}

// SetCheerCallback sets a callback to be called when a cheer event is received.
func (c *Client) SetCheerCallback(callback CheerCallback) {
	c.cheerCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.cheerCallback)
}

// SetRaidCallback sets a callback to be called when a raid event is received.
func (c *Client) SetRaidCallback(callback RaidCallback) {
	c.raidCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.raidCallback)
}

// SetSubscriptionCallback sets a callback to be called when a subscription event is received.
func (c *Client) SetSubscriptionCallback(callback SubscriptionCallback) {
	c.subscriptionCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.subscriptionCallback)
}

// SetPointsRedemptionCallback sets a callback to be called when a points redemption event is received.
func (c *Client) SetPointsRedemptionCallback(callback PointsRedemptionCallback) {
	c.pointsRedemptionCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.pointsRedemptionCallback)
}

// SetHypeTrainBeginCallback sets a callback to be called when a hype train begin event is received.
func (c *Client) SetHypeTrainBeginCallback(callback HypeTrainBeginCallback) {
	c.hypeTrainBeginCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.hypeTrainBeginCallback)
}

// SetHypeTrainProgressCallback sets a callback to be called when a hype train progress event is received.
func (c *Client) SetHypeTrainProgressCallback(callback HypeTrainProgressCallback) {
	c.hypeTrainProgressCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.hypeTrainProgressCallback)
}

// SetHypeTrainEndedCallback sets a callback to be called when a hype train ended event is received.
func (c *Client) SetHypeTrainEndedCallback(callback HypeTrainEndCallback) {
	c.hypeTrainEndedCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.hypeTrainEndedCallback)
}
