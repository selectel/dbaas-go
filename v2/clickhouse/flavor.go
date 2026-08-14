package clickhouse

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/common"
)

// FlavorResponse is the API response for the clickhouse flavor.
type FlavorResponse struct {
	ID       string                `json:"id"`
	DiskType common.FlavorDiskType `json:"disk_type"`
	FlSize   string                `json:"fl_size"`
	Disk     int                   `json:"disk"`
	RAM      int                   `json:"ram"`
	VCPUs    int                   `json:"vcpus"`
	// "type": "FIXED",
	// "subtype": "STANDARD"
}

// FlavorForNodeGroupCreate is body to create NodeGroup.
type FlavorForNodeGroupCreate struct {
	ID       string                `json:"id,omitempty"`
	DiskType common.FlavorDiskType `json:"disk_type"`
	Type     common.FlavorType     `json:"type"`
	Disk     int                   `json:"disk,omitempty"`
	RAM      int                   `json:"ram,omitempty"`
	VCPUs    int                   `json:"vcpus,omitempty"`
}

func (f FlavorForNodeGroupCreate) validate() error {
	switch f.Type {
	case common.FlavorTypeFIXED:
		if err := uuid.Validate(f.ID); err != nil {
			return fmt.Errorf("validate flavor id: %w", err)
		}

	case common.FlavorTypeFlexible:
		if f.ID != "" {
			return errors.New("flavor.id must not be specified for FLEXIBLE flavor") //nolint:goerr113 // Dynamic error
		}

		if f.Disk <= 0 {
			return errors.New("flavor.disk must be greater than 0") //nolint:goerr113 // Dynamic error
		}

		if f.RAM <= 0 {
			return errors.New("flavor.ram must be greater than 0") //nolint:goerr113 // Dynamic error
		}

		if f.VCPUs <= 0 {
			return errors.New("flavor.vcpus must be greater than 0") //nolint:goerr113 // Dynamic error
		}

	default:
		return fmt.Errorf("unsupported flavor type %q", f.Type)
	}

	// API requires DiskType for bouth types.

	return nil
}
