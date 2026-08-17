package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type validator interface {
	validate() error
}

type validationTest struct {
	name   string
	body   validator
	errMsg string
}

func checkValidationTests(t *testing.T, tests []validationTest) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.body.validate()

			if tt.errMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDatastoreCreateRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "datastore without name",
			body:   DatastoreCreateRequest{},
			errMsg: "datastore.name is required",
		},
		{
			name: "datastore without password",
			body: DatastoreCreateRequest{
				Name: "Test",
			},
			errMsg: "datastore.password is required",
		},
		{
			name: "datastore without type_id",
			body: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
			},
			errMsg: "datastore.type_id must be a valid UUID",
		},
		{
			name: "datastore without subnet_id",
			body: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
			},
			errMsg: "datastore.subnet_id must be a valid UUID",
		},
		{
			name: "datastore without node_groups",
			body: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
			},
			errMsg: "datastore.node_groups must be at least one",
		},
		{
			name: "datastore with invalid node_group",
			body: DatastoreCreateRequest{
				Name:     "Test",
				Password: "1234",
				TypeID:   "550e8400-e29b-41d4-a716-446655440000",
				SubnetID: "550e8400-e29b-41d4-a716-446655440001",
				NodeGroups: []NodeGroupCreateRequest{
					{Name: ""},
				},
			},
			errMsg: "node_groups[0]: node_group.name is required",
		},
	}
	checkValidationTests(t, tests)
}

func TestDatastoreUpdateRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without name",
			body:   DatastoreUpdateRequest{},
			errMsg: "name is required",
		},
		{
			name: "body with name",
			body: DatastoreUpdateRequest{Name: "Test"},
		},
	}

	checkValidationTests(t, tests)
}

func TestDatastoreUpdatePasswordRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without new password",
			body:   DatastoreUpdatePasswordRequest{},
			errMsg: "password is required",
		},
		{
			name: "body with name",
			body: DatastoreUpdatePasswordRequest{NewPassword: "Test"},
		},
	}

	checkValidationTests(t, tests)
}

func TestDatastoreSecurityGroupsRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name: "body without security groups",
			body: DatastoreSecurityGroupsRequest{},
		},
		{
			name: "body with security groups",
			body: DatastoreSecurityGroupsRequest{
				SecurityGroups: []string{
					"550e8400-e29b-41d4-a716-446655440000",
				},
			},
		},
		{
			name:   "body with invalid ids",
			body:   DatastoreSecurityGroupsRequest{SecurityGroups: []string{"Test"}},
			errMsg: "security_groups[0]: invalid UUID length",
		},
	}

	checkValidationTests(t, tests)
}

func TestDatastoreLogPlatformRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without log_group",
			body:   DatastoreLogPlatformRequest{},
			errMsg: "log_platform.log_group is required",
		},
		{
			name: "body with log_group",
			body: DatastoreLogPlatformRequest{
				LogPlatform: DatastoreLogGroup{
					LogGroup: "TestGroup",
				},
			},
		},
	}

	checkValidationTests(t, tests)
}

func TestDatastoreConfigRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without config",
			body:   DatastoreConfigRequest{},
			errMsg: "config is required",
		},
		{
			name: "body with config",
			body: DatastoreConfigRequest{
				Config: map[string]any{
					"param_first":  1,
					"param_second": "test",
				},
			},
		},
	}

	checkValidationTests(t, tests)
}
