package go_tau

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_SetRawCallback(t *testing.T) {
	client := Client{}
	callback := func(msg []byte) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.rawCallback)
	client.SetRawCallback(callback)
	require.NotNil(t, client.rawCallback)
}

func TestClient_SetErrorCallback(t *testing.T) {
	client := Client{}
	callback := func(err error) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.errorCallback)
	client.SetErrorCallback(callback)
	require.NotNil(t, client.errorCallback)
}

func TestClient_SetFollowCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *FollowMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.followCallback)
	client.SetFollowCallback(callback)
	require.NotNil(t, client.followCallback)
}

func TestClient_SetStreamUpdateCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *StreamUpdateMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.streamUpdateCallback)
	client.SetStreamUpdateCallback(callback)
	require.NotNil(t, client.streamUpdateCallback)
}

func TestClient_SetCheerCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *CheerMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.cheerCallback)
	client.SetCheerCallback(callback)
	require.NotNil(t, client.cheerCallback)
}

func TestClient_SetRaidCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *RaidMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.raidCallback)
	client.SetRaidCallback(callback)
	require.NotNil(t, client.raidCallback)
}

func TestClient_SetSubscriptionCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *SubscriptionMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.subscriptionCallback)
	client.SetSubscriptionCallback(callback)
	require.NotNil(t, client.subscriptionCallback)
}

func TestClient_SetHypeTrainBeginCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *HypeTrainBeginMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.hypeTrainBeginCallback)
	client.SetHypeTrainBeginCallback(callback)
	require.NotNil(t, client.hypeTrainBeginCallback)
}

func TestClient_SetHypeTrainProgressCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *HypeTrainProgressMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.hypeTrainProgressCallback)
	client.SetHypeTrainProgressCallback(callback)
	require.NotNil(t, client.hypeTrainProgressCallback)
}

func TestClient_SetHypeTrainEndedCallback(t *testing.T) {
	client := Client{}
	callback := func(msg *HypeTrainEndedMsg) {
	}
	require.NotNil(t, callback)
	require.Nil(t, client.hypeTrainEndedCallback)
	client.SetHypeTrainEndedCallback(callback)
	require.NotNil(t, client.hypeTrainEndedCallback)
}
