package common

import (
	"context"

	"github.com/selectel/dbaas-go/internal/transport"
)

// DatastoreTypeResponse is the API response for common datastore type.
type DatastoreTypeResponse struct {
	ID      string `json:"id"`
	Engine  string `json:"engine"`
	Version string `json:"version"`
}

// DatastoreTypeListResponse is the API response for datastore type list.
type DatastoreTypeListResponse struct {
	Errors         string                  `json:"errors"`
	DatastoreTypes []DatastoreTypeResponse `json:"datastore-types"`
}

type DatastoreTypeService struct {
	*baseService
}

// GetDatastoreTypeList returns datastore type list from api.
func (s *DatastoreTypeService) GetDatastoreTypeList(ctx context.Context) (DatastoreTypeListResponse, error) {
	response := DatastoreTypeListResponse{}

	err := s.Get(ctx, "", &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func NewDatastoreTypeService(client transport.Client) *DatastoreTypeService {
	root := "/datastore-types"

	return &DatastoreTypeService{
		baseService: &baseService{
			client:   client,
			rootPath: root,
		},
	}
}
