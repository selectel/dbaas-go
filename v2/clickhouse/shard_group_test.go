package clickhouse

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
)

const sgID = "000e8400-e29b-41d4-a716-446655440002"

//nolint:gochecknoglobals
var simpleShardGroupResponse = `{
	"id": "000e8400-e29b-41d4-a716-446655440002",
	"name": "shardGR1"
}`

var shardGroupEndpoint = datastoreEndpoint + "/shard-groups/" + sgID //nolint:gochecknoglobals

func newShardGroupService(t *testing.T, serverURL string) *ShardGroupService {
	t.Helper()

	client, err := transport.NewHTTPClient(http.DefaultClient, "token", serverURL+"/v2")
	require.NoError(t, err)

	engine := common.NewEngineService(
		client,
		common.EngineClickHouse,
	)

	return &ShardGroupService{
		EngineService: engine,
	}
}

func newShardGroupServiceWithMockClient() *ShardGroupService {
	mockClient := mockClient{}

	engine := common.NewEngineService(
		mockClient,
		common.EngineClickHouse,
	)

	return &ShardGroupService{
		EngineService: engine,
	}
}

const testShardGroups = `{
    "shard_groups": [
        {
            "id": "2097b78b-6636-4bd9-bf3e-51e189d67329",
            "name": "shardgroup1",
            "shard_ids": [
                "24a9cbdf-5c39-4984-941f-173dd7163b7a"
            ],
            "description": ""
        },
        {
            "id": "2097b78b-6636-4bd9-bf3e-51e189d67300",
            "name": "shardgroup2",
            "shard_ids": [
                "24a9cbdf-5c39-4984-941f-173dd7163b7a",
				"24a9cbdf-5c39-4984-941f-173dd7163b7b"
            ],
            "description": "Test"
        }
    ]
}
`

var ShardGroupsExpected = []ShardGroupResponse{ //nolint:gochecknoglobals
	{
		ID:       "2097b78b-6636-4bd9-bf3e-51e189d67329",
		Name:     "shardgroup1",
		ShardIDs: []string{"24a9cbdf-5c39-4984-941f-173dd7163b7a"},
	},
	{
		ID:          "2097b78b-6636-4bd9-bf3e-51e189d67300",
		Name:        "shardgroup2",
		ShardIDs:    []string{"24a9cbdf-5c39-4984-941f-173dd7163b7a", "24a9cbdf-5c39-4984-941f-173dd7163b7b"},
		Description: "Test",
	},
}

func TestShardGroupService_GetShardGroups_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse/04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1/shard-groups", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testShardGroups))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newShardGroupService(t, server.URL)

	resp, err := srv.GetShardGroups(context.Background(), "04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1")

	require.NoError(t, err)

	require.Equal(t, ShardGroupsExpected, resp)
}

func TestShardGroupService_CreateShardGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, datastoreEndpoint+"/shard-groups", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req ShardGroupCreateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Len(t, req.ShardIDs, 1)
		require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", req.ShardIDs[0])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Mock response from API
		_, err = w.Write([]byte(simpleShardGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newShardGroupService(t, server.URL)

	req := ShardGroupCreateRequest{
		Name:     "shardGR1",
		ShardIDs: []string{"550e8400-e29b-41d4-a716-446655440000"},
	}

	result, err := srv.CreateShardGroup(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, sgID, result.ID)
	require.Equal(t, "shardGR1", result.Name)
}

func TestShardGroupService_CreateShardGroup_InvalidRequest(t *testing.T) {
	srv := newShardGroupServiceWithMockClient()
	body := ShardGroupCreateRequest{
		Name:     "shardGR1",
		ShardIDs: []string{},
	}
	_, err := srv.CreateShardGroup(context.Background(), dsID, body)
	require.Error(t, err)
	require.Equal(
		t, "validate body: shard_group.shard_ids: must be at least one value", err.Error(),
	)
}

func TestShardGroupService_DeleteShardGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, shardGroupEndpoint, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Empty(t, body)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	srv := newShardGroupService(t, server.URL)

	err := srv.DeleteShardGroup(context.Background(), dsID, sgID)

	require.NoError(t, err)
}

func TestShardGroupService_DeleteShardGroup_InvalidRequest(t *testing.T) {
	srv := newShardGroupServiceWithMockClient()

	err := srv.DeleteShardGroup(context.Background(), dsID, "sgid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate shard group id:")
}

func TestShardGroupService_UpdateShardGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, shardGroupEndpoint, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req ShardGroupUpdateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Len(t, *req.ShardIDs, 1)
		require.Empty(t, req.Description)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleShardGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newShardGroupService(t, server.URL)

	shIDs := []string{"550e8400-e29b-41d4-a716-446655440000"}
	req := ShardGroupUpdateRequest{ShardIDs: &shIDs}

	result, err := srv.UpdateShardGroup(context.Background(), dsID, sgID, req)

	require.NoError(t, err)
	require.Equal(t, sgID, result.ID)
	require.Equal(t, "shardGR1", result.Name)
}

func TestShardGroupService_UpdateShardGroup_InvalidRequest(t *testing.T) {
	srv := newShardGroupServiceWithMockClient()
	shIDs := []string{"1"}
	body := ShardGroupUpdateRequest{ShardIDs: &shIDs}
	_, err := srv.UpdateShardGroup(context.Background(), dsID, sgID, body)
	require.Error(t, err)
	require.Equal(t, "validate body: shard_group.shard_ids[0]: invalid UUID length: 1", err.Error())
}
