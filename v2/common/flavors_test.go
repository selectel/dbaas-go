package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/selectel/dbaas-go/internal/transport"
)

func newFlavorService(t *testing.T, serverURL string) *FlavorService {
	t.Helper()

	client, err := transport.NewHTTPClient(http.DefaultClient, "token", serverURL+"/v2")
	require.NoError(t, err)

	return NewFlavorService(client)
}

const testFlavorListResponse = `{
    "errors": "",
    "flavors": [
		{
			"id": "4882f722-5d20-406a-a119-ca5cedaa4426",
			"disk_type": "LOCAL",
			"disk": 35,
			"fl_size": "STANDARD",
			"ram": 4096,
			"type": "FIXED",
			"vcpus": 2,
			"subtype": "STANDARD",
			"allowed_roles": [
				"DATA",
				"KEEPER"
			],
			"datastore_type_ids": [
				"00000000-0000-0000-0000-000000000001",
				"00000000-0000-0000-0000-000000000002"
			]
		},
		{
			"id": "ac3284cd-8730-4a51-9537-383d49556ced",
			"disk_type": "LOCAL",
			"disk": 32,
			"fl_size": "STANDARD",
			"ram": 8192,
			"type": "FIXED",
			"vcpus": 2,
			"subtype": "MEMORY",
			"allowed_roles": [
				"MANAGER",
				"DATA"
			],
			"datastore_type_ids": [
				"770379e9-3008-43bf-9e62-7fbb6bd4fa08"
			]
		},
		{
			"id": "b72e03b8-8aad-4157-8dd0-bf4dfe5da814",
			"name": "2-8192-64-redis-high-freq",
			"description": null,
			"vcpus": 2,
			"ram": 8192,
			"disk": 64,
			"fl_size": "high_freq",
			"subtype": null,
			"datastore_type_ids": [
				"2c8d1fff-20b2-4fe0-8055-df9ee127929f"
			],
			"host": null
		},
		{
		    "id": "04cc38a1-fc84-49ef-b8e8-7fb481ec8def",
            "name": "micro-dedicated-host",
            "description": null,
            "vcpus": 14,
            "ram": 32768,
            "disk": 512,
            "fl_size": "dedicated_host",
            "subtype": null,
            "datastore_type_ids": [
                "dab55ac0-6d85-407c-b2b0-09f67d0d0a16"
            ],
            "host": {
                "line": "Micro",
                "processor": "Intel E-2488",
                "available_count": 2
            }
        }
    ]
}`

var flavorListExpected FlavorListResponse = FlavorListResponse{ //nolint:gochecknoglobals
	Flavors: []FlavorResponse{
		// For clickhouse
		{
			ID:           "4882f722-5d20-406a-a119-ca5cedaa4426",
			DiskType:     "LOCAL",
			Disk:         35,
			FlSize:       "STANDARD",
			RAM:          4096,
			Type:         "FIXED",
			VCPUs:        2,
			AllowedRoles: []string{"DATA", "KEEPER"},
			DatastoreTypeIDs: []string{
				"00000000-0000-0000-0000-000000000001",
				"00000000-0000-0000-0000-000000000002",
			},
			// subtype: "STANDARD",
		},
		// For opensearch
		{
			ID:               "ac3284cd-8730-4a51-9537-383d49556ced",
			DiskType:         "LOCAL",
			Disk:             32,
			FlSize:           "STANDARD",
			RAM:              8192,
			Type:             "FIXED",
			VCPUs:            2,
			AllowedRoles:     []string{"MANAGER", "DATA"},
			DatastoreTypeIDs: []string{"770379e9-3008-43bf-9e62-7fbb6bd4fa08"},
			// "subtype": "MEMORY",
		},
		// For Evort
		{
			ID: "b72e03b8-8aad-4157-8dd0-bf4dfe5da814",
			// Name:             "2-8192-64-redis-high-freq",
			// Description:      nil,
			VCPUs:            2,
			RAM:              8192,
			Disk:             64,
			FlSize:           "high_freq",
			DatastoreTypeIDs: []string{"2c8d1fff-20b2-4fe0-8055-df9ee127929f"},
			Host:             nil,
			// "subtype": null,
		},
		{
			ID: "04cc38a1-fc84-49ef-b8e8-7fb481ec8def",
			// Name:        "micro-dedicated-host",
			// Description: nil,
			VCPUs:  14,
			RAM:    32768,
			Disk:   512,
			FlSize: "dedicated_host",
			// "subtype": null,
			DatastoreTypeIDs: []string{"dab55ac0-6d85-407c-b2b0-09f67d0d0a16"},
			Host: &FlavorHost{
				Line:           "Micro",
				Processor:      "Intel E-2488",
				AvailableCount: 2,
			},
		},
	},
}

func TestFlavorService_GetFlavorList_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/flavors", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testFlavorListResponse))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newFlavorService(t, server.URL)

	resp, err := srv.GetFlavorList(context.Background())

	require.NoError(t, err)

	require.Equal(t, flavorListExpected, resp)
}
