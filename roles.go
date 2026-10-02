package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// User is the API response for the users.
type Role struct {
	ID              string `json:"id"`
	DatastoreTypeID string `json:"datastore_type_id"`
	Name            string `json:"name"`
}

const RolesURI = "/roles"

// Roles returns all roles.
func (api *API) Roles(ctx context.Context) ([]Role, error) {
	resp, err := api.makeRequest(ctx, http.MethodGet, RolesURI, nil)
	if err != nil {
		return []Role{}, err
	}

	var result struct {
		Roles []Role `json:"roles"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Role{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Roles, nil
}
