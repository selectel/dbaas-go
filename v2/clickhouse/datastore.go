package clickhouse

import (
	"context"

	"github.com/selectel/dbaas-go/v2/common"
	"github.com/selectel/dbaas-go/v2/internal"
)

// DatastoreLogGroup represents log platform schema.
type DatastoreLogGroup struct {
	LogGroup string `json:"log_group"`
}

// DatastoreResponse is the API response for the clickhouse datastore.
type DatastoreResponse struct {
	Config         map[string]any         `json:"config"`
	ID             string                 `json:"id"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
	ProjectID      string                 `json:"project_id"`
	Name           string                 `json:"name"`
	TypeID         string                 `json:"type_id"`
	SubnetID       string                 `json:"subnet_id"`
	Status         common.DatastoreStatus `json:"status"`
	State          common.DatastoreState  `json:"state"`
	SecurityGroups []string               `json:"security_groups"`
	LogPlatform    DatastoreLogGroup      `json:"log_platform"`
	NodeGroups     []NodeGroup            `json:"node_groups"`
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
