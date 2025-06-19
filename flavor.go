package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type FlavorHost struct {
	Line           string `json:"line"`
	Processor      string `json:"processor"`
	AvailableCount int    `json:"available_count"`
}

// FlavorResponse is the API response for the flavors.
type FlavorResponse struct {
	Host             *FlavorHost `json:"host"`
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	FlSize           string      `json:"fl_size"`
	DatastoreTypeIDs []string    `json:"datastore_type_ids"`
	Vcpus            int         `json:"vcpus"`
	RAM              int         `json:"ram"`
	Disk             int         `json:"disk"`
}

const FlavorsURI = "/flavors"

// Flavors returns all flavors.
func (api *API) Flavors(ctx context.Context) ([]FlavorResponse, error) {
	resp, err := api.makeRequest(ctx, http.MethodGet, FlavorsURI, nil)
	if err != nil {
		return []FlavorResponse{}, err
	}

	var result struct {
		Flavors []FlavorResponse `json:"flavors"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []FlavorResponse{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Flavors, nil
}

// Flavor returns a flavor based on the ID.
func (api *API) Flavor(ctx context.Context, flavorID string) (FlavorResponse, error) {
	uri := fmt.Sprintf("%s/%s", FlavorsURI, flavorID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return FlavorResponse{}, err
	}

	var result struct {
		Flavor FlavorResponse `json:"flavor"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return FlavorResponse{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Flavor, nil
}
