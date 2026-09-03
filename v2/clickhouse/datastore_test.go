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

const dsID = "550e8400-e29b-41d4-a716-446655440000"

var datastoreEndpoint = "/v2/datastores/clickhouse/" + dsID //nolint:gochecknoglobals

//nolint:gochecknoglobals
var simpleDatastoreResponse = `{
	"id": "550e8400-e29b-41d4-a716-446655440000",
	"name": "NewNameDS"
}`

func newDatastoreService(t *testing.T, serverURL string) *DatastoreService {
	t.Helper()

	client, err := transport.NewHTTPClient(http.DefaultClient, "token", serverURL+"/v2")
	require.NoError(t, err)

	engine := common.NewEngineService(
		client,
		common.EngineClickHouse,
	)

	return &DatastoreService{
		EngineService: engine,
	}
}

type mockClient struct{}

func (m mockClient) Do(ctx context.Context, method, path string, body, result any) error {
	return nil
}

func newDatastoreServiceWithMockClient() *DatastoreService {
	mockClient := mockClient{}

	engine := common.NewEngineService(
		mockClient,
		common.EngineClickHouse,
	)

	return &DatastoreService{
		EngineService: engine,
	}
}

const testDatastoreConfigurationParameters = `{
	"configuration-parameters": [
		{
			"id": "0fa4bb95-6f55-4249-94d0-a70cb85dc09b",
			"datastore_type_id": "00000000-0000-0000-0000-000000000001",
			"name": "server_settings.async_insert_threads",
			"type": "int",
			"choices": null,
			"min": "0",
			"max": "18446744073709551615",
			"default_value": 16,
			"invalid_values": null,
			"is_restart_required": true,
			"can_be_empty": false,
			"is_multiple_choice_available": false,
			"is_changeable": true
		},
		{
			"id": "084d0355-2d1a-47cb-83de-54bebd42e100",
			"datastore_type_id": "00000000-0000-0000-0000-000000000001",
			"name": "server_settings.background_merges_mutations_concurrency_ratio",
			"type": "float",
			"choices": null,
			"min": "0",
			"max": "3.4028235e+38F",
			"default_value": 2.0,
			"invalid_values": null,
			"is_restart_required": true,
			"can_be_empty": false,
			"is_multiple_choice_available": false,
			"is_changeable": true
		},
		{
			"id": "eee271c7-0afe-44d3-aa84-a4f56d9cf7f4",
			"datastore_type_id": "00000000-0000-0000-0000-000000000001",
			"name": "merge_tree_settings.deduplicate_merge_projection_mode",
			"type": "str",
			"choices": [
				"ignore",
				"throw",
				"drop",
				"rebuild"
			],
			"min": null,
			"max": null,
			"default_value": "throw",
			"invalid_values": null,
			"is_restart_required": true,
			"can_be_empty": false,
			"is_multiple_choice_available": false,
			"is_changeable": true
		},
		{
			"id": "11cc979c-b912-4efa-8eee-ee63bbc10283",
			"datastore_type_id": "00000000-0000-0000-0000-000000000001",
			"name": "server_settings.dictionaries_lazy_load",
			"type": "bool",
			"choices": null,
			"min": null,
			"max": null,
			"default_value": true,
			"invalid_values": null,
			"is_restart_required": true,
			"can_be_empty": false,
			"is_multiple_choice_available": false,
			"is_changeable": true
		}
	]
}`

func ptr[T any](v T) *T { return &v }

var DatastoreConfigurationParametersExpected = []DatastoreConfigurationParameterResponse{ //nolint:gochecknoglobals
	{
		ID:                        "0fa4bb95-6f55-4249-94d0-a70cb85dc09b",
		DatastoreTypeID:           "00000000-0000-0000-0000-000000000001",
		Name:                      "server_settings.async_insert_threads",
		Type:                      "int",
		Choices:                   nil,
		MinValue:                  ptr("0"),
		MaxValue:                  ptr("18446744073709551615"),
		DefaultValue:              float64(16),
		InvalidValues:             nil,
		IsRestartRequired:         true,
		CanBeEmpty:                false,
		IsMultipleChoiceAvailable: false,
		IsChangeable:              true,
	},
	{
		ID:                        "084d0355-2d1a-47cb-83de-54bebd42e100",
		DatastoreTypeID:           "00000000-0000-0000-0000-000000000001",
		Name:                      "server_settings.background_merges_mutations_concurrency_ratio",
		Type:                      "float",
		Choices:                   nil,
		MinValue:                  ptr("0"),
		MaxValue:                  ptr("3.4028235e+38F"),
		DefaultValue:              float64(2),
		InvalidValues:             nil,
		IsRestartRequired:         true,
		CanBeEmpty:                false,
		IsMultipleChoiceAvailable: false,
		IsChangeable:              true,
	},
	{
		ID:                        "eee271c7-0afe-44d3-aa84-a4f56d9cf7f4",
		DatastoreTypeID:           "00000000-0000-0000-0000-000000000001",
		Name:                      "merge_tree_settings.deduplicate_merge_projection_mode",
		Type:                      "str",
		Choices:                   []string{"ignore", "throw", "drop", "rebuild"},
		MinValue:                  nil,
		MaxValue:                  nil,
		DefaultValue:              "throw",
		InvalidValues:             nil,
		IsRestartRequired:         true,
		CanBeEmpty:                false,
		IsMultipleChoiceAvailable: false,
		IsChangeable:              true,
	},
	{
		ID:                        "11cc979c-b912-4efa-8eee-ee63bbc10283",
		DatastoreTypeID:           "00000000-0000-0000-0000-000000000001",
		Name:                      "server_settings.dictionaries_lazy_load",
		Type:                      "bool",
		Choices:                   nil,
		MinValue:                  nil,
		MaxValue:                  nil,
		DefaultValue:              true,
		InvalidValues:             nil,
		IsRestartRequired:         true,
		CanBeEmpty:                false,
		IsMultipleChoiceAvailable: false,
		IsChangeable:              true,
	},
}

