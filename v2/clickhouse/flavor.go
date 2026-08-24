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

// FlavorForNodeGroupRequest is body to create, resize NodeGroup.
type FlavorForNodeGroupRequest struct {
	ID       string                `json:"id,omitempty"`
	DiskType common.FlavorDiskType `json:"disk_type"`
	Type     common.FlavorType     `json:"type"`
	Disk     int                   `json:"disk,omitempty"`
	RAM      int                   `json:"ram,omitempty"`
	VCPUs    int                   `json:"vcpus,omitempty"`
}

func (f FlavorForNodeGroupRequest) validate() error {
	switch f.Type {
	case common.FlavorTypeFIXED:
		if err := uuid.Validate(f.ID); err != nil {
			return fmt.Errorf("flavor.id must be a valid UUID for FIXED flavor: %w", err)
		}

	case common.FlavorTypeFlexible:
		if f.ID != "" {
			return errors.New("flavor.id must not be specified for FLEXIBLE flavor") //nolint:goerr113 // Dynamic error
		}

		if f.Disk <= 0 {
			return fmt.Errorf("flavor.disk: %w", common.ErrPositiveIntegerRequired)
		}

		if f.RAM <= 0 {
			return fmt.Errorf("flavor.ram: %w", common.ErrPositiveIntegerRequired)
		}

		if f.VCPUs <= 0 {
			return fmt.Errorf("flavor.vcpus: %w", common.ErrPositiveIntegerRequired)
		}

	default:
		return fmt.Errorf("%w: %q", common.ErrUnsupportedFlavorType, f.Type)
	}

	// API requires DiskType for both types. It seems that the fixed flavor should not have this field as required.
	// if f.DiskType == "" {
	// 	return fmt.Errorf("flavor.disk_type %w", common.ErrFieldRequired)
	// }

	return nil
}
