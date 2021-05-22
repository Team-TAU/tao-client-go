package go_tau

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	type tmp struct {
		Ts Time `json:"ts"`
	}
	timeData := []byte("{\"ts\": \"2021-05-22T18:31:59.588379+00:00\"}")
	timestamp := new(tmp)
	err := json.Unmarshal(timeData, timestamp)
	require.NoError(t, err)
	require.Equal(t, 2021, timestamp.Ts.Year())
	require.Equal(t, time.Month(5), timestamp.Ts.Month())
	require.Equal(t, 22, timestamp.Ts.Day())
	require.Equal(t, 18, timestamp.Ts.Hour())
	require.Equal(t, 31, timestamp.Ts.Minute())
	require.Equal(t, 59, timestamp.Ts.Second())
}

func TestTime_UnmarshalJSONReturnsError(t *testing.T) {
	type tmp struct {
		Ts Time `json:"ts"`
	}
	// Validate that an invalid string causes an error
	timeData := []byte("{\"ts\": \"2021-05-22-18:31:59.588379+00:00\"}")
	timestamp := new(tmp)
	err := json.Unmarshal(timeData, timestamp)
	require.Error(t, err)
}