const testDatastoresResponse = `{
	"datastores": [
		{
			"id": "04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1",
			"name": "Sinead",
			"created_at": "2026-08-07T13:12:45.295388",
			"updated_at": "2026-08-12T13:34:13.142778",
			"node_groups": [
				{
					"id": "08c2341d-7af7-4914-aee6-7b54156064eb",
					"name": "shard1",
					"has_public_ips": false,
					"role": "DATA",
					"node_count": 1,
					"flavor": {
						"id": "4882f722-5d20-406a-a119-ca5cedaa4426",
						"disk_type": "LOCAL",
						"disk": 35,
						"fl_size": "STANDARD",
						"ram": 4096,
						"type": "FIXED",
						"vcpus": 2,
						"subtype": "STANDARD"
					},
					"instances": [
						{
							"id": "f812b882-74ae-4b0b-9901-14f603d8221a",
							"ip": "192.168.2.78",
							"flavor_id": "4882f722-5d20-406a-a119-ca5cedaa4426",
							"floating_ip": null,
							"status": "ACTIVE",
							"availability_zone": "ru-2a",
							"hostname": "f812b882-74ae-4b0b-9901-14f603d8221a.ru-2.c.dbaas.selcloud.org"
						}
					],
					"status": "ACTIVE",
					"weight": 100,
					"shard_groups": [
						"c377547d-68e9-42ab-942f-a4e3305dea0f"
					]
				}
			],
			"state": "RUNNING",
			"status": "ACTIVE",
			"project_id": "295b45d99b6f4afea59ad504ac73e5f6",
			"subnet_id": "7160539b-3cf6-4a65-a75b-24071390d80f",
			"type_id": "00000000-0000-0000-0000-000000000002",
			"security_groups": [
				"9c922408-51de-4598-9894-e69ca00e51c8"
			],
			"log_platform": {
				"log_group": null
			},
			"config": {}
		},
		{
			"id": "ds-2", 
			"name": "clickhouse-2"
		}
	]
}`

var datastoreListExpected DatastoreListResponse = DatastoreListResponse{ //nolint:gochecknoglobals
	Datastores: []DatastoreResponse{
		{
			ID:        "04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1",
			Name:      "Sinead",
			CreatedAt: "2026-08-07T13:12:45.295388",
			UpdatedAt: "2026-08-12T13:34:13.142778",
			State:     common.DatastoreStateRunning,
			Status:    common.DatastoreStatusActive,
			ProjectID: "295b45d99b6f4afea59ad504ac73e5f6",
			SubnetID:  "7160539b-3cf6-4a65-a75b-24071390d80f",
			TypeID:    "00000000-0000-0000-0000-000000000002",
			SecurityGroups: []string{
				"9c922408-51de-4598-9894-e69ca00e51c8",
			},
			LogPlatform: DatastoreLogGroup{},
			Config:      map[string]any{},
			NodeGroups: []NodeGroupResponse{
				{
					ID:           "08c2341d-7af7-4914-aee6-7b54156064eb",
					Name:         "shard1",
					HasPublicIPs: false,
					Role:         NodeGroupRoleData,
					NodeCount:    1,
					Status:       NodeGroupStatusActive,
					Weight:       100,
					ShardGroups: []string{
						"c377547d-68e9-42ab-942f-a4e3305dea0f",
					},
					Flavor: FlavorResponse{
						ID:       "4882f722-5d20-406a-a119-ca5cedaa4426",
						DiskType: common.FlavorDiskLocal,
						Disk:     35,
						FlSize:   "STANDARD",
						RAM:      4096,
						VCPUs:    2,
						Type:     common.FlavorTypeFIXED,
					},
					Instances: []InstanceResponse{
						{
							ID:               "f812b882-74ae-4b0b-9901-14f603d8221a",
							IP:               "192.168.2.78",
							FlavorID:         "4882f722-5d20-406a-a119-ca5cedaa4426",
							FloatingIP:       "",
							Status:           InstanceStatusActive,
							AvailabilityZone: "ru-2a",
							Hostname:         "f812b882-74ae-4b0b-9901-14f603d8221a.ru-2.c.dbaas.selcloud.org",
						},
					},
				},
			},
		},
		{
			ID:   "ds-2",
			Name: "clickhouse-2",
		},
	},
}

