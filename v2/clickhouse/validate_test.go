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
			errMsg: "datastore.name: required field",
		},
		{
			name: "datastore without password",
			body: DatastoreCreateRequest{
				Name: "Test",
			},
			errMsg: "datastore.password: required field",
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
			errMsg: "datastore.node_groups: must be at least one value",
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
			errMsg: "node_groups[0]: node_group.name: required field",
		},
	}
	checkValidationTests(t, tests)
}

func TestDatastoreUpdateRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without name",
			body:   DatastoreUpdateRequest{},
			errMsg: "name: required field",
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
			errMsg: "password: required field",
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
			errMsg: "log_platform.log_group: required field",
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
			errMsg: "config: required field",
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

func TestNodeGroupCreateRequest_validate(t *testing.T) {
	weight := 100

	tests := []validationTest{
		{
			name: "node group without name",
			body: NodeGroupCreateRequest{
				Role:   NodeGroupRoleData,
				Flavor: FlavorForNodeGroupRequest{},
			},
			errMsg: "node_group.name: required field",
		},
		{
			name: "node group without role",
			body: NodeGroupCreateRequest{
				Name:   "TestNg",
				Flavor: FlavorForNodeGroupRequest{},
			},
			errMsg: "node_group.role: unsupported node group role",
		},
		{
			name: "node group without flavor",
			body: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
			},
			errMsg: "node_group.flavor: unsupported flavor type: \"\"",
		},
		{
			name: "node group with invalid flavor",
			body: NodeGroupCreateRequest{
				Name:   "TestNg",
				Role:   NodeGroupRoleData,
				Flavor: FlavorForNodeGroupRequest{Type: "X"},
			},
			errMsg: "node_group.flavor: unsupported flavor type: \"X\"",
		},
		{
			name: "node group without node_count",
			body: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
				Flavor: FlavorForNodeGroupRequest{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
			},
			errMsg: "node_group.node_count: value must be greater than 0",
		},
		{
			name: "node group data with weight",
			body: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleData,
				Flavor: FlavorForNodeGroupRequest{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
				NodeCount: 1,
				Weight:    &weight,
			},
		},
		{
			name: "node group keeper with weight",
			body: NodeGroupCreateRequest{
				Name: "TestNg",
				Role: NodeGroupRoleKeeper,
				Flavor: FlavorForNodeGroupRequest{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
				NodeCount: 1,
				Weight:    &weight,
			},
			errMsg: "node_group.role KEEPER could not have weight",
		},
	}

	checkValidationTests(t, tests)
}

func TestNodeGroupResizeRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name: "body without flavor",
			body: NodeGroupResizeRequest{
				NodeCount: 1,
			},
			errMsg: "validate flavor: unsupported flavor type",
		},
		{
			name: "body without node_count",
			body: NodeGroupResizeRequest{
				Flavor: FlavorForNodeGroupRequest{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
			},
			errMsg: "node_count: value must be greater than 0",
		},
		{
			name: "good body",
			body: NodeGroupResizeRequest{
				Flavor: FlavorForNodeGroupRequest{
					ID:   "550e8400-e29b-41d4-a716-446655440000",
					Type: "FIXED",
				},
				NodeCount: 2,
			},
		},
	}

	checkValidationTests(t, tests)
}

func TestNodeGroupDeleteInstancesRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without instances",
			body:   NodeGroupDeleteInstancesRequest{},
			errMsg: "instances: must be at least one value",
		},
		{
			name: "body with invalid instances",
			body: NodeGroupDeleteInstancesRequest{
				Instances: []string{"550e8400-e29b-41d4-a716-446655440000", "invalidID"},
			},
			errMsg: "instances[1]: invalid UUID",
		},
		{
			name: "good body",
			body: NodeGroupDeleteInstancesRequest{
				Instances: []string{"550e8400-e29b-41d4-a716-446655440000"},
			},
		},
	}

	checkValidationTests(t, tests)
}

