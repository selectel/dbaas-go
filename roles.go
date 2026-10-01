package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// User is the API response for the users.
type Roles struct {
	ID              string   `json:"id"`
	DatastoreTypeID string   `json:"datastore_type_id"`
	Name            string   `json:"name"`
}

const RolesURI = "/roles"

// ConfigurationParameters returns all configuration parameters.
func (api *API) Roles(ctx context.Context) ([]Roles, error) {
	resp, err := api.makeRequest(ctx, http.MethodGet, RolesURI, nil)
	if err != nil {
		return []Roles{}, err
	}

	var result struct {
		Roles []Roles `json:"roles"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Roles{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Roles, nil
}