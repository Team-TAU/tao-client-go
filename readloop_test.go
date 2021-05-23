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
		msg                             string
		testCase                        string
		expectFollowCallback            bool
		expectStreamUpdateCallback      bool
		expectStreamOnlineCallback      bool
		expectStreamOfflineCallback     bool
		expectPointsRedemptionCallback  bool
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
		{
			msg:                        "{\"id\":null,\"event_id\":\"78e4825d-7496-44a5-b4f2-d71af9f040d8\",\"event_type\":\"stream-online\",\"event_source\":\"TestCall\",\"event_data\":{\"id\":\"9001\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"type\":\"live\",\"started_at\":\"2020-10-11T10:11:12.123Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                   "Expect Stream Online Callback",
			expectStreamOnlineCallback: true,
		},
		{
			msg:                         "{\"id\":null,\"event_id\":\"097896c2-9de8-4840-bcdf-7bc5cada9c28\",\"event_type\":\"stream-offline\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                    "Expect Stream Offline Callback",
			expectStreamOfflineCallback: true,
		},
		{
			msg:                          "{\"id\":null,\"event_id\":\"a4fc9436-74f6-438a-9cc0-5f3c9bce76cf\",\"event_type\":\"hype-train-begin\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"total\":137,\"progress\":137,\"goal\":500,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                     "Expect Hype Train Begin Callback",
			expectHypeTrainBeginCallback: true,
		},
		{
			msg:                             "{\"id\":null,\"event_id\":\"25c54242-154e-44d0-8d52-8a91c783bca4\",\"event_type\":\"hype-train-progress\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":700,\"progress\":200,\"goal\":1000,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                        "Expect Hype Train Progress Callback",
			expectHypeTrainProgressCallback: true,
		},
		{
			msg:                        "{\"id\":null,\"event_id\":\"f46bdbdc-4890-4199-9e7d-34f4ee9c2bdd\",\"event_type\":\"hype-train-end\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":137,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"ended_at\":\"2020-07-15T17:16:11.17106713Z\",\"cooldown_ends_at\":\"2020-07-15T18:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                   "Expect Hype Train End Callback",
			expectHypeTrainEndCallback: true,
		},
		{
			msg:                            "{\"id\":null,\"event_id\":\"4dc3e5c5-3e38-4d2d-bfe2-9c89c65a8c2a\",\"event_type\":\"point-redemption\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\",\"id\":\"649995ea-b88b-446d-a011-0cc183588bd4\",\"user_name\":\"FiniteSingularity\",\"user_id\":1234,\"user_login\":\"finitesingularity\",\"user_input\":\"\",\"status\":\"unfilled\",\"redeemed_at\":\"2021-05-22T20:36:06.427Z\",\"reward\":{\"id\":\"3859c466-8cff-4480-9e9e-b7e9814b405d\",\"title\":\"Free tier 1 sub\",\"prompt\":\"Sean will gift you a tier one sub to his channel for one month\",\"cost\":20000}},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}",
			testCase:                       "Expect Points Redemption Callback",
			expectPointsRedemptionCallback: true,
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
			require.True(t, test.expectHypeTrainBeginCallback, test.testCase)
		}
		hypeTrainProgress := func(msg *HypeTrainProgressMsg) {
			require.True(t, test.expectHypeTrainProgressCallback, test.testCase)
		}
		hypeTrainEnded := func(msg *HypeTrainEndedMsg) {
			require.True(t, test.expectHypeTrainEndCallback, test.testCase)
		}
		streamOnline := func(msg *StreamOnlineMsg) {
			require.True(t, test.expectStreamOnlineCallback, test.testCase)
		}
		streamOffline := func(msg *StreamOfflineMsg) {
			require.True(t, test.expectStreamOfflineCallback, test.testCase)
		}
		pointsCallback := func(msg *PointsRedemptionMsg) {
			require.True(t, test.expectPointsRedemptionCallback, test.testCase)
		}
		client.SetFollowCallback(followCallback)
		client.SetStreamUpdateCallback(streamUpdateCallback)
		client.SetCheerCallback(cheerCallback)
		client.SetRaidCallback(raidCallback)
		client.SetSubscriptionCallback(subCallback)
		client.SetHypeTrainBeginCallback(hypeTrainBegin)
		client.SetHypeTrainProgressCallback(hypeTrainProgress)
		client.SetHypeTrainEndedCallback(hypeTrainEnded)
		client.SetStreamOnlineCallback(streamOnline)
		client.SetStreamOfflineCallback(streamOffline)
		client.SetPointsRedemptionCallback(pointsCallback)

		client.handleMessage([]byte(test.msg))
	}
}

