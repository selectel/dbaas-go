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

func newDatastoreService(t *testing.T, serverURL string) *DatastoreService {
	t.Helper()

	client, err := transport.NewHTTPClient(
		http.DefaultClient,
		"token",
		serverURL+"/v2",
	)
	require.NoError(t, err)

	engine := internal.NewEngineService(
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

func newDastastoreServiceWithMockClient() *DatastoreService {
	mockClient := mockClient{}

	engine := internal.NewEngineService(
		mockClient,
		common.EngineClickHouse,
	)

	return &DatastoreService{
		EngineService: engine,
	}
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
		require.Empty(t, dataNG.Weight)
		require.Empty(t, dataNG.HasPublicIPs)

		keeperNG := req.NodeGroups[1]
		require.Equal(t, NodeGroupRoleKeeper, keeperNG.Role)
		require.Equal(t, common.FlavorTypeFlexible, keeperNG.Flavor.Type)
		require.Equal(t, 25, keeperNG.Flavor.Disk)
		require.Empty(t, dataNG.Weight)
		require.Empty(t, dataNG.HasPublicIPs)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Mock response from API
		_, err = w.Write([]byte(`{
			"id": "123",
			"name": "Test_cluster"
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

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
				Flavor: FlavorForNodeGroupCreate{
					Type: "FIXED",
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					// API requires DiskType filed.
				},
			},
			{
				Name:      "Keepers",
				Role:      "KEEPER",
				NodeCount: 3,
				Flavor: FlavorForNodeGroupCreate{
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
	require.Equal(t, "123", result.ID)
	require.Equal(t, "Test_cluster", result.Name)
}

func TestDatastoreService_CreateDatastore_InvalidRequest(t *testing.T) {
	srv := newDastastoreServiceWithMockClient()
	req := DatastoreCreateRequest{
		Name:     "Test_cluster",
		Password: "pass123",
		TypeID:   "000e0000-e29b-41d4-a716-446655000000",
		SubnetID: "00000000-e29b-41d4-a716-446655000000",
		NodeGroups: []NodeGroupCreateRequest{
			{
				Name: "shard1",
				Role: "DATA",
				Flavor: FlavorForNodeGroupCreate{
					Type: "FIXED",
					ID:   "550e8400-e29b-41d4-a716-446655440000",
				},
			},
		},
	}

	result, err := srv.CreateDatastore(context.Background(), req)
	require.Error(t, err)
	require.Equal(t, "validate body: node_groups[0]: node_group.node_count must be greater than 0", err.Error())
	require.Equal(t, DatastoreResponse{}, result)
}

func TestDatastoreCreateRequest_validate(t *testing.T) {
	tests := []struct {
		name      string
		errMsg    string
		datastore DatastoreCreateRequest
		wantErr   bool
	}{
		{
			name:      "datastore without name",
			datastore: DatastoreCreateRequest{},
			wantErr:   true,
			errMsg:    "datastore.name is required",
		},
		{
			name: "datastore without password",
			datastore: DatastoreCreateRequest{
				Name: "Test",
			},
			wantErr: true,
			errMsg:  "datastore.password is required",
		},
		{
			name: "datastore without type_id",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
			},
			wantErr: true,
			errMsg:  "datastore.type_id must be a valid UUID",
		},
		{
			name: "datastore without subnet_id",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
			errMsg:  "datastore.subnet_id must be a valid UUID",
		},
		{
			name: "datastore without node_groups",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
			},
			wantErr: true,
			errMsg:  "datastore.node_groups must be at least one",
		},
		{
			name: "datastore with invalid node_group",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
				NodeGroups: []NodeGroupCreateRequest{
					{Name: ""},
				},
			},
			wantErr: true,
			errMsg:  "node_groups[0]: node_group.name is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.datastore.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDatastoreService_UpdateDatastore_Success(t *testing.T) {
	dsID := "550e8400-e29b-41d4-a716-446655440000"
	url := "/v2/datastores/clickhouse/" + dsID
	newName := "NewNameDS"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, url, r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req DatastoreUpdateRequest
		require.NoError(t, json.Unmarshal(body, &req))

		require.Equal(t, newName, req.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Mock response from API
		_, err = w.Write([]byte(`{
			"id": "550e8400-e29b-41d4-a716-446655440000",
			"name": "NewNameDS"
		}`))
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

func TestDatastoreUpdateRequest_validate(t *testing.T) {
	tests := []struct {
		name      string
		errMsg    string
		datastore DatastoreUpdateRequest
		wantErr   bool
	}{
		{
			name:      "datastore without name",
			datastore: DatastoreUpdateRequest{},
			wantErr:   true,
			errMsg:    "datastore.name is required",
		},
		{
			name:      "datastore with name",
			datastore: DatastoreUpdateRequest{Name: "Test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.datastore.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
