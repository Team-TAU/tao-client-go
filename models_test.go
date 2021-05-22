package go_tau

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEvent_CreatedAsTime(t *testing.T) {
	event := new(Event)
	event.Created = "2021-05-22T05:20:06.120452+00:00"
	timestamp, err := event.CreatedAsTime()
	require.NoError(t, err)
	require.Equal(t, 2021, timestamp.Year())
	require.Equal(t, time.Month(5), timestamp.Month())
	require.Equal(t, 22, timestamp.Day())
	require.Equal(t, 5, timestamp.Hour())
	require.Equal(t, 20, timestamp.Minute())
	require.Equal(t, 06, timestamp.Second())
}

func TestHypeTrainBeginMsg_StartedAtAsTime(t *testing.T) {
	h := new(HypeTrainBeginMsg)
	data := "{\"event_data\": {\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"total\":137,\"progress\":137,\"goal\":500,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.StartedAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 03, timestamp.Second())
}

func TestHypeTrainBeginMsg_ExpiresAtAsTime(t *testing.T) {
	h := new(HypeTrainBeginMsg)
	data := "{\"event_data\": {\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"total\":137,\"progress\":137,\"goal\":500,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.ExpiresAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 11, timestamp.Second())
}

func TestHypeTrainProgressMsg_StartedAtAsTime(t *testing.T) {
	h := new(HypeTrainProgressMsg)
	data := "{\"event_data\": {\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":700,\"progress\":200,\"goal\":1000,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.StartedAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 03, timestamp.Second())
}

func TestHypeTrainProgressMsg_ExpiresAtAsTime(t *testing.T) {
	h := new(HypeTrainProgressMsg)
	data := "{\"event_data\": {\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":700,\"progress\":200,\"goal\":1000,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"last_contribution\":{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"expires_at\":\"2020-07-15T17:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.ExpiresAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 11, timestamp.Second())
}

func TestHypeTrainEndedMsg_StartedAtAsTime(t *testing.T) {
	h := new(HypeTrainEndedMsg)
	data := "{\"event_data\":{\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":137,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"ended_at\":\"2020-07-15T17:16:11.17106713Z\",\"cooldown_ends_at\":\"2020-07-15T18:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.StartedAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 03, timestamp.Second())
}

func TestHypeTrainEndedMsg_EndedAtAsTime(t *testing.T) {
	h := new(HypeTrainEndedMsg)
	data := "{\"event_data\":{\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":137,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"ended_at\":\"2020-07-15T17:16:11.17106713Z\",\"cooldown_ends_at\":\"2020-07-15T18:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.EndedAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 17, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 11, timestamp.Second())
}

func TestHypeTrainEndedMsg_CooldownEndsAtAsTime(t *testing.T) {
	h := new(HypeTrainEndedMsg)
	data := "{\"event_data\":{\"id\":\"1b0AsbInCHZW2SQFQkCzqN07Ib2\",\"broadcaster_user_id\":\"1337\",\"broadcaster_user_login\":\"cool_user\",\"broadcaster_user_name\":\"Cool_User\",\"level\":2,\"total\":137,\"top_contributions\":[{\"user_id\":\"123\",\"user_login\":\"pogchamp\",\"user_name\":\"PogChamp\",\"type\":\"bits\",\"total\":50},{\"user_id\":\"456\",\"user_login\":\"kappa\",\"user_name\":\"Kappa\",\"type\":\"subscription\",\"total\":45}],\"started_at\":\"2020-07-15T17:16:03.17106713Z\",\"ended_at\":\"2020-07-15T17:16:11.17106713Z\",\"cooldown_ends_at\":\"2020-07-15T18:16:11.17106713Z\"}}"
	err := json.Unmarshal([]byte(data), h)
	require.NoError(t, err)
	timestamp, err := h.CooldownEndsAtAsTime()
	require.NoError(t, err)
	require.Equal(t, 2020, timestamp.Year())
	require.Equal(t, time.Month(7), timestamp.Month())
	require.Equal(t, 15, timestamp.Day())
	require.Equal(t, 18, timestamp.Hour())
	require.Equal(t, 16, timestamp.Minute())
	require.Equal(t, 11, timestamp.Second())
}
