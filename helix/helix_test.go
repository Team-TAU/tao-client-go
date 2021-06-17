package helix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClientReturnsClient(t *testing.T) {
	client, err := NewClient("tau.example.com", 443, "abcdefg", true)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, "tau.example.com", client.hostname)
	require.Equal(t, 443, client.port)
	require.Equal(t, "abcdefg", client.token)
	require.True(t, client.hasSSL)
}
