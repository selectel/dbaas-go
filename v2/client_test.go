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
	require.NotNil(t, api.ClickHouse.ShardGroupService)

	ptrDirectD := reflect.ValueOf(api.ClickHouse.GetDatastore).Pointer()
	ptrViaInnerD := reflect.ValueOf(api.ClickHouse.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirectD, ptrViaInnerD, "Methods must refer to the same function")

	ptrDirectN := reflect.ValueOf(api.ClickHouse.ResizeNodeGroup).Pointer()
	ptrViaInnerN := reflect.ValueOf(api.ClickHouse.NodeGroupService.ResizeNodeGroup).Pointer()
	require.Equal(t, ptrDirectN, ptrViaInnerN, "Methods must refer to the same function")

	ptrDirectS := reflect.ValueOf(api.ClickHouse.CreateShardGroup).Pointer()
	ptrViaInnerS := reflect.ValueOf(api.ClickHouse.ShardGroupService.CreateShardGroup).Pointer()
	require.Equal(t, ptrDirectS, ptrViaInnerS, "Methods must refer to the same function")

	require.NotNil(t, api.DatastoreType)
	require.NotNil(t, api.Flavor)
}
