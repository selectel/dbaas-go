package clickhouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/internal"
)

// NodeGroupStatus represents custom type for DBaaS v2 clickhouse node group statuses.
type NodeGroupStatus string

const (
	NodeGroupStatusActive   NodeGroupStatus = "ACTIVE"
	NodeGroupStatusCreating NodeGroupStatus = "CREATING"
	NodeGroupStatusDeleted  NodeGroupStatus = "DELETED"
	NodeGroupStatusDeleting NodeGroupStatus = "DELETING"
	NodeGroupStatusError    NodeGroupStatus = "ERROR"
	NodeGroupStatusResizing NodeGroupStatus = "RESIZING"
	NodeGroupStatusUpdating NodeGroupStatus = "UPDATING"
)

// NodeGroupRole.
type NodeGroupRole string

const (
	NodeGroupRoleData   NodeGroupRole = "DATA"
	NodeGroupRoleKeeper NodeGroupRole = "KEEPER"
)

// NodeGroupResponse is the API response for the clickhouse node group.
type NodeGroupResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Role         NodeGroupRole      `json:"role"`
	Status       NodeGroupStatus    `json:"status"`
	ShardGroups  []string           `json:"shard_groups"`
	Instances    []InstanceResponse `json:"instances"`
	Flavor       FlavorResponse     `json:"flavor"`
	NodeCount    int                `json:"node_count"`
	Weight       int                `json:"weight"`
	HasPublicIPs bool               `json:"has_public_ips"`
}

// NodeGroupResponse is the API response for the clickhouse node group.
type NodeGroupCreateRequest struct {
	Weight       *int                      `json:"weight,omitempty"`
	HasPublicIPs *bool                     `json:"has_public_ips,omitempty"`
	Name         string                    `json:"name"`
	Role         NodeGroupRole             `json:"role"`
	Flavor       FlavorForNodeGroupRequest `json:"flavor"`
	NodeCount    int                       `json:"node_count"`
}

func (n NodeGroupCreateRequest) validate() error {
	if n.Name == "" {
		return errors.New("node_group.name is required") //nolint:goerr113 // Dynamic error
	}

	if n.Role != NodeGroupRoleData && n.Role != NodeGroupRoleKeeper {
		return errors.New("node_group.role must be DATA or KEEPER") //nolint:goerr113 // Dynamic error
	}

	if err := n.Flavor.validate(); err != nil {
		return fmt.Errorf("node_group.flavor: %w", err)
	}

	if n.NodeCount <= 0 {
		return errors.New("node_group.node_count must be greater than 0") //nolint:goerr113 // Dynamic error
	}

	if n.Weight != nil && n.Role == NodeGroupRoleKeeper {
		return errors.New("node_group.role KEEPER could not have weight") //nolint:goerr113 // Dynamic error
	}
	return nil
}

type NodeGroupResizeRequest struct {
	Flavor    FlavorForNodeGroupRequest `json:"flavor"`
	NodeCount int                       `json:"node_count"`
}

func (r NodeGroupResizeRequest) validate() error {
	if err := r.Flavor.validate(); err != nil {
		return fmt.Errorf("validate flavor: %w", err)
	}

	if r.NodeCount <= 0 {
		return errors.New("node_count must be greater than 0") //nolint:goerr113 // Dynamic error
	}
	return nil
}

// DatastoreService is service to interact with clickhouse datastore resource.
type NodegroupService struct {
	*internal.EngineService
}

func (s *NodegroupService) nodeGroupsPath(datastoreID string, parts ...string) string {
	path := []string{"node_groups"}
	path = append(path, parts...)

	return s.DatastorePath(datastoreID, path...)
}

func (s *NodegroupService) CreateNodeGroup(
	ctx context.Context, datastoreID string, body NodeGroupCreateRequest,
) (NodeGroupResponse, error) {
	response := NodeGroupResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}
	err := s.Post(ctx, s.nodeGroupsPath(datastoreID), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

func (s *NodegroupService) DeleteNodeGroup(ctx context.Context, datastoreID, nodeGroupID string) error {
	if err := uuid.Validate(datastoreID); err != nil {
		return fmt.Errorf("validate datastore id: %w", err)
	}

	if err := uuid.Validate(nodeGroupID); err != nil {
		return fmt.Errorf("validate node group id: %w", err)
	}

	err := s.Delete(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID))
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func (s *NodegroupService) ResizeNodeGroup(
	ctx context.Context, datastoreID, nodeGroupID string, body NodeGroupResizeRequest,
) (NodeGroupResponse, error) {
	response := NodeGroupResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := uuid.Validate(nodeGroupID); err != nil {
		return response, fmt.Errorf("validate node group id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Patch(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID, "resize"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}
