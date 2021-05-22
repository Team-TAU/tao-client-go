package go_tau

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Basic testing to verify that messages get routed properly.  More specific validation
//on the message done in individual tests.
func TestHandleMessage(t *testing.T) {
	client := Client{}
	type testData struct {
		msg                        string
		testCase                   string
		expectFollowCallback       bool
		expectStreamUpdateCallback bool
		//expectPointsRedemptionCallback  bool
		expectCheerCallback             bool
		expectRaidCallback              bool
		expectSubscriptionCallback      bool
		expectHypeTrainBeginCallback    bool
		expectHypeTrainProgressCallback bool
		expectHypeTrainEndCallback      bool
	}

	data := []testData{
		{
			msg:                        "{\"id\":null,\"event_id\":\"286244d9-c382-4a6e-81ed-5e80bcd2c94e\",\"event_type\":\"update\",\"event_source\":\"TestCall\",\"event_data\":{\"title\":\"foo\",\"languate\":\"en\",\"is_mature\":false,\"category_id\":12345,\"category_name\":\"Science & Technology\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\"},\"created\":\"2021-05-22T04:56:23.545683+00:00\",\"origin\":\"test\"}",
			testCase:                   "Expect Stream Update Callback",
			expectStreamUpdateCallback: true,
		},
		{
			msg:                  "{\"id\":null,\"event_id\":\"69db248a-b980-4de6-ad06-b528feb1294e\",\"event_type\":\"follow\",\"event_source\":\"TestCall\",\"event_data\":{\"user_name\":\"FiniteSingularity\",\"user_id\":\"\",\"user_login\":\"finitesingularity\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\"},\"created\":\"2021-05-22T05:16:21.506340+00:00\",\"origin\":\"test\"}",
			testCase:             "Expect Follow Callback",
			expectFollowCallback: true,
		},
		{
			msg:                 "{\"id\":null,\"event_id\":\"ae567342-62c8-4b45-a41b-7da1472003b9\",\"event_type\":\"cheer\",\"event_source\":\"TestCall\",\"event_data\":{\"is_anonymous\":false,\"user_id\":\"536397236\",\"user_name\":\"FiniteSingularity\",\"user_login\":\"finitesingularity\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\",\"bits\":1000,\"message\":\"hello world\"},\"created\":\"2021-05-22T05:17:40.208431+00:00\",\"origin\":\"test\"}",
			testCase:            "Expect Cheer Callback",
			expectCheerCallback: true,
		},
		{
			msg:                "{\"id\":null,\"event_id\":\"3e62303f-55d3-478d-8f09-c83712e7c3b8\",\"event_type\":\"raid\",\"event_source\":\"TestCall\",\"event_data\":{\"from_broadcaster_user_name\":\"FiniteSingularity\",\"from_broadcaster_user_id\":\"536397236\",\"from_broadcaster_user_login\":\"finitesingularity\",\"to_broadcaster_user_id\":\"47073625\",\"to_broadcaster_user_login\":\"wwsean08\",\"to_broadcaster_user_name\":\"wwsean08\",\"viewers\":42},\"created\":\"2021-05-22T05:19:05.406987+00:00\",\"origin\":\"test\"}",
			testCase:           "Expect Raid Callback",
			expectRaidCallback: true,
		},
		{
			msg:                        "{\"id\":null,\"event_id\":null,\"event_type\":\"subscribe\",\"event_source\":\"TestCall\",\"event_data\":{\"data\":{\"topic\":\"channel-subscribe-events-v1.47073625\",\"message\":{\"benefit_end_month\":0,\"user_name\":\"finitesingularity\",\"display_name\":\"FiniteSingularity\",\"channel_name\":\"wwsean08\",\"user_id\":\"536397236\",\"channel_id\":\"47073625\",\"time\":\"2021-05-22T05:20:06.015Z\",\"sub_message\":{\"message\":\"hello world\",\"emotes\":null},\"sub_plan\":\"1000\",\"sub_plan_name\":\"Channel Subscription (wwsean08)\",\"months\":0,\"cumulative_months\":42,\"context\":\"resub\",\"is_gift\":false,\"multi_month_duration\":0,\"streak_months\":42}},\"type\":\"MESSAGE\"},\"created\":\"2021-05-22T05:20:06.120452+00:00\",\"origin\":\"test\"}",
			testCase:                   "Expect Subscription Callback",
			expectSubscriptionCallback: true,
		},
	}

	for _, test := range data {
		followCallback := func(msg *FollowMsg) {
			require.True(t, test.expectFollowCallback, test.testCase)
		}
		streamUpdateCallback := func(msg *StreamUpdateMsg) {
			require.True(t, test.expectStreamUpdateCallback, test.testCase)
		}
		cheerCallback := func(msg *CheerMsg) {
			require.True(t, test.expectCheerCallback, test.testCase)
		}
		raidCallback := func(msg *RaidMsg) {
			require.True(t, test.expectRaidCallback, test.testCase)
		}
		subCallback := func(msg *SubscriptionMsg) {
			require.True(t, test.expectSubscriptionCallback, test.testCase)
		}
		hypeTrainBegin := func(msg *HypeTrainBeginMsg) {
			require.True(t, test.expectHypeTrainBeginCallback)
		}
		hypeTrainProgress := func(msg *HypeTrainProgressMsg) {
			require.True(t, test.expectHypeTrainProgressCallback)
		}
		hypeTrainEnded := func(msg *HypeTrainEndedMsg) {
			require.True(t, test.expectHypeTrainEndCallback)
		}
		client.SetFollowCallback(followCallback)
		client.SetStreamUpdateCallback(streamUpdateCallback)
		client.SetCheerCallback(cheerCallback)
		client.SetRaidCallback(raidCallback)
		client.SetSubscriptionCallback(subCallback)
		client.SetHypeTrainBeginCallback(hypeTrainBegin)
		client.SetHypeTrainProgressCallback(hypeTrainProgress)
		client.SetHypeTrainEndedCallback(hypeTrainEnded)

		client.handleMessage([]byte(test.msg))
	}
}