const testDatastoreResponse = `{
	"id": "04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1",
	"name": "Sinead",
	"created_at": "2026-08-07T13:12:45.295388",
	"updated_at": "2026-08-12T13:34:13.142778",
	"node_groups": [
		{
			"id": "08c2341d-7af7-4914-aee6-7b54156064eb",
			"name": "shard1",
			"has_public_ips": false,
			"role": "DATA",
			"node_count": 1,
			"flavor": {
				"id": "4882f722-5d20-406a-a119-ca5cedaa4426",
				"disk_type": "LOCAL",
				"disk": 35,
				"fl_size": "STANDARD",
				"ram": 4096,
				"type": "FIXED",
				"vcpus": 2,
				"subtype": "STANDARD"
			},
			"instances": [
				{
					"id": "f812b882-74ae-4b0b-9901-14f603d8221a",
					"ip": "192.168.2.78",
					"flavor_id": "4882f722-5d20-406a-a119-ca5cedaa4426",
					"floating_ip": null,
					"status": "ACTIVE",
					"availability_zone": "ru-2a",
					"hostname": "f812b882-74ae-4b0b-9901-14f603d8221a.ru-2.c.dbaas.selcloud.org"
				}
			],
			"status": "ACTIVE",
			"weight": 100,
			"shard_groups": [
				"c377547d-68e9-42ab-942f-a4e3305dea0f"
			]
		}
	],
	"state": "RUNNING",
	"status": "ACTIVE",
	"project_id": "295b45d99b6f4afea59ad504ac73e5f6",
	"subnet_id": "7160539b-3cf6-4a65-a75b-24071390d80f",
	"type_id": "00000000-0000-0000-0000-000000000002",
	"security_groups": [
		"9c922408-51de-4598-9894-e69ca00e51c8"
	],
	"log_platform": {
		"log_group": null
	},
	"config": {}
}`

func TestDatastoreService_GetDatastoreList_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testDatastoresResponse))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	resp, err := srv.GetDatastoreList(context.Background())

	require.NoError(t, err)

	require.Equal(t, datastoreListExpected, resp)
}

func TestDatastoreService_GetDatastoreConfigurationParameters_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse/configuration-parameters", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testDatastoreConfigurationParameters))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	resp, err := srv.GetDatastoreConfigurationParameters(context.Background())

	require.NoError(t, err)

	require.Equal(t, DatastoreConfigurationParametersExpected, resp)
}

func TestDatastoreService_GetDatastore_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse/04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testDatastoreResponse))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	resp, err := srv.GetDatastore(context.Background(), "04c5bf32-bb4b-40b4-b0c7-3cb4eb1d99d1")

	require.NoError(t, err)

	require.Equal(t, datastoreListExpected.Datastores[0], resp)
}

func TestDatastoreService_CreateDatastore_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreCreateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Len(t, req.NodeGroups, 2)
		require.Empty(t, req.SecurityGroups)
		require.Empty(t, req.LogPlatform)
		require.Empty(t, req.Config)

		dataNG := req.NodeGroups[0]
		require.Equal(t, NodeGroupRoleData, dataNG.Role)
		require.Equal(t, common.FlavorTypeFIXED, dataNG.Flavor.Type)
		require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", dataNG.Flavor.ID)
		require.Empty(t, dataNG.Flavor.Disk)
		require.Empty(t, dataNG.Flavor.RAM)
		require.Empty(t, dataNG.Flavor.VCPUs)
		require.Equal(t, 100, *dataNG.Weight)
		require.Empty(t, dataNG.HasPublicIPs)

		keeperNG := req.NodeGroups[1]
		require.Equal(t, NodeGroupRoleKeeper, keeperNG.Role)
		require.Equal(t, common.FlavorTypeFlexible, keeperNG.Flavor.Type)
		require.Equal(t, 25, keeperNG.Flavor.Disk)
		require.Empty(t, keeperNG.Weight)
		require.Empty(t, keeperNG.HasPublicIPs)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)
	weight := 100

	req := DatastoreCreateRequest{
		Name:     "Test_cluster",
		Password: "pass123",
		TypeID:   "000e0000-e29b-41d4-a716-446655000000",
		SubnetID: "00000000-e29b-41d4-a716-446655000000",
		NodeGroups: []NodeGroupCreateRequest{
			{
				Name:      "shard1",
				Role:      "DATA",
				NodeCount: 1,
				Flavor: FlavorForNodeGroupRequest{
					Type: "FIXED",
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					// API requires DiskType field.
				},
				Weight: &weight,
			},
			{
				Name:      "Keepers",
				Role:      "KEEPER",
				NodeCount: 3,
				Flavor: FlavorForNodeGroupRequest{
					Type:     "FLEXIBLE",
					DiskType: "NETWORK-ULTRA",
					RAM:      4096,
					VCPUs:    2,
					Disk:     25,
				},
			},
		},
	}

	result, err := srv.CreateDatastore(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
	require.Equal(t, "NewNameDS", result.Name)
}

