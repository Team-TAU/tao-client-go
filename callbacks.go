package go_tau

import "fmt"

type RawCallback func(msg []byte)
type FollowCallback func(msg *FollowMsg)
type StreamUpdateCallback func(msg *StreamUpdateMsg)
type CheerCallback func(msg *CheerMsg)
type RaidCallback func(msg *RaidMsg)
type SubscriptionCallback func(msg *SubscriptionMsg)

// SetRawCallback sets a callback to be called on all received messages.
func (c *Client) SetRawCallback(callback RawCallback) {
	c.rawCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetFollowCallback sets a callback to be called on a follow event received.
func (c *Client) SetFollowCallback(callback FollowCallback) {
	c.followCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetStreamUpdateCallback sets a callback to be called on stream update events received
func (c *Client) SetStreamUpdateCallback(callback StreamUpdateCallback) {
	c.streamUpdateCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetCheerCallback sets a callback to be called when a cheer event is received.
func (c *Client) SetCheerCallback(callback CheerCallback) {
	c.cheerCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetRaidCallback sets a callback to be called when a raid event is received.
func (c *Client) SetRaidCallback(callback RaidCallback) {
	c.raidCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}

// SetSubscriptionCallback sets a callback to be called when a subscription event is received.
func (c *Client) SetSubscriptionCallback(callback SubscriptionCallback) {
	c.subscriptionCallback = callback
	// Attempt to fix the heisenbug where if I don't acknowledge the callback it will be null
	// TODO: Figure out an ACTUAL fix
	_ = fmt.Sprintf("%p", c.rawCallback)
}
