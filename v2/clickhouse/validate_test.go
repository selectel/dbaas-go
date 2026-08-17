package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatastoreCreateRequest_validate(t *testing.T) {
	tests := []struct {
		name      string
		errMsg    string
		datastore DatastoreCreateRequest
		wantErr   bool
	}{
		{
			name:      "datastore without name",
			datastore: DatastoreCreateRequest{},
			wantErr:   true,
			errMsg:    "datastore.name is required",
		},
		{
			name: "datastore without password",
			datastore: DatastoreCreateRequest{
				Name: "Test",
			},
			wantErr: true,
			errMsg:  "datastore.password is required",
		},
		{
			name: "datastore without type_id",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
			},
			wantErr: true,
			errMsg:  "datastore.type_id must be a valid UUID",
		},
		{
			name: "datastore without subnet_id",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
			errMsg:  "datastore.subnet_id must be a valid UUID",
		},
		{
			name: "datastore without node_groups",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
			},
			wantErr: true,
			errMsg:  "datastore.node_groups must be at least one",
		},
		{
			name: "datastore with invalid node_group",
			datastore: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
				NodeGroups: []NodeGroupCreateRequest{
					{Name: ""},
				},
			},
			wantErr: true,
			errMsg:  "node_groups[0]: node_group.name is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.datastore.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDatastoreUpdateRequest_validate(t *testing.T) {
	tests := []struct {
		name      string
		errMsg    string
		datastore DatastoreUpdateRequest
		wantErr   bool
	}{
		{
			name:      "datastore without name",
			datastore: DatastoreUpdateRequest{},
			wantErr:   true,
			errMsg:    "name is required",
		},
		{
			name:      "datastore with name",
			datastore: DatastoreUpdateRequest{Name: "Test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.datastore.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDatastoreUpdatePasswordRequest_validate(t *testing.T) {
	tests := []struct {
		name      string
		errMsg    string
		datastore DatastoreUpdatePasswordRequest
		wantErr   bool
	}{
		{
			name:      "datastore without new password",
			datastore: DatastoreUpdatePasswordRequest{},
			wantErr:   true,
			errMsg:    "password is required",
		},
		{
			name:      "datastore with name",
			datastore: DatastoreUpdatePasswordRequest{NewPassword: "Test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.datastore.validate()

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