func TestNodeGroupNodeGroupUpdateWeightRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name: "good body",
			body: NodeGroupUpdateWeightRequest{Weight: 100},
		},
		{
			name:   "body with invalid weight",
			body:   NodeGroupUpdateWeightRequest{Weight: -2},
			errMsg: "node group weight must be greater than or equal to zero",
		},
	}

	checkValidationTests(t, tests)
}

func TestNodeGroupUpdateShardGroupsRequest_validate(t *testing.T) {
	tests := []validationTest{
		{
			name:   "body without shard_groups",
			body:   NodeGroupUpdateShardGroupsRequest{},
			errMsg: "shard_groups: must be at least one value",
		},
		{
			name: "body with invalid shard_groups",
			body: NodeGroupUpdateShardGroupsRequest{
				ShardGroups: []string{"550e8400-e29b-41d4-a716-446655440000", "invalidID"},
			},
			errMsg: "shard_groups[1]: invalid UUID",
		},
		{
			name: "good body",
			body: NodeGroupUpdateShardGroupsRequest{
				ShardGroups: []string{"550e8400-e29b-41d4-a716-446655440000"},
			},
		},
	}

	checkValidationTests(t, tests)
}

func TestShardGroupCreateRequest_validate(t *testing.T) {
	description := "Test"
	tests := []validationTest{
		{
			name: "shard group without name",
			body: ShardGroupCreateRequest{
				ShardIDs: []string{"1"},
			},
			errMsg: "shard_group.name: required field",
		},
		{
			name: "shard group without shard_ids",
			body: ShardGroupCreateRequest{
				Name:        "TestSg",
				Description: &description,
			},
			errMsg: "shard_group.shard_ids: must be at least one value",
		},
		{
			name: "shard group with empty shard_ids",
			body: ShardGroupCreateRequest{
				Name:        "TestSg",
				Description: &description,
				ShardIDs:    []string{},
			},
			errMsg: "shard_group.shard_ids: must be at least one value",
		},
		{
			name: "shard group with invalid shard_ids",
			body: ShardGroupCreateRequest{
				Name:        "TestSg",
				Description: &description,
				ShardIDs:    []string{"550e8400-e29b-41d4-a716-446655440000", "1"},
			},
			errMsg: "shard_group.shard_ids[1]: invalid UUID",
		},
		{
			name: "good shard group",
			body: ShardGroupCreateRequest{
				Name:        "TestSg",
				Description: &description,
				ShardIDs:    []string{"550e8400-e29b-41d4-a716-446655440000"},
			},
		},
	}

	checkValidationTests(t, tests)
}

func TestShardGroupUpdateRequest_validate(t *testing.T) {
	description := "Test"
	sGIDsInvalidOne := []string{"1"}
	sGIDsInvalidTwo := []string{}
	sGIDsValid := []string{"550e8400-e29b-41d4-a716-446655440000"}

	tests := []validationTest{
		{
			name: "shard group with invalid ids only",
			body: ShardGroupUpdateRequest{
				ShardIDs: &sGIDsInvalidOne,
			},
			errMsg: "shard_group.shard_ids[0]: invalid UUID",
		},
		{
			name: "shard group with valid ids only",
			body: ShardGroupUpdateRequest{
				ShardIDs: &sGIDsValid,
			},
		},
		{
			name: "shard group with valid ids and description",
			body: ShardGroupUpdateRequest{
				Description: &description,
				ShardIDs:    &sGIDsValid,
			},
		},
		{
			name: "shard group with description only",
			body: ShardGroupUpdateRequest{
				Description: &description,
			},
		},
		{
			name: "shard group with invalid ids and description",
			body: ShardGroupUpdateRequest{
				Description: &description,
				ShardIDs:    &sGIDsInvalidTwo,
			},
			errMsg: "shard_group.shard_ids: must be at least one value",
		},
	}

	checkValidationTests(t, tests)
}
