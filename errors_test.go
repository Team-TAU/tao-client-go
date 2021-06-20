package gotau

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthorizationError_Error(t *testing.T) {
	err := AuthorizationError{
		err: "User unauthorized",
	}

	require.Error(t, err)
	require.Equal(t, "User unauthorized", err.Error())
}

func TestBadRequestError_Error(t *testing.T) {
	err := BadRequestError{
		err: "Malformed request",
	}

	require.Error(t, err)
	require.Equal(t, "Malformed request", err.Error())
}

func TestGenericError_Error(t *testing.T) {
	err := GenericError{
		err: "Generic Error",
	}

	require.Error(t, err)
	require.Equal(t, "Generic Error", err.Error())
}
