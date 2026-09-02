package clickhouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/common"
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

func (r DatastoreResponse) GetState() string {
	return string(r.State)
}

func (r DatastoreResponse) GetStatus() string {
	return string(r.Status)
}

// DatastoreListResponse is the API response for the clickhouse datastore list.
type DatastoreListResponse struct {
	Datastores []DatastoreResponse `json:"datastores"`
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
		return fmt.Errorf("datastore.name: %w", common.ErrFieldRequired)
	}

	if r.Password == "" {
		return fmt.Errorf("datastore.password: %w", common.ErrFieldRequired)
	}

	if err := uuid.Validate(r.TypeID); err != nil {
		return fmt.Errorf("datastore.type_id must be a valid UUID: %w", err)
	}

	if err := uuid.Validate(r.SubnetID); err != nil {
		return fmt.Errorf("datastore.subnet_id must be a valid UUID: %w", err)
	}

	if len(r.NodeGroups) == 0 {
		return fmt.Errorf("datastore.node_groups: %w", common.ErrFieldEmptySlice)
	}

	for i, group := range r.NodeGroups {
		if err := group.validate(); err != nil {
			return fmt.Errorf("node_groups[%d]: %w", i, err)
		}
	}

	return nil
}

// DatastoreUpdateRequest represents options for the datastore Update request.
type DatastoreUpdateRequest struct {
	Name string `json:"name"`
}

func (r DatastoreUpdateRequest) validate() error {
	if r.Name == "" {
		return fmt.Errorf("name: %w", common.ErrFieldRequired)
	}
	return nil
}

// DatastoreUpdatePasswordRequest represents options for the datastore Update password request.
type DatastoreUpdatePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (r DatastoreUpdatePasswordRequest) validate() error {
	if r.NewPassword == "" {
		return fmt.Errorf("new_password: %w", common.ErrFieldRequired)
	}
	return nil
}

// DatastoreSecurityGroupsRequest represents update options for the Datastore security groups.
type DatastoreSecurityGroupsRequest struct {
	SecurityGroups []string `json:"security_groups"`
}

func (r DatastoreSecurityGroupsRequest) validate() error {
	// add check min length 1 ?

	for i, sgID := range r.SecurityGroups {
		if err := uuid.Validate(sgID); err != nil {
			return fmt.Errorf("security_groups[%d]: %w", i, err)
		}
	}
	return nil
}

// DatastoreLogPlatformRequest represents enable options to update a datastore log platform.
type DatastoreLogPlatformRequest struct {
	LogPlatform DatastoreLogGroup `json:"log_platform"`
}

func (r DatastoreLogPlatformRequest) validate() error {
	if r.LogPlatform.LogGroup == "" {
		return fmt.Errorf("log_platform.log_group: %w", common.ErrFieldRequired)
	}
	return nil
}

// DatastoreConfigRequest represents options for the datastore's configuration parameters Update request.
type DatastoreConfigRequest struct {
	Config map[string]any `json:"config"`
}

func (r DatastoreConfigRequest) validate() error {
	if r.Config == nil {
		return fmt.Errorf("config: %w", common.ErrFieldRequired)
	}
	return nil
}

// DatastoreService is service to interact with clickhouse datastore resource.
type DatastoreService struct {
	*common.EngineService
}

// GetDatastoreList returns datastore list from api.
func (s *DatastoreService) GetDatastoreList(ctx context.Context) (DatastoreListResponse, error) {
	response := DatastoreListResponse{}

	err := s.Get(ctx, "", &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// GetDatastore returns a datastore based on the ID.
func (s *DatastoreService) GetDatastore(ctx context.Context, datastoreID string) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	err := s.Get(ctx, s.DatastorePath(datastoreID), &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// Datastore creates a new datastore.
func (s *DatastoreService) CreateDatastore(
	ctx context.Context, body DatastoreCreateRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Post(ctx, "", body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// UpdateDatastore updates an existing datastore.
func (s *DatastoreService) UpdateDatastore(
	ctx context.Context, datastoreID string, body DatastoreUpdateRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Patch(ctx, s.DatastorePath(datastoreID), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// UpdateDatastorePassword updates a password of the existing datastore.
func (s *DatastoreService) UpdateDatastorePassword(
	ctx context.Context, datastoreID string, body DatastoreUpdatePasswordRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Patch(ctx, s.DatastorePath(datastoreID, "password"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// UpdateDatastoreSecurityGroups updates a security groups of the existing datastore.
func (s *DatastoreService) UpdateDatastoreSecurityGroups(
	ctx context.Context, datastoreID string, body DatastoreSecurityGroupsRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Put(ctx, s.DatastorePath(datastoreID, "security-groups"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// EnableLogPlatform updates a log platform params of the existing datastore.
func (s *DatastoreService) EnableLogPlatform(
	ctx context.Context, datastoreID string, body DatastoreLogPlatformRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Put(ctx, s.DatastorePath(datastoreID, "log-platform"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// DisableLogPlatform disables the log platform for the existing datastore.
func (s *DatastoreService) DisableLogPlatform(ctx context.Context, datastoreID string) error {
	if err := uuid.Validate(datastoreID); err != nil {
		return fmt.Errorf("validate datastore id: %w", err)
	}

	err := s.Delete(ctx, s.DatastorePath(datastoreID, "log-platform"))
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

// UpdateDatastoreConfig updates a config of the existing datastore.
func (s *DatastoreService) UpdateDatastoreConfig(
	ctx context.Context, datastoreID string, body DatastoreConfigRequest,
) (DatastoreResponse, error) {
	response := DatastoreResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Put(ctx, s.DatastorePath(datastoreID, "config"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

// DeleteDatastore deletes a datastore by ID.
func (s *DatastoreService) DeleteDatastore(ctx context.Context, datastoreID string) error {
	if err := uuid.Validate(datastoreID); err != nil {
		return fmt.Errorf("validate datastore id: %w", err)
	}

	err := s.Delete(ctx, s.DatastorePath(datastoreID))
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}
