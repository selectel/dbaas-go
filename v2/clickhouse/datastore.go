package clickhouse

import (
	"context"

	"github.com/selectel/dbaas-go/v2/internal"
)

// Datastore is the API response for the clickhouse datastores.
type DatastoreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DatastoreListResponse struct {
	Datastores []DatastoreResponse `json:"datastores"`
}

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
