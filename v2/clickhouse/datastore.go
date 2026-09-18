package clickhouse

import (
	"context"

	"github.com/selectel/dbaas-go/v2/internal"
)

// DatastoreResponse is the API response for the clickhouse datastore.
type DatastoreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DatastoreListResponse is the API response for the clickhouse datastore list.
type DatastoreListResponse struct {
	Datastores []DatastoreResponse `json:"datastores"`
}

// DatastoreService is service to interact with clickhouse datastore resource.
type DatastoreService struct {
	*internal.EngineService
}

func (s *DatastoreService) List(ctx context.Context) (DatastoreListResponse, error) {
	response := DatastoreListResponse{}

	err := s.Get(ctx, "", &response)
	if err != nil {
		//nolint:wrapcheck
		return response, err
	}

	return response, nil
}
