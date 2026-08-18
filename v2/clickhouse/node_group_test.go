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
	"github.com/selectel/dbaas-go/v2/internal"
)

const ngID = "000e8400-e29b-41d4-a716-446655440001"

//nolint:gochecknoglobals
var simpleNodeGroupResponse = `{
	"id": "000e8400-e29b-41d4-a716-446655440001",
	"name": "NewNameNG"
}`

var nodeGroupEndpoint = datastoreEndpoint + "/node_groups/" + ngID //nolint:gochecknoglobals

func newNodeGroupService(t *testing.T, serverURL string) *NodeGroupService {
	t.Helper()

	client, err := transport.NewHTTPClient(http.DefaultClient, "token", serverURL+"/v2")
	require.NoError(t, err)

	engine := internal.NewEngineService(
		client,
		common.EngineClickHouse,
	)

	return &NodeGroupService{
		EngineService: engine,
	}
}

func newNodeGroupServiceWithMockClient() *NodeGroupService {
	mockClient := mockClient{}

	engine := internal.NewEngineService(
		mockClient,
		common.EngineClickHouse,
	)

	return &NodeGroupService{
		EngineService: engine,
	}
}

func TestNodeGroupService_CreateNodeGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, datastoreEndpoint+"/node_groups", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req NodeGroupCreateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, 1, req.NodeCount)
		require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", req.Flavor.ID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Mock response from API
		_, err = w.Write([]byte(simpleNodeGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	req := NodeGroupCreateRequest{
		Name:      "shard1",
		Role:      "DATA",
		NodeCount: 1,
		Flavor: FlavorForNodeGroupRequest{
			Type: "FIXED",
			ID:   "550e8400-e29b-41d4-a716-446655440000",
			// API requires DiskType filed.
		},
	}

	result, err := srv.CreateNodeGroup(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, ngID, result.ID)
	require.Equal(t, "NewNameNG", result.Name)
}

func TestNodeGroupService_CreateNodeGroup_InvalidRequest(t *testing.T) {
	srv := newNodeGroupServiceWithMockClient()
	body := NodeGroupCreateRequest{
		Name:      "shard1",
		Role:      "DATA",
		NodeCount: 1,
		Flavor: FlavorForNodeGroupRequest{
			Type: "FLEXIBLE",
			ID:   "550e8400-e29b-41d4-a716-446655440000",
			// API requires DiskType filed.
		},
	}
	_, err := srv.CreateNodeGroup(context.Background(), dsID, body)
	require.Error(t, err)
	require.Equal(
		t, "validate body: node_group.flavor: flavor.id must not be specified for FLEXIBLE flavor", err.Error(),
	)
}

func TestNodeGroupService_DeleteNodeGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, nodeGroupEndpoint, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Empty(t, body)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	err := srv.DeleteNodeGroup(context.Background(), dsID, ngID)

	require.NoError(t, err)
}

func TestNodeGroupService_DeleteNodeGroup_InvalidRequest(t *testing.T) {
	srv := newNodeGroupServiceWithMockClient()

	err := srv.DeleteNodeGroup(context.Background(), dsID, "ngid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate node group id:")
}

func TestNodeGroupService_ResizeNodeGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, nodeGroupEndpoint+"/resize", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req NodeGroupResizeRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, 1, req.NodeCount)
		require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", req.Flavor.ID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleNodeGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	req := NodeGroupResizeRequest{
		NodeCount: 1,
		Flavor: FlavorForNodeGroupRequest{
			Type: "FIXED",
			ID:   "550e8400-e29b-41d4-a716-446655440000",
			// API requires DiskType filed.
		},
	}

	result, err := srv.ResizeNodeGroup(context.Background(), dsID, ngID, req)

	require.NoError(t, err)
	require.Equal(t, ngID, result.ID)
	require.Equal(t, "NewNameNG", result.Name)
}

func TestNodeGroupService_ResizeNodeGroup_InvalidRequest(t *testing.T) {
	srv := newNodeGroupServiceWithMockClient()
	body := NodeGroupResizeRequest{NodeCount: 2}
	_, err := srv.ResizeNodeGroup(context.Background(), dsID, ngID, body)
	require.Error(t, err)
	require.Equal(t, "validate body: validate flavor: unsupported flavor type: \"\"", err.Error())
}

func TestNodeGroupService_DeleteNodeGroupInstances_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, nodeGroupEndpoint+"/instances", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req NodeGroupDeleteInstancesRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Len(t, req.Instances, 2)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleNodeGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	req := NodeGroupDeleteInstancesRequest{
		Instances: []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"550e8400-e29b-41d4-a716-446655440001",
		},
	}

	result, err := srv.DeleteNodeGroupInstances(context.Background(), dsID, ngID, req)

	require.NoError(t, err)
	require.Equal(t, ngID, result.ID)
	require.Equal(t, "NewNameNG", result.Name)
}

func TestNodeGroupService_UpdateNodeGroupWeight_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, nodeGroupEndpoint+"/weight", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req NodeGroupUpdateWeightRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, 100, req.Weight)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleNodeGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	req := NodeGroupUpdateWeightRequest{Weight: 100}

	result, err := srv.UpdateNodeGroupWeight(context.Background(), dsID, ngID, req)

	require.NoError(t, err)
	require.Equal(t, ngID, result.ID)
	require.Equal(t, "NewNameNG", result.Name)
}

func TestNodeGroupService_UpdateNodeGroupShardGroups_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, nodeGroupEndpoint+"/shard-groups", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req NodeGroupUpdateShardGroupsRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Len(t, req.ShardGroups, 2)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleNodeGroupResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newNodeGroupService(t, server.URL)

	req := NodeGroupUpdateShardGroupsRequest{
		ShardGroups: []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"550e8400-e29b-41d4-a716-446655440001",
		},
	}

	result, err := srv.UpdateNodeGroupShardGroups(context.Background(), dsID, ngID, req)

	require.NoError(t, err)
	require.Equal(t, ngID, result.ID)
	require.Equal(t, "NewNameNG", result.Name)
}
