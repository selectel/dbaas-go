package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FloatingIPsOpts represents create|delete options for creating|deleting
// Floating IP in|from existed cluster.
type FloatingIPsOpts struct {
	InstanceID string `json:"instance_id"`
}

const FloatingIPsURI = "/floating-ips"

// CreateFloatingIP creates FloatingIP for provided instance of an existing datastore.
func (api *API) CreateFloatingIP(ctx context.Context, opts FloatingIPsOpts) error {
	floatingIPsOpts := struct {
		FloatingIP FloatingIPsOpts `json:"floating-ip"`
	}{
		FloatingIP: opts,
	}
	requestBody, err := json.Marshal(floatingIPsOpts)
	if err != nil {
		return fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	_, err = api.makeRequest(ctx, http.MethodPost, FloatingIPsURI, requestBody)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFloatingIP deletes FloatingIP from provided instance of an existing datastore.
func (api *API) DeleteFloatingIP(ctx context.Context, opts FloatingIPsOpts) error {
	floatingIPsOpts := struct {
		FloatingIP FloatingIPsOpts `json:"floating-ip"`
	}{
		FloatingIP: opts,
	}
	requestBody, err := json.Marshal(floatingIPsOpts)
	if err != nil {
		return fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	_, err = api.makeRequest(ctx, http.MethodDelete, FloatingIPsURI, requestBody)
	if err != nil {
		return err
	}

	return nil
}
