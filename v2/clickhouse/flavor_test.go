package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/selectel/dbaas-go/v2/common"
)

func TestFlavorForNodeGroupRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		errMsg  string
		flavor  FlavorForNodeGroupRequest
		wantErr bool
	}{
		{
			name: "fixed with valid id",
			flavor: FlavorForNodeGroupRequest{
				Type: "FIXED",
				ID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "fixed without id",
			flavor: FlavorForNodeGroupRequest{
				Type: "FIXED",
			},
			wantErr: true,
			errMsg:  "flavor.id must be a valid UUID for FIXED flavor",
		},
		{
			name: "fixed with invalid id",
			flavor: FlavorForNodeGroupRequest{
				Type: "FIXED",
				ID:   "invalid",
			},
			wantErr: true,
			errMsg:  "flavor.id must be a valid UUID for FIXED flavor",
		},
		{
			name: "flexible without id",
			flavor: FlavorForNodeGroupRequest{
				Type:     "FLEXIBLE",
				DiskType: common.FlavorDiskNetworkUltra,
				Disk:     10,
				RAM:      5,
				VCPUs:    2,
			},
			wantErr: false,
		},
		{
			name: "flexible without disk",
			flavor: FlavorForNodeGroupRequest{
				Type:     "FLEXIBLE",
				DiskType: common.FlavorDiskNetworkUltra,
				RAM:      5,
				VCPUs:    2,
			},
			wantErr: true,
			errMsg:  "flavor.disk must be greater than 0",
		},
		{
			name: "flexible without ram",
			flavor: FlavorForNodeGroupRequest{
				Type:     "FLEXIBLE",
				DiskType: common.FlavorDiskNetworkUltra,
				Disk:     5,
				VCPUs:    2,
			},
			wantErr: true,
			errMsg:  "flavor.ram must be greater than 0",
		},
		{
			name: "flexible without vcpus",
			flavor: FlavorForNodeGroupRequest{
				Type:     "FLEXIBLE",
				DiskType: common.FlavorDiskNetworkUltra,
				Disk:     5,
				RAM:      2,
			},
			wantErr: true,
			errMsg:  "flavor.vcpus must be greater than 0",
		},
		{
			name: "flexible with id",
			flavor: FlavorForNodeGroupRequest{
				Type: "FLEXIBLE",
				ID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
			errMsg:  "flavor.id must not be specified for FLEXIBLE flavor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.flavor.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
