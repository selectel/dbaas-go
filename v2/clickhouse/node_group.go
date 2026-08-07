package clickhouse

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

// NodeGroup is the API response for the clickhouse node group.
type NodeGroup struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Role         NodeGroupRole   `json:"role"`
	Status       NodeGroupStatus `json:"status"`
	ShardGroups  []string        `json:"shard_groups"`
	Instances    []Instance      `json:"instances"`
	Flavor       Flavor          `json:"flavor"`
	NodeCount    int             `json:"node_count"`
	Weight       int             `json:"weight"`
	HasPublicIPs bool            `json:"has_public_ips"`
}
