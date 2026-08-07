package clickhouse

// InstanceStatus presents custom type for DBaaS v2 clickhouse instance statuses.
type InstanceStatus string

const (
	InstanceStatusActive   InstanceStatus = "ACTIVE"
	InstanceStatusCreating InstanceStatus = "CREATING"
	InstanceStatusDeleted  InstanceStatus = "DELETED"
	InstanceStatusDeleting InstanceStatus = "DELETING"
	InstanceStatusError    InstanceStatus = "ERROR"
	InstanceStatusResizing InstanceStatus = "RESIZING"
	InstanceStatusUpdating InstanceStatus = "UPDATING"
	InstanceStatusDiskFull InstanceStatus = "DISK_FULL"
	InstanceStatusDown     InstanceStatus = "DOWN"
)

type Instance struct {
	ID               string         `json:"id"`
	IP               string         `json:"ip"`
	FlavorID         string         `json:"flavor_id"`
	FloatingIP       string         `json:"floating_ip"`
	Status           InstanceStatus `json:"status"`
	AvailabilityZone string         `json:"availability_zone"`
	Hostname         string         `json:"hostname"`
}