func TestHandleMessage_FollowEvent(t *testing.T) {
	callback := func(msg *FollowMsg) {
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "69db248a-b980-4de6-ad06-b528feb1294e", msg.EventID)
		require.Equal(t, FOLLOWEVENT, msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 5, msg.Created.Hour())
		require.Equal(t, 16, msg.Created.Minute())
		require.Equal(t, 21, msg.Created.Second())
		require.NotNil(t, msg.EventData)
		require.Equal(t, "FiniteSingularity", msg.EventData.UserName)
		require.Zero(t, msg.EventData.UserID)
		require.Equal(t, "finitesingularity", msg.EventData.UserLogin)
		require.Equal(t, "47073625", msg.EventData.BroadcasterID)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterName)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterLogin)
	}
	client := Client{
		followCallback: callback,
	}

	msg := "{\"id\":null,\"event_id\":\"69db248a-b980-4de6-ad06-b528feb1294e\",\"event_type\":\"follow\",\"event_source\":\"TestCall\",\"event_data\":{\"user_name\":\"FiniteSingularity\",\"user_id\":\"\",\"user_login\":\"finitesingularity\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\"},\"created\":\"2021-05-22T05:16:21.506340+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
}

func TestHandleMessage_UpdateEvent(t *testing.T) {
	callback := func(msg *StreamUpdateMsg) {
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "286244d9-c382-4a6e-81ed-5e80bcd2c94e", msg.EventID)
		require.Equal(t, UPDATEEVENT, msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 4, msg.Created.Hour())
		require.Equal(t, 56, msg.Created.Minute())
		require.Equal(t, 23, msg.Created.Second())
		require.Equal(t, "foo", msg.EventData.Title)
		require.Equal(t, "en", msg.EventData.Language)
		require.False(t, msg.EventData.IsMature)
		require.Equal(t, 12345, msg.EventData.CategoryID)
		require.Equal(t, "Science & Technology", msg.EventData.CategoryName)
		require.Equal(t, "47073625", msg.EventData.BroadcasterID)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterName)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterLogin)
	}
	client := Client{
		streamUpdateCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"286244d9-c382-4a6e-81ed-5e80bcd2c94e\",\"event_type\":\"update\",\"event_source\":\"TestCall\",\"event_data\":{\"title\":\"foo\",\"language\":\"en\",\"is_mature\":false,\"category_id\":12345,\"category_name\":\"Science & Technology\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\"},\"created\":\"2021-05-22T04:56:23.545683+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
}

func TestHandleMessage_CheerEvent(t *testing.T) {
	callback := func(msg *CheerMsg) {
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "ae567342-62c8-4b45-a41b-7da1472003b9", msg.EventID)
		require.Equal(t, CHEEREVENT, msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 5, msg.Created.Hour())
		require.Equal(t, 17, msg.Created.Minute())
		require.Equal(t, 40, msg.Created.Second())
		require.NotNil(t, msg.EventData)
		require.False(t, msg.EventData.IsAnonymous)
		require.Equal(t, "536397236", msg.EventData.UserID)
		require.Equal(t, "FiniteSingularity", msg.EventData.UserName)
		require.Equal(t, "finitesingularity", msg.EventData.UserLogin)
		require.Equal(t, "47073625", msg.EventData.BroadcasterID)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterName)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterLogin)
		require.Equal(t, 1000, msg.EventData.Bits)
		require.Equal(t, "hello world", msg.EventData.Message)

	}
	client := Client{
		cheerCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"ae567342-62c8-4b45-a41b-7da1472003b9\",\"event_type\":\"cheer\",\"event_source\":\"TestCall\",\"event_data\":{\"is_anonymous\":false,\"user_id\":\"536397236\",\"user_name\":\"FiniteSingularity\",\"user_login\":\"finitesingularity\",\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\",\"bits\":1000,\"message\":\"hello world\"},\"created\":\"2021-05-22T05:17:40.208431+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
}

