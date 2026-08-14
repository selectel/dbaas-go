package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeGroupCreateRequest_validate(t *testing.T) {
	weight := 100

	tests := []struct {
		name      string
		errMsg    string
		NodeGroup NodeGroupCreateRequest
		wantErr   bool
	}{
		{
			name: "node group without name",
			NodeGroup: NodeGroupCreateRequest{
				Role:   NodeGroupRoleData,
				Flavor: FlavorForNodeGroupCreate{},
			},
			wantErr: true,
			errMsg:  "node_group.name is required",
		},
		{
			name: "node group without role",
			NodeGroup: NodeGroupCreateRequest{
				Name:   "TestNg",
				Flavor: FlavorForNodeGroupCreate{},
			},
			wantErr: true,
			errMsg:  "node_group.role must be DATA or KEEPER",
		},
		{
			name: "node group without flavor",
			NodeGroup: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
			},
			wantErr: true,
			errMsg:  "node_group.flavor: unsupported flavor type: \"\"",
		},
		{
			name: "node group with invalid flavor",
			NodeGroup: NodeGroupCreateRequest{
				Name:   "TestNg",
				Role:   NodeGroupRoleData,
				Flavor: FlavorForNodeGroupCreate{Type: "X"},
			},
			wantErr: true,
			errMsg:  "node_group.flavor: unsupported flavor type: \"X\"",
		},
		{
			name: "node group without node_count",
			NodeGroup: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
				Flavor: FlavorForNodeGroupCreate{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
			},
			wantErr: true,
			errMsg:  "node_group.node_count must be greater than 0",
		},
		{
			name: "node group data with weight",
			NodeGroup: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
				Flavor: FlavorForNodeGroupCreate{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
				NodeCount: 1,
				Weight:    &weight,
			},
		},
		{
			name: "node group keeper with weight",
			NodeGroup: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleKeeper,
				Flavor: FlavorForNodeGroupCreate{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
				NodeCount: 1,
				Weight:    &weight,
			},
			wantErr: true,
			errMsg:  "node_group.role KEEPER could not have weight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.NodeGroup.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
