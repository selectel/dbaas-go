package clickhouse

import "github.com/selectel/dbaas-go/v2/common"

type FlavorType string

const (
	FlavorTypeFIXED    FlavorType = "FIXED"
	FlavorTypeFlexible FlavorType = "FLEXIBLE"
)

// Flavor storers required fields for flavor model in api.
type Flavor struct {
	ID       string                `json:"id"`
	DiskType common.FlavorDiskType `json:"disk_type"`
	FlSize   string                `json:"fl_size"`
	Disk     int                   `json:"disk"`
	RAM      int                   `json:"ram"`
	VCPUs    int                   `json:"vcpus"`
}
