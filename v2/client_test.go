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

	ptrDirectCh := reflect.ValueOf(api.ClickHouse.GetDatastore).Pointer()
	ptrViaInnerCh := reflect.ValueOf(api.ClickHouse.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirectCh, ptrViaInnerCh, "ClickHouse methods must refer to the same function")

	require.NotNil(t, api.Opensearch)
	require.NotNil(t, api.Opensearch.DatastoreService)
	require.NotNil(t, api.Opensearch.NodeGroupService)

	ptrDirectOs := reflect.ValueOf(api.Opensearch.GetDatastore).Pointer()
	ptrViaInnerOs := reflect.ValueOf(api.Opensearch.DatastoreService.GetDatastore).Pointer()
	require.Equal(t, ptrDirectOs, ptrViaInnerOs, "Opensearch methods must refer to the same function")

	require.NotNil(t, api.DatastoreType)
	require.NotNil(t, api.Flavor)
}
