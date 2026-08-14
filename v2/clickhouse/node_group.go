package clickhouse

import (
	"errors"
	"fmt"
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
	Weight       *int                     `json:"weight,omitempty"`
	HasPublicIPs *bool                    `json:"has_public_ips,omitempty"`
	Name         string                   `json:"name"`
	Role         NodeGroupRole            `json:"role"`
	Flavor       FlavorForNodeGroupCreate `json:"flavor"`
	NodeCount    int                      `json:"node_count"`
}

func (n NodeGroupCreateRequest) validate() error {
	if n.Name == "" {
		return errors.New("node_group.name is required") //nolint:goerr113 // Dynamic error
	}

	if err := n.validateRole(); err != nil {
		return err
	}

	if err := n.Flavor.validate(); err != nil {
		return fmt.Errorf("node_group flavor: %w", err)
	}

	if n.NodeCount <= 0 {
		return errors.New("node_count must be greater than 0") //nolint:goerr113 // Dynamic error
	}

	if n.Weight != nil && n.Role == NodeGroupRoleKeeper {
		return errors.New("KEEPER could not have Weight") //nolint:goerr113 // Dynamic error
	}
	return nil
}

func (n NodeGroupCreateRequest) validateRole() error {
	switch n.Role {
	case NodeGroupRoleData, NodeGroupRoleKeeper:
		// valid
		return nil
	default:
		return fmt.Errorf("unsupported group role %q", n.Role)
	}
}
