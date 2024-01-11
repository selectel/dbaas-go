package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ACL is the API response for the acls.
type ACL struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ProjectID   string `json:"project_id"`
	DatastoreID string `json:"datastore_id"`
	Pattern     string `json:"pattern"`
	PatternType string `json:"pattern_type"`
	UserID      string `json:"user_id"`
	Status      Status `json:"status"`
	AllowRead   bool   `json:"allow_read"`
	AllowWrite  bool   `json:"allow_write"`
}

// ACLCreateOpts represents options for the acl Create request.
type ACLCreateOpts struct {
	DatastoreID string `json:"datastore_id"`
	Pattern     string `json:"pattern,omitempty"`
	PatternType string `json:"pattern_type"`
	UserID      string `json:"user_id"`
	AllowRead   bool   `json:"allow_read"`
	AllowWrite  bool   `json:"allow_write"`
}

// ACLUpdateOpts represents options for the acl Update request.
type ACLUpdateOpts struct {
	AllowRead  bool `json:"allow_read"`
	AllowWrite bool `json:"allow_write"`
}

// ACLQueryParams represents available query parameters for the acl.
type ACLQueryParams struct {
	ID          string `json:"id,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	DatastoreID string `json:"datastore_id,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
	PatternType string `json:"pattern_type,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Status      Status `json:"status,omitempty"`
}

// ACLs returns all ACLs.
func (api *API) ACLs(ctx context.Context, params *ACLQueryParams) ([]ACL, error) {
	uri, err := setQueryParams("/acls", params)
	if err != nil {
		return []ACL{}, err
	}

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []ACL{}, err
	}

	var result struct {
		ACLs []ACL `json:"acls"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []ACL{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ACLs, nil
}

// ACL returns an ACL based on the ID.
func (api *API) ACL(ctx context.Context, aclID string) (ACL, error) {
	uri := fmt.Sprintf("/acls/%s", aclID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return ACL{}, err
	}

	var result struct {
		ACL ACL `json:"acl"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ACL{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ACL, nil
}

// CreateACL creates a new acl.
func (api *API) CreateACL(ctx context.Context, opts ACLCreateOpts) (ACL, error) {
	uri := "/acls"
	createACLOpts := struct {
		ACL ACLCreateOpts `json:"acl"`
	}{
		ACL: opts,
	}
	requestBody, err := json.Marshal(createACLOpts)
	if err != nil {
		return ACL{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, uri, requestBody)
	if err != nil {
		return ACL{}, err
	}

	var result struct {
		ACL ACL `json:"acl"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ACL{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ACL, nil
}

// UpdateACL updates an existing acl.
func (api *API) UpdateACL(ctx context.Context, aclID string, opts ACLUpdateOpts) (ACL, error) {
	uri := fmt.Sprintf("/acls/%s", aclID)
	updateACLOpts := struct {
		ACL ACLUpdateOpts `json:"acl"`
	}{
		ACL: opts,
	}
	requestBody, err := json.Marshal(updateACLOpts)
	if err != nil {
		return ACL{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return ACL{}, err
	}

	var result struct {
		ACL ACL `json:"acl"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ACL{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ACL, nil
}

// DeleteACL deletes an existing acl.
func (api *API) DeleteACL(ctx context.Context, aclID string) error {
	uri := fmt.Sprintf("/acls/%s", aclID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}