func TestHandleMessage_FollowEvent(t *testing.T) {
	called := false
	callback := func(msg *FollowMsg) {
		called = true
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
	require.True(t, called)
}

func TestHandleMessage_UpdateEvent(t *testing.T) {
	called := false
	callback := func(msg *StreamUpdateMsg) {
		called = true
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
	require.True(t, called)
}

func TestHandleMessage_CheerEvent(t *testing.T) {
	called := false
	callback := func(msg *CheerMsg) {
		called = true
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
	require.True(t, called)
}

func TestHandleMessage_RaidEvent(t *testing.T) {
	called := false
	callback := func(msg *RaidMsg) {
		called = true
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
	require.True(t, called)
}

func TestHandleMessage_SubEvent(t *testing.T) {
	called := false
	callback := func(msg *SubscriptionMsg) {
		called = true
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
	require.True(t, called)
}

func TestHandleMessage_StreamOnlineEvent(t *testing.T) {
	called := false
	callback := func(msg *StreamOnlineMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "78e4825d-7496-44a5-b4f2-d71af9f040d8", msg.EventID)
		require.Equal(t, "stream-online", msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())

		require.NotNil(t, msg.EventData)
		require.Equal(t, "9001", msg.EventData.ID)
		require.Equal(t, "1337", msg.EventData.BroadcasterID)
		require.Equal(t, "cool_user", msg.EventData.BroadcasterLogin)
		require.Equal(t, "Cool_User", msg.EventData.BroadcasterName)
		require.Equal(t, "live", msg.EventData.Type)

		startedAt := msg.EventData.StartedAt
		require.Equal(t, 2020, startedAt.Year())
		require.Equal(t, time.Month(10), startedAt.Month())
		require.Equal(t, 11, startedAt.Day())
		require.Equal(t, 10, startedAt.Hour())
		require.Equal(t, 11, startedAt.Minute())
		require.Equal(t, 12, startedAt.Second())
	}
	client := Client{
		streamOnlineCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"78e4825d-7496-44a5-b4f2-d71af9f040d8\",\"event_type\":\"stream-online\",\"event_source\":\"TestCall\",\"event_data\":{\"id\":\"9001\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"type\":\"live\",\"started_at\":\"2020-10-11T10:11:12.123Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}

func TestHandleMessage_StreamOfflineEvent(t *testing.T) {
	called := false
	callback := func(msg *StreamOfflineMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "097896c2-9de8-4840-bcdf-7bc5cada9c28", msg.EventID)
		require.Equal(t, "stream-offline", msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())

		require.Equal(t, "1337", msg.EventData.BroadcasterID)
		require.Equal(t, "cool_user", msg.EventData.BroadcasterLogin)
		require.Equal(t, "Cool_User", msg.EventData.BroadcasterName)
	}
	client := Client{
		streamOfflineCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"097896c2-9de8-4840-bcdf-7bc5cada9c28\",\"event_type\":\"stream-offline\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}

func TestHandleMessage_HypeTrainBeginEvent(t *testing.T) {
	called := false
	callback := func(msg *HypeTrainBeginMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "a4fc9436-74f6-438a-9cc0-5f3c9bce76cf", msg.EventID)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())

		require.NotZero(t, msg.EventData)
		require.Equal(t, "1337", msg.EventData.BroadcasterID)
		require.Equal(t, "cool_user", msg.EventData.BroadcasterLogin)
		require.Equal(t, "Cool_User", msg.EventData.BroadcasterName)
		require.Equal(t, 137, msg.EventData.Total)
		require.Equal(t, 137, msg.EventData.Progress)
		require.Equal(t, 500, msg.EventData.Goal)
		require.Len(t, msg.EventData.TopContributions, 2)

		top := msg.EventData.TopContributions[0]
		require.Equal(t, "123", top.UserID)
		require.Equal(t, "pogchamp", top.UserLogin)
		require.Equal(t, "PogChamp", top.UserName)
		require.Equal(t, "bits", top.Type)
		require.Equal(t, 50, top.Total)

		last := msg.EventData.LastContribution
		require.Equal(t, "123", last.UserID)
		require.Equal(t, "pogchamp", last.UserLogin)
		require.Equal(t, "PogChamp", last.UserName)
		require.Equal(t, "bits", last.Type)
		require.Equal(t, 50, last.Total)

		started := msg.EventData.StartedAt
		require.Equal(t, 2020, started.Year())
		require.Equal(t, time.Month(7), started.Month())
		require.Equal(t, 15, started.Day())
		require.Equal(t, 17, started.Hour())
		require.Equal(t, 16, started.Minute())
		require.Equal(t, 03, started.Second())

		expires := msg.EventData.ExpiresAt
		require.Equal(t, 2020, expires.Year())
		require.Equal(t, time.Month(7), expires.Month())
		require.Equal(t, 15, expires.Day())
		require.Equal(t, 17, expires.Hour())
		require.Equal(t, 16, expires.Minute())
		require.Equal(t, 11, expires.Second())
	}
	client := Client{
		hypeTrainBeginCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"a4fc9436-74f6-438a-9cc0-5f3c9bce76cf\",\"event_type\":\"hype-train-begin\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"total\":137,\"progress\":137,\"goal\":500,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}

func TestHandleMessage_HypeTrainProgressEvent(t *testing.T) {
	called := false
	callback := func(msg *HypeTrainProgressMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "25c54242-154e-44d0-8d52-8a91c783bca4", msg.EventID)
		require.Equal(t, "hype-train-progress", msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())
		require.NotNil(t, msg.EventData)
		require.Equal(t, "1337", msg.EventData.BroadcasterID)
		require.Equal(t, "cool_user", msg.EventData.BroadcasterLogin)
		require.Equal(t, "Cool_User", msg.EventData.BroadcasterName)
		require.Equal(t, 2, msg.EventData.Level)
		require.Equal(t, 700, msg.EventData.Total)
		require.Equal(t, 200, msg.EventData.Progress)
		require.Equal(t, 1000, msg.EventData.Goal)
		require.Len(t, msg.EventData.TopContributions, 2)

		top := msg.EventData.TopContributions[0]
		require.Equal(t, "123", top.UserID)
		require.Equal(t, "pogchamp", top.UserLogin)
		require.Equal(t, "PogChamp", top.UserName)
		require.Equal(t, "bits", top.Type)
		require.Equal(t, 50, top.Total)

		last := msg.EventData.LastContribution
		require.Equal(t, "123", last.UserID)
		require.Equal(t, "pogchamp", last.UserLogin)
		require.Equal(t, "PogChamp", last.UserName)
		require.Equal(t, "bits", last.Type)
		require.Equal(t, 50, last.Total)

		started := msg.EventData.StartedAt
		require.Equal(t, 2020, started.Year())
		require.Equal(t, time.Month(7), started.Month())
		require.Equal(t, 15, started.Day())
		require.Equal(t, 17, started.Hour())
		require.Equal(t, 16, started.Minute())
		require.Equal(t, 03, started.Second())

		expires := msg.EventData.ExpiresAt
		require.Equal(t, 2020, expires.Year())
		require.Equal(t, time.Month(7), expires.Month())
		require.Equal(t, 15, expires.Day())
		require.Equal(t, 17, expires.Hour())
		require.Equal(t, 16, expires.Minute())
		require.Equal(t, 11, expires.Second())
	}
	client := Client{
		hypeTrainProgressCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"25c54242-154e-44d0-8d52-8a91c783bca4\",\"event_type\":\"hype-train-progress\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":700,\"progress\":200,\"goal\":1000,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}

func TestHandleMessage_HypeTrainEndEvent(t *testing.T) {
	called := false
	callback := func(msg *HypeTrainEndedMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "f46bdbdc-4890-4199-9e7d-34f4ee9c2bdd", msg.EventID)
		require.Equal(t, "hype-train-end", msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())
		require.Equal(t, "1337", msg.EventData.BroadcasterID)
		require.Equal(t, "cool_user", msg.EventData.BroadcasterLogin)
		require.Equal(t, "Cool_User", msg.EventData.BroadcasterName)
		require.Equal(t, 2, msg.EventData.Level)
		require.Equal(t, 137, msg.EventData.Total)
		require.Len(t, msg.EventData.TopContributions, 2)

		top := msg.EventData.TopContributions[0]
		require.Equal(t, "123", top.UserID)
		require.Equal(t, "pogchamp", top.UserLogin)
		require.Equal(t, "PogChamp", top.UserName)
		require.Equal(t, "bits", top.Type)
		require.Equal(t, 50, top.Total)

		started := msg.EventData.StartedAt
		require.Equal(t, 2020, started.Year())
		require.Equal(t, time.Month(7), started.Month())
		require.Equal(t, 15, started.Day())
		require.Equal(t, 17, started.Hour())
		require.Equal(t, 16, started.Minute())
		require.Equal(t, 03, started.Second())

		endedAt := msg.EventData.EndedAt
		require.Equal(t, 2020, endedAt.Year())
		require.Equal(t, time.Month(7), endedAt.Month())
		require.Equal(t, 15, endedAt.Day())
		require.Equal(t, 17, endedAt.Hour())
		require.Equal(t, 16, endedAt.Minute())
		require.Equal(t, 11, endedAt.Second())

		coolDown := msg.EventData.CooldownEndsAt
		require.Equal(t, 2020, coolDown.Year())
		require.Equal(t, time.Month(7), coolDown.Month())
		require.Equal(t, 15, coolDown.Day())
		require.Equal(t, 18, coolDown.Hour())
		require.Equal(t, 16, coolDown.Minute())
		require.Equal(t, 11, coolDown.Second())
	}
	client := Client{
		hypeTrainEndedCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"f46bdbdc-4890-4199-9e7d-34f4ee9c2bdd\",\"event_type\":\"hype-train-end\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":137,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"ended_at\":\"2020-07-15T17:16:11.17106713Z\",\"cooldown_ends_at\":\"2020-07-15T18:16:11.17106713Z\"},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}

func TestHandleMessage_PointsRedemptionEvent(t *testing.T) {
	called := false
	callback := func(msg *PointsRedemptionMsg) {
		called = true
		require.NotNil(t, msg)
		require.Zero(t, msg.ID)
		require.Equal(t, "4dc3e5c5-3e38-4d2d-bfe2-9c89c65a8c2a", msg.EventID)
		require.Equal(t, "point-redemption", msg.EventType)
		require.Equal(t, "TestCall", msg.EventSource)
		require.Equal(t, "test", msg.Origin)
		require.Equal(t, 2021, msg.Created.Year())
		require.Equal(t, time.Month(5), msg.Created.Month())
		require.Equal(t, 22, msg.Created.Day())
		require.Equal(t, 20, msg.Created.Hour())
		require.Equal(t, 36, msg.Created.Minute())
		require.Equal(t, 7, msg.Created.Second())
		require.Equal(t, "47073625", msg.EventData.BroadcasterID)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterName)
		require.Equal(t, "wwsean08", msg.EventData.BroadcasterLogin)
		require.Equal(t, "649995ea-b88b-446d-a011-0cc183588bd4", msg.EventData.ID)
		require.Equal(t, "FiniteSingularity", msg.EventData.UserName)
		require.Equal(t, "finitesingularity", msg.EventData.UserLogin)
		require.Equal(t, "1234", msg.EventData.UserID)
		require.Empty(t, msg.EventData.UserInput)
		require.Equal(t, "unfilled", msg.EventData.Status)

		redeemedAt := msg.EventData.RedeemedAt
		require.Equal(t, 2021, redeemedAt.Year())
		require.Equal(t, time.Month(5), redeemedAt.Month())
		require.Equal(t, 22, redeemedAt.Day())
		require.Equal(t, 20, redeemedAt.Hour())
		require.Equal(t, 36, redeemedAt.Minute())
		require.Equal(t, 6, redeemedAt.Second())

		reward := msg.EventData.Reward
		require.Equal(t, "3859c466-8cff-4480-9e9e-b7e9814b405d", reward.ID)
		require.Equal(t, "Free tier 1 sub", reward.Title)
		require.Equal(t, "Sean will gift you a tier one sub to his channel for one month", reward.Prompt)
		require.Equal(t, 20000, reward.Cost)
	}
	client := Client{
		pointsRedemptionCallback: callback,
	}
	msg := "{\"id\":null,\"event_id\":\"4dc3e5c5-3e38-4d2d-bfe2-9c89c65a8c2a\",\"event_type\":\"point-redemption\",\"event_source\":\"TestCall\",\"event_data\":{\"broadcaster_user_id\":\"47073625\",\"broadcaster_user_name\":\"wwsean08\",\"broadcaster_user_login\":\"wwsean08\",\"id\":\"649995ea-b88b-446d-a011-0cc183588bd4\",\"user_name\":\"FiniteSingularity\",\"user_id\":\"1234\",\"user_login\":\"finitesingularity\",\"user_input\":\"\",\"status\":\"unfilled\",\"redeemed_at\":\"2021-05-22T20:36:06.427Z\",\"reward\":{\"id\":\"3859c466-8cff-4480-9e9e-b7e9814b405d\",\"title\":\"Free tier 1 sub\",\"prompt\":\"Sean will gift you a tier one sub to his channel for one month\",\"cost\":20000}},\"created\":\"2021-05-22T20:36:07.969806+00:00\",\"origin\":\"test\"}"
	client.handleMessage([]byte(msg))
	require.True(t, called)
}
