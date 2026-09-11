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

	ptrDirectDCh := reflect.ValueOf(api.ClickHouse.GetDatastore).Pointer()
	ptrViaInnerDCh := reflect.ValueOf(api.ClickHouse.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirectDCh, ptrViaInnerDCh, "Methods must refer to the same function")

	ptrDirectNCh := reflect.ValueOf(api.ClickHouse.ResizeNodeGroup).Pointer()
	ptrViaInnerNCh := reflect.ValueOf(api.ClickHouse.NodeGroupService.ResizeNodeGroup).Pointer()
	require.Equal(t, ptrDirectNCh, ptrViaInnerNCh, "Methods must refer to the same function")

	ptrDirectSCh := reflect.ValueOf(api.ClickHouse.CreateShardGroup).Pointer()
	ptrViaInnerSCh := reflect.ValueOf(api.ClickHouse.ShardGroupService.CreateShardGroup).Pointer()
	require.Equal(t, ptrDirectSCh, ptrViaInnerSCh, "Methods must refer to the same function")

	require.NotNil(t, api.Opensearch)
	require.NotNil(t, api.Opensearch.DatastoreService)
	require.NotNil(t, api.Opensearch.NodeGroupService)

	ptrDirectDOs := reflect.ValueOf(api.Opensearch.GetDatastore).Pointer()
	ptrViaInnerDOs := reflect.ValueOf(api.Opensearch.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirectDOs, ptrViaInnerDOs, "Methods must refer to the same function")

	ptrDirectNOs := reflect.ValueOf(api.Opensearch.ResizeNodeGroup).Pointer()
	ptrViaInnerNOs := reflect.ValueOf(api.Opensearch.NodeGroupService.ResizeNodeGroup).Pointer()
	require.Equal(t, ptrDirectNOs, ptrViaInnerNOs, "Methods must refer to the same function")

	require.NotNil(t, api.DatastoreType)
	require.NotNil(t, api.Flavor)
}
