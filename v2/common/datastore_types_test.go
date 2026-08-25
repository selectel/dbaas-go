package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/selectel/dbaas-go/internal/transport"
)

func newDatastoreTypeService(t *testing.T, serverURL string) *DatastoreTypeService {
	t.Helper()

	client, err := transport.NewHTTPClient(http.DefaultClient, "token", serverURL+"/v2")
	require.NoError(t, err)

	return NewDatastoreTypeService(client)
}

const testDatastoreTypesResponse = `{
    "errors": "",
    "datastore-types": [
        {
            "id": "60de5498-37b5-49b9-9f7d-47c5e23b6aa1",
            "engine": "postgresql",
            "version": "18"
        },
        {
            "id": "770379e9-3008-43bf-9e62-7fbb6bd4fa08",
            "engine": "opensearch",
            "version": "2.18"
        },
        {
            "id": "00000000-0000-0000-0000-000000000002",
            "engine": "clickhouse",
            "version": "26.3.12.3"
        }
    ]
}`

var datastoreTypeListExpected DatastoreTypeListResponse = DatastoreTypeListResponse{ //nolint:gochecknoglobals
	DatastoreTypes: []DatastoreTypeResponse{
		{
			ID:      "60de5498-37b5-49b9-9f7d-47c5e23b6aa1",
			Engine:  "postgresql",
			Version: "18",
		},
		{
			ID:      "770379e9-3008-43bf-9e62-7fbb6bd4fa08",
			Engine:  "opensearch",
			Version: "2.18",
		},
		{
			ID:      "00000000-0000-0000-0000-000000000002",
			Engine:  "clickhouse",
			Version: "26.3.12.3",
		},
	},
}

func TestDatastoreTypeService_GetDatastoreList_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/datastore-types", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testDatastoreTypesResponse))

		require.NoError(t, err)
	}))
	defer server.Close()

	srv := newDatastoreTypeService(t, server.URL)

	resp, err := srv.GetDatastoreTypeList(context.Background())

	require.NoError(t, err)

	require.Equal(t, datastoreTypeListExpected, resp)
}
