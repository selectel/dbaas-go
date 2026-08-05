package clickhouse

import (
	"context"
	"encoding/json"
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

	client := transport.NewHTTPClient(
		http.DefaultClient,
		"token",
		serverURL+"/v2",
		"test-agent",
	)

	engine := internal.NewEngineService(
		client,
		common.EngineClickHouse,
	)

	return &DatastoreService{
		EngineService: engine,
	}
}

func TestDatastoreService_List_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastores/clickhouse", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		require.NoError(t, json.NewEncoder(w).Encode(DatastoreListResponse{
			Datastores: []DatastoreResponse{
				{
					ID:   "1",
					Name: "cluster-1",
				},
				{
					ID:   "2",
					Name: "cluster-2",
				},
			},
		}))
	}))
	defer server.Close()

	srv := newDatastoreService(t, server.URL)

	resp, err := srv.List(context.Background())

	require.NoError(t, err)
	require.Len(t, resp.Datastores, 2)
}
