package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LogicalReplicationSlot is the API response for the logical replication slot.
type LogicalReplicationSlot struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DatastoreID string `json:"datastore_id"`
	DatabaseID  string `json:"database_id"`
	Status      Status `json:"status"`
}

type LogicalReplicationSlotCreateOpts struct {
	Name        string `json:"name"`
	DatastoreID string `json:"datastore_id"`
	DatabaseID  string `json:"database_id"`
}

type LogicalReplicationSlotQueryParams struct {
	ID          string `json:"id,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	Name        string `json:"name,omitempty"`
	DatastoreID string `json:"datastore_id,omitempty"`
	DatabaseID  string `json:"database_id,omitempty"`
	Status      Status `json:"status,omitempty"`
}

const LogicalReplicationSlotsURI = "/logical-replication-slots"

// LogicalReplicationSlots returns all slots.
func (api *API) LogicalReplicationSlots(
	ctx context.Context,
	params *LogicalReplicationSlotQueryParams,
) ([]LogicalReplicationSlot, error) {
	uri, err := setQueryParams(LogicalReplicationSlotsURI, params)
	if err != nil {
		return []LogicalReplicationSlot{}, err
	}

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []LogicalReplicationSlot{}, err
	}

	var result struct {
		LogicalReplicationSlots []LogicalReplicationSlot `json:"logical-replication-slots"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []LogicalReplicationSlot{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.LogicalReplicationSlots, nil
}

// LogicalReplicationSlot returns a slot based on the ID.
func (api *API) LogicalReplicationSlot(ctx context.Context, slotID string) (LogicalReplicationSlot, error) {
	uri := fmt.Sprintf("%s/%s", LogicalReplicationSlotsURI, slotID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return LogicalReplicationSlot{}, err
	}

	var result struct {
		LogicalReplicationSlot LogicalReplicationSlot `json:"logical-replication-slot"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return LogicalReplicationSlot{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.LogicalReplicationSlot, nil
}

// CreateLogicalReplicationSlot creates a new slot.
func (api *API) CreateLogicalReplicationSlot(
	ctx context.Context,
	opts LogicalReplicationSlotCreateOpts,
) (LogicalReplicationSlot, error) {
	createLogicalReplicationSlotsOpts := struct {
		LogicalReplicationSlot LogicalReplicationSlotCreateOpts `json:"logical-replication-slot"`
	}{
		LogicalReplicationSlot: opts,
	}
	requestBody, err := json.Marshal(createLogicalReplicationSlotsOpts)
	if err != nil {
		return LogicalReplicationSlot{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, LogicalReplicationSlotsURI, requestBody)
	if err != nil {
		return LogicalReplicationSlot{}, err
	}
	var result struct {
		LogicalReplicationSlot LogicalReplicationSlot `json:"logical-replication-slot"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return LogicalReplicationSlot{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.LogicalReplicationSlot, nil
}

// DeleteLogicalReplicationSlot deletes an existing slot.
func (api *API) DeleteLogicalReplicationSlot(ctx context.Context, slotID string) error {
	uri := fmt.Sprintf("%s/%s", LogicalReplicationSlotsURI, slotID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}
