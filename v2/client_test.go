package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPI(t *testing.T) {
	t.Parallel()

	api, err := NewAPI(
		"token",
		"https://api.example.com/v2",
	)

	require.NoError(t, err)
	require.NotNil(t, api)

	require.NotNil(t, api.ClickHouse)
	require.NotNil(t, api.ClickHouse.Datastore)
	require.NotNil(t, api.ClickHouse.NodeGroup)

	require.NotNil(t, api.DatastoreTypes)
}
