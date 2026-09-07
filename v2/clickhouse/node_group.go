package clickhouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/common"
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

// NodeGroupCreateRequest is the request body for creating a clickhouse node group.
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
		return fmt.Errorf("node_group.name: %w", common.ErrFieldRequired)
	}

	if n.Role != NodeGroupRoleData && n.Role != NodeGroupRoleKeeper {
		return fmt.Errorf("node_group.role: %w", common.ErrUnsupportedNodeGroupRole)
	}

	if err := n.Flavor.validate(); err != nil {
		return fmt.Errorf("node_group.flavor: %w", err)
	}

	if n.NodeCount <= 0 {
		return fmt.Errorf("node_group.node_count: %w", common.ErrPositiveIntegerRequired)
	}

	if n.Weight != nil && n.Role == NodeGroupRoleKeeper {
		return errors.New("node_group.role KEEPER could not have weight") //nolint:goerr113 // Dynamic error
	}

	if n.Role == NodeGroupRoleData && n.Weight == nil {
		return errors.New( //nolint:goerr113 // Dynamic error
			"node_group.weight: required for node_group.role DATA")
	}
	return nil
}

// NodeGroupResizeRequest is the request body to resize a node group.
type NodeGroupResizeRequest struct {
	Flavor    FlavorForNodeGroupRequest `json:"flavor"`
	NodeCount int                       `json:"node_count"`
}

func (r NodeGroupResizeRequest) validate() error {
	if err := r.Flavor.validate(); err != nil {
		return fmt.Errorf("validate flavor: %w", err)
	}

	if r.NodeCount <= 0 {
		return fmt.Errorf("node_count: %w", common.ErrPositiveIntegerRequired)
	}
	return nil
}

// NodeGroupDeleteInstancesRequest is the request body to resize a node group by reducing instances.
type NodeGroupDeleteInstancesRequest struct {
	Instances []string `json:"instances"`
}

func (r NodeGroupDeleteInstancesRequest) validate() error {
	if len(r.Instances) == 0 {
		return fmt.Errorf("instances: %w", common.ErrFieldEmptySlice)
	}
	for i, ID := range r.Instances {
		if err := uuid.Validate(ID); err != nil {
			return fmt.Errorf("instances[%d]: %w", i, err)
		}
	}
	return nil
}

// NodeGroupUpdateWeightRequest is request body to update node group weight.
type NodeGroupUpdateWeightRequest struct {
	Weight int `json:"weight"`
}

func (r NodeGroupUpdateWeightRequest) validate() error {
	if r.Weight < 0 {
		return errors.New("node group weight must be greater than or equal to zero") //nolint:goerr113 // Dynamic error
	}
	return nil
}

// NodeGroupUpdateShardGroupsRequest is the request body to update shard goups for a node group.
type NodeGroupUpdateShardGroupsRequest struct {
	ShardGroups []string `json:"shard_groups"`
}

func (r NodeGroupUpdateShardGroupsRequest) validate() error {
	if len(r.ShardGroups) == 0 {
		return fmt.Errorf("shard_groups: %w", common.ErrFieldEmptySlice)
	}
	for i, ID := range r.ShardGroups {
		if err := uuid.Validate(ID); err != nil {
			return fmt.Errorf("shard_groups[%d]: %w", i, err)
		}
	}
	return nil
}

// NodeGroupUpdateFloatingIPs is the request body to update public ips for a node group.
type NodeGroupUpdateFloatingIPsRequest struct {
	HasPublicIPs bool `json:"has_public_ips"`
}

// NodeGroupService is service to interact with clickhouse node group resource.
type NodeGroupService struct {
	*common.EngineService
}

func (s *NodeGroupService) nodeGroupsPath(datastoreID string, parts ...string) string {
	path := []string{"node_groups"}
	path = append(path, parts...)

	return s.DatastorePath(datastoreID, path...)
}

func (s *NodeGroupService) CreateNodeGroup(
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

func (s *NodeGroupService) DeleteNodeGroup(ctx context.Context, datastoreID, nodeGroupID string) error {
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

func (s *NodeGroupService) ResizeNodeGroup(
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

func (s *NodeGroupService) DeleteNodeGroupInstances(
	ctx context.Context, datastoreID, nodeGroupID string, body NodeGroupDeleteInstancesRequest,
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

	err := s.Patch(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID, "instances"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

func (s *NodeGroupService) UpdateNodeGroupWeight(
	ctx context.Context, datastoreID, nodeGroupID string, body NodeGroupUpdateWeightRequest,
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

	err := s.Patch(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID, "weight"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

func (s *NodeGroupService) UpdateNodeGroupShardGroups(
	ctx context.Context, datastoreID, nodeGroupID string, body NodeGroupUpdateShardGroupsRequest,
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

	err := s.Patch(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID, "shard-groups"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

func (s *NodeGroupService) UpdateNodeGroupFloatingIPs(
	ctx context.Context, datastoreID, nodeGroupID string, body NodeGroupUpdateFloatingIPsRequest,
) (NodeGroupResponse, error) {
	response := NodeGroupResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := uuid.Validate(nodeGroupID); err != nil {
		return response, fmt.Errorf("validate node group id: %w", err)
	}

	err := s.Patch(ctx, s.nodeGroupsPath(datastoreID, nodeGroupID, "floating_ips"), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}