func TestHandleMessage_RaidEvent(t *testing.T) {
	callback := func(msg *RaidMsg) {
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "3e62303f-55d3-478d-8f09-c83712e7c3b8", msg.EventID)
		require.Equal(t, RAIDEVENT, msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 5, msg.Created.Hour())
		require.Equal(t, 19, msg.Created.Minute())
		require.Equal(t, 5, msg.Created.Second())
		require.NotNil(t, msg.EventData)
		require.Equal(t, "536397236", msg.EventData.FromBroadcasterID)
		require.Equal(t, "FiniteSingularity", msg.EventData.FromBroadcasterName)
		require.Equal(t, "finitesingularity", msg.EventData.FromBroadcasterLogin)
		require.Equal(t, "47073625", msg.EventData.ToBroadcasterID)
		require.Equal(t, "wwsean08", msg.EventData.ToBroadcasterLogin)
		require.Equal(t, "wwsean08", msg.EventData.ToBroadcasterName)
		require.Equal(t, 42, msg.EventData.Viewers)
	}
	client := Client{
		raidCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"3e62303f-55d3-478d-8f09-c83712e7c3b8\",\"event_type\":\"raid\",\"event_source\":\"TestCall\",\"event_data\":{\"from_broadcaster_user_name\":\"FiniteSingularity\",\"from_broadcaster_user_id\":\"536397236\",\"from_broadcaster_user_login\":\"finitesingularity\",\"to_broadcaster_user_id\":\"47073625\",\"to_broadcaster_user_login\":\"wwsean08\",\"to_broadcaster_user_name\":\"wwsean08\",\"viewers\":42},\"created\":\"2021-05-22T05:19:05.406987+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
}

func TestHandleMessage_SubEvent(t *testing.T) {
	callback := func(msg *SubscriptionMsg) {
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Zero(t, msg.EventID)
		require.Equal(t, SUBSCRIPTIONEVENT, msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 5, msg.Created.Hour())
		require.Equal(t, 20, msg.Created.Minute())
		require.Equal(t, 6, msg.Created.Second())
		require.Equal(t, "MESSAGE", msg.EventData.Type)

		data := msg.EventData.Data
		require.NotNil(t, data)
		require.Equal(t, "channel-subscribe-events-v1.47073625", data.Topic)

		message := data.Message
		require.NotNil(t, message)
		require.Equal(t, 0, message.BenefitEndMonth)
		require.Equal(t, "finitesingularity", message.UserName)
		require.Equal(t, "FiniteSingularity", message.DisplayName)
		require.Equal(t, "536397236", message.UserID)
		require.Equal(t, "wwsean08", message.ChannelName)
		require.Equal(t, "47073625", message.ChannelID)
		require.Equal(t, 2021, message.Time.Year())
		require.Equal(t, time.Month(5), message.Time.Month())
		require.Equal(t, 22, message.Time.Day())
		require.Equal(t, 5, message.Time.Hour())
		require.Equal(t, 20, message.Time.Minute())
		require.Equal(t, 6, message.Time.Second())
		require.Equal(t, "1000", message.SubPlan)
		require.Equal(t, "Channel Subscription (wwsean08)", message.SubPlanName)
		require.Equal(t, 0, message.Months)
		require.Equal(t, 42, message.CumulativeMonths)
		require.Equal(t, "resub", message.Context)
		require.False(t, message.IsGift)
		require.Equal(t, 0, message.MultiMonthDuration)
		require.Equal(t, 42, message.StreakMonths)

		subMessage := message.SubMessage
		require.NotNil(t, subMessage)
		require.Equal(t, "hello world", subMessage.Message)
		require.Nil(t, subMessage.Emotes)
	}
	client := Client{
		subscriptionCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":null,\"event_type\":\"subscribe\",\"event_source\":\"TestCall\",\"event_data\":{\"data\":{\"topic\":\"channel-subscribe-events-v1.47073625\",\"message\":{\"benefit_end_month\":0,\"user_name\":\"finitesingularity\",\"display_name\":\"FiniteSingularity\",\"channel_name\":\"wwsean08\",\"user_id\":\"536397236\",\"channel_id\":\"47073625\",\"time\":\"2021-05-22T05:20:06.015Z\",\"sub_message\":{\"message\":\"hello world\",\"emotes\":null},\"sub_plan\":\"1000\",\"sub_plan_name\":\"Channel Subscription (wwsean08)\",\"months\":0,\"cumulative_months\":42,\"context\":\"resub\",\"is_gift\":false,\"multi_month_duration\":0,\"streak_months\":42}},\"type\":\"MESSAGE\"},\"created\":\"2021-05-22T05:20:06.120452+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
}
