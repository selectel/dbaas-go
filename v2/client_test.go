package v2

import (
	"reflect"
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
	require.NotNil(t, api.ClickHouse.DatastoreService)
	require.NotNil(t, api.ClickHouse.NodeGroupService)

	ptrDirect := reflect.ValueOf(api.ClickHouse.GetDatastore).Pointer()
	ptrViaInner := reflect.ValueOf(api.ClickHouse.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirect, ptrViaInner, "Methods must refer to the same function")

	require.NotNil(t, api.DatastoreType)
	require.NotNil(t, api.Flavor)
}
