package clickhouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/selectel/dbaas-go/v2/common"
)

// ShardGroupResponse is the API response for the clickhouse shard group.
type ShardGroupResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShardIDs    []string `json:"shard_ids"`
}

// ShardGroupCreateRequest is the request body to create a clickhouse shard group.
type ShardGroupCreateRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	ShardIDs    []string `json:"shard_ids"`
}

func (s ShardGroupCreateRequest) validate() error {
	if s.Name == "" {
		return fmt.Errorf("shard_group.name: %w", common.ErrFieldRequired)
	}

	if len(s.ShardIDs) == 0 {
		return fmt.Errorf("shard_group.shard_ids: %w", common.ErrFieldEmptySlice)
	}

	for i, grID := range s.ShardIDs {
		if err := uuid.Validate(grID); err != nil {
			return fmt.Errorf("shard_group.shard_ids[%d]: %w", i, err)
		}
	}
	return nil
}

// ShardGroupUpdateRequest is the request body to update a clickhouse shard group.
type ShardGroupUpdateRequest struct {
	Description *string   `json:"description"`
	ShardIDs    *[]string `json:"shard_ids"`
}

func (s ShardGroupUpdateRequest) validate() error {
	if s.ShardIDs != nil {
		if len(*s.ShardIDs) == 0 {
			return fmt.Errorf("shard_group.shard_ids: %w", common.ErrFieldEmptySlice)
		}

		for i, grID := range *s.ShardIDs {
			if err := uuid.Validate(grID); err != nil {
				return fmt.Errorf("shard_group.shard_ids[%d]: %w", i, err)
			}
		}
	}

	return nil
}

// ShardGroupService is service to interact with clickhouse shard group resource.
type ShardGroupService struct {
	*common.EngineService
}

func (s *ShardGroupService) shardGroupsPath(datastoreID string, parts ...string) string {
	path := []string{"shard-groups"}
	path = append(path, parts...)

	return s.DatastorePath(datastoreID, path...)
}

func (s *ShardGroupService) GetShardGroups(ctx context.Context, datastoreID string) ([]ShardGroupResponse, error) {
	var response struct {
		ShardGroups []ShardGroupResponse `json:"shard_groups"`
	}

	err := s.Get(ctx, s.shardGroupsPath(datastoreID), &response)
	if err != nil {
		return response.ShardGroups, err //nolint:wrapcheck
	}

	return response.ShardGroups, nil
}

func (s *ShardGroupService) CreateShardGroup(
	ctx context.Context, datastoreID string, body ShardGroupCreateRequest,
) (ShardGroupResponse, error) {
	response := ShardGroupResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}
	err := s.Post(ctx, s.shardGroupsPath(datastoreID), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}

func (s *ShardGroupService) DeleteShardGroup(ctx context.Context, datastoreID, shardGroupID string) error {
	if err := uuid.Validate(datastoreID); err != nil {
		return fmt.Errorf("validate datastore id: %w", err)
	}

	if err := uuid.Validate(shardGroupID); err != nil {
		return fmt.Errorf("validate shard group id: %w", err)
	}

	err := s.Delete(ctx, s.shardGroupsPath(datastoreID, shardGroupID))
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func (s *ShardGroupService) UpdateShardGroup(
	ctx context.Context, datastoreID, shardGroupID string, body ShardGroupUpdateRequest,
) (ShardGroupResponse, error) {
	response := ShardGroupResponse{}

	if err := uuid.Validate(datastoreID); err != nil {
		return response, fmt.Errorf("validate datastore id: %w", err)
	}

	if err := uuid.Validate(shardGroupID); err != nil {
		return response, fmt.Errorf("validate shard group id: %w", err)
	}

	if err := body.validate(); err != nil {
		return response, fmt.Errorf("validate body: %w", err)
	}

	err := s.Patch(ctx, s.shardGroupsPath(datastoreID, shardGroupID), body, &response)
	if err != nil {
		return response, err //nolint:wrapcheck
	}

	return response, nil
}
