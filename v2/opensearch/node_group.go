package opensearch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/common"
)

// NodeGroupStatus represents custom type for DBaaS v2 opensearch node group statuses.
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
	NodeGroupRoleData      NodeGroupRole = "DATA"
	NodeGroupRoleDashboard NodeGroupRole = "DASHBOARD"
	NodeGroupRoleManager   NodeGroupRole = "MANAGER"
)

// NodeGroupResponse is the API response for the opensearch node group.
type NodeGroupResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Role         NodeGroupRole      `json:"role"`
	Status       NodeGroupStatus    `json:"status"`
	Instances    []InstanceResponse `json:"instances"`
	Flavor       FlavorResponse     `json:"flavor"`
	NodeCount    int                `json:"node_count"`
	HasPublicIPs bool               `json:"has_public_ips"`
}

// NodeGroupCreateRequest is the request body for creating a opensearch node group.
type NodeGroupCreateRequest struct {
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

	if n.Role != NodeGroupRoleData && n.Role != NodeGroupRoleDashboard && n.Role != NodeGroupRoleManager {
		return fmt.Errorf("node_group.role: %w", common.ErrUnsupportedNodeGroupRole)
	}

	if err := n.Flavor.validate(); err != nil {
		return fmt.Errorf("node_group.flavor: %w", err)
	}

	if n.NodeCount <= 0 {
		return fmt.Errorf("node_group.node_count: %w", common.ErrPositiveIntegerRequired)
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

// NodeGroupUpdateFloatingIPs is the request body to update public ips for a node group.
type NodeGroupUpdateFloatingIPsRequest struct {
	HasPublicIPs bool `json:"has_public_ips"`
}

// DatastoreService is service to interact with opensearch datastore resource.
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
