package clickhouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

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
	NodeGroups     []NodeGroupResponse    `json:"node_groups"`
}

// DatastoreListResponse is the API response for the clickhouse datastore list.
type DatastoreListResponse struct {
	Datastores []DatastoreResponse `json:"datastores"`
}

// DatastoreService is service to interact with clickhouse datastore resource.
type DatastoreService struct {
	*internal.EngineService
}

// DatastoreCreateRequest represents body for the datastore Create request.
type DatastoreCreateRequest struct {
	Config         map[string]any           `json:"config,omitempty"`
	LogPlatform    *DatastoreLogGroup       `json:"log_platform,omitempty"`
	Name           string                   `json:"name"`
	TypeID         string                   `json:"type_id"`
	SubnetID       string                   `json:"subnet_id"`
	Password       string                   `json:"password"`
	SecurityGroups []string                 `json:"security_groups,omitempty"`
	NodeGroups     []NodeGroupCreateRequest `json:"node_groups"`
}

func (r DatastoreCreateRequest) validate() error {
	if r.Name == "" {
		return errors.New("datastore.name is required") //nolint:goerr113 // Dynamic error
	}

	if r.Password == "" {
		return errors.New("password is required") //nolint:goerr113 // Dynamic error
	}

	if err := uuid.Validate(r.TypeID); err != nil {
		return fmt.Errorf("type_id must be a valid UUID: %w", err)
	}

	if err := uuid.Validate(r.SubnetID); err != nil {
		return fmt.Errorf("subnet_id must be a valid UUID: %w", err)
	}

	if len(r.NodeGroups) == 0 {
		return errors.New("at least one node_group is required") //nolint:goerr113 // Dynamic error
	}

	for i, group := range r.NodeGroups {
		if err := group.validate(); err != nil {
			return fmt.Errorf("groups[%d]: %w", i, err)
		}
	}

	return nil
}

// {
// 	"config": {
// 	  "simple_parameter": 1
// 	},
// 	"log_platform": {
// 	  "log_group": "s/dbaas/Logs"
// 	},
// 	"name": "Johnny_Cage",
// 	"node_groups": [
// 	  {
// 		"flavor": {
// 		  "disk": 256,
// 		  "disk_type": "LOCAL",
// 		  "ram": 4096,
// 		  "type": "FLEXIBLE",
// 		  "vcpus": 4
// 		},
// 		"name": "My pretty keeper group",
// 		"role": "KEEPER"
// 	  },
// 	  {
// 		"flavor": {
// 		  "disk_type": "LOCAL",
// 		  "id": "e62d8ce9-58ec-4330-912a-7aa4b62b7b35",
// 		  "type": "FIXED"
// 		},
// 		"has_public_ips": false,
// 		"name": "My pretty data group",
// 		"node_count": 1,
// 		"role": "DATA",
// 		"weight": 100
// 	  }
// 	],
// 	"password": "Adm!n123",
// 	"security_groups": [
// 	  "cdb6331d-1967-4434-9263-7be18941c698",
// 	  "dd1dce88-8662-4685-99b7-8c9b4fe14a7e"
// 	],
// 	"subnet_id": "4e5588d8-53aa-4057-8a3f-e6c25b06bfba",
// 	"type_id": "275293b2-7b82-49cc-8e56-0278562b08ac"
//   }

// GetDatastoreList returns datastore list from api.
func (s *DatastoreService) GetDatastoreList(ctx context.Context) (DatastoreListResponse, error) {
	response := DatastoreListResponse{}

	err := s.Get(ctx, "", &response)
	if err != nil {
		//nolint:wrapcheck
		return response, err
	}

	return response, nil
}

// GetDatastore returns a datastore based on the ID.
func (s *DatastoreService) GetDatastore(ctx context.Context, datastoreID string) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	err := s.Get(ctx, "/"+datastoreID, &response)
	if err != nil {
		//nolint:wrapcheck
		return response, err
	}

	return response, nil
}

// Datastore creates a new datastore.
func (s *DatastoreService) CreateDatastore(
	ctx context.Context, body DatastoreCreateRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate datastore: %w", err)
	}

	err := s.Post(ctx, "", body, &response)
	if err != nil {
		//nolint:wrapcheck
		return response, err
	}

	return response, nil
}