func TestDatastoreService_CreateDatastore_InvalidRequest(t *testing.T) {
	srv := newDatastoreServiceWithMockClient()
	req := DatastoreCreateRequest{
		Name:     "Test_cluster",
		Password: "pass123",
		TypeID:   "000e0000-e29b-41d4-a716-446655000000",
		SubnetID: "00000000-e29b-41d4-a716-446655000000",
		NodeGroups: []NodeGroupCreateRequest{
			{
				Name: "shard1",
				Role: "DATA",
				Flavor: FlavorForNodeGroupRequest{
					Type: "FIXED",
					ID:   "550e8400-e29b-41d4-a716-446655440000",
				},
			},
		},
	}

	result, err := srv.CreateDatastore(context.Background(), req)
	require.Error(t, err)
	require.Equal(t, "validate body: node_groups[0]: node_group.node_count: value must be greater than 0", err.Error())
	require.Equal(t, DatastoreResponse{}, result)
}

func TestDatastoreService_UpdateDatastore_Success(t *testing.T) {
	newName := "NewNameDS"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, datastoreEndpoint, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreUpdateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, newName, req.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	req := DatastoreUpdateRequest{
		Name: newName,
	}

	result, err := srv.UpdateDatastore(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
	require.Equal(t, newName, result.Name)
}

func TestDatastoreService_UpdateDatastorePassword_Success(t *testing.T) {
	url := datastoreEndpoint + "/password"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, url, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreUpdatePasswordRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, "secret", req.NewPassword)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	req := DatastoreUpdatePasswordRequest{
		NewPassword: "secret",
	}

	result, err := srv.UpdateDatastorePassword(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
}

func TestDatastoreService_UpdateDatastoreSecurityGroups_Success(t *testing.T) {
	url := datastoreEndpoint + "/security-groups"
	sGroups := []string{
		"550e8400-e29b-41d4-a716-446655440002",
		"550e8400-e29b-41d4-a716-446655440001",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, url, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreSecurityGroupsRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, sGroups, req.SecurityGroups)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	req := DatastoreSecurityGroupsRequest{
		SecurityGroups: sGroups,
	}

	result, err := srv.UpdateDatastoreSecurityGroups(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
}

func TestDatastoreService_EnableLogPlatform_Success(t *testing.T) {
	url := datastoreEndpoint + "/log-platform"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, url, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreLogPlatformRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, "TestGroup", req.LogPlatform.LogGroup)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	req := DatastoreLogPlatformRequest{
		LogPlatform: DatastoreLogGroup{
			LogGroup: "TestGroup",
		},
	}

	result, err := srv.EnableLogPlatform(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
}

func TestDatastoreService_DisableLogPlatform_Success(t *testing.T) {
	url := datastoreEndpoint + "/log-platform"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, url, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	err := srv.DisableLogPlatform(context.Background(), dsID)

	require.NoError(t, err)
}

func TestDatastoreService_UpdateDatastoreConfig_Success(t *testing.T) {
	url := datastoreEndpoint + "/config"
	config := map[string]any{
		"param_first":  1,
		"param_second": "test",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, url, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreConfigRequest
		require.NoError(t, json.Unmarshal(body, &req))
		require.Equal(t, float64(1), req.Config["param_first"])
		require.Equal(t, "test", req.Config["param_second"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(simpleDatastoreResponse))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	req := DatastoreConfigRequest{
		Config: config,
	}

	result, err := srv.UpdateDatastoreConfig(context.Background(), dsID, req)

	require.NoError(t, err)
	require.Equal(t, dsID, result.ID)
}

func TestDatastoreService_DeleteDatastore_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, datastoreEndpoint, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Empty(t, body)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	err := srv.DeleteDatastore(context.Background(), dsID)

	require.NoError(t, err)
}
