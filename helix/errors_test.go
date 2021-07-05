package helix

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRateLimitError_Error(t *testing.T) {
	err := RateLimitError{
		err: "this is a test",
	}
	require.Error(t, err)
	require.Equal(t, "this is a test", err.Error())
}

func TestRateLimitError_ResetTime(t *testing.T) {
	now := time.Now()
	err := RateLimitError{
		err:   "this is a test",
		reset: &now,
	}

	require.Error(t, err)
	require.Equal(t, now, *err.ResetTime())
}
