# Go Tau
[![pipeline status](https://gitlab.com/wwsean08/go-tau/badges/main/pipeline.svg)](https://gitlab.com/wwsean08/go-tau/-/commits/main)
[![coverage report](https://gitlab.com/wwsean08/go-tau/badges/main/coverage.svg)](https://gitlab.com/wwsean08/go-tau/-/commits/main)
[![Go Reference](https://pkg.go.dev/badge/gitlab.com/wwsean08/go-tau.svg)](https://pkg.go.dev/gitlab.com/wwsean08/go-tau)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/wwsean08/go-tau)](https://goreportcard.com/report/gitlab.com/wwsean08/go-tau)
[![discord](https://img.shields.io/discord/493162062524973056?label=discord)](https://discord.gg/PPVYMeP)

This library is designed to make integrating with [TAU](https://github.com/FiniteSingularity/tau) in go, taking care of parsing out the messages and calling the registered receivers, leaving you to handle the logic.

## Callbacks
You can implement any/all of these callbacks in your own code to take advantage of this library allowing you to take care of the business logic, while this library takes care of parsing messages into usable structs.

* `RawCallback(msg []byte)` - All websocket messages received will be forwarded to this callback untouched.  Note this will not stop processing via more specific callbacks.
* `ErrorCallback(err error)` - Used to handle websocket type errors like when the websocket is closed by the remote source, if this is not handled then errors here will cause a `panic`.
* `StreamOnlineCallback(msg *StreamOnlineMsg)` - Called when a streamer goes online
* `StreamOfflineCallback(msg *StreamOnlineMsg)` - Called when a streamer goes offline
* `FollowCallback(msg *FollowMsg)` - Called when a new follow event is received.
* `StreamUpdateCallback(msg *StreamUpdateMsg)` - Called when a stream update event is received.
* `CheerCallback(msg *CheerMsg)` - Called when a cheer event is received.
* `RaidCallback(msg *RaidMsg)` - Called when a raid event is received.
* `SubscriptionCallback(msg *SubscriptionMsg)` - Called when a subscription event is received.
* `PointsRedemptionCallback(msg *PointsRedemptionMsg)` - Called when a points redemption event is received. _Breaking change coming in a future version based on an upcoming TAU change_
* `HypeTrainBeginCallback(msg *HypeTrainBeginMsg)` - Called when a hype train beginning event is received.
* `HypeTrainProgressCallback(msg *HypeTrainProgressMsg)` - Called when a hype train progress event is received.
* `HypeTrainEndCallback(msg *HypeTrainEndedMsg)` - Called when a hype train end event is received.

## Utility Functions
* `GetAuthToken` - Allows for getting your auth token via username and password.  Ideally you would keep your auth token in your config, but this just gives you another option of how to get the data.
* `Reconnect` - Can be used to reconnect to the websocket on a connection error via the `ErrorCallback`.