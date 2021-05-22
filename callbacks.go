package go_tau

import "fmt"

type RawCallback func(msg []byte)
type ErrorCallback func(err error)
type FollowCallback func(msg *FollowMsg)
type StreamUpdateCallback func(msg *StreamUpdateMsg)
type CheerCallback func(msg *CheerMsg)
type RaidCallback func(msg *RaidMsg)
type SubscriptionCallback func(msg *SubscriptionMsg)
type HypeTrainBeginCallback func(msg *HypeTrainBeginMsg)
type HypeTrainProgressCallback func(msg *HypeTrainProgressMsg)
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
