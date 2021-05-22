package go_tau

import (
	"github.com/stretchr/testify/require"
	"testing"
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
