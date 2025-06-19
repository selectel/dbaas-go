package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ConfigurationParameter is the API response for the configuration parameters.
type ConfigurationParameter struct {
	ID                string `json:"id"`
	DatastoreTypeID   string `json:"datastore_type_id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Unit              string `json:"unit"`
	Min               any    `json:"min"`
	Max               any    `json:"max"`
	DefaultValue      any    `json:"default_value"`
	Choices           []any  `json:"choices"`
	InvalidValues     []any  `json:"invalid_values"`
	IsRestartRequired bool   `json:"is_restart_required"`
	IsChangeable      bool   `json:"is_changeable"`
}

const ConfigurationParametersURI = "/configuration-parameters"

// ConfigurationParameters returns all configuration parameters.
func (api *API) ConfigurationParameters(ctx context.Context) ([]ConfigurationParameter, error) {
	resp, err := api.makeRequest(ctx, http.MethodGet, ConfigurationParametersURI, nil)
	if err != nil {
		return []ConfigurationParameter{}, err
	}

	var result struct {
		ConfigurationParameters []ConfigurationParameter `json:"configuration-parameters"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []ConfigurationParameter{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ConfigurationParameters, nil
}

// ConfigurationParameter returns a configuration parameter based on the ID.
func (api *API) ConfigurationParameter(
	ctx context.Context,
	configurationParameterID string,
) (ConfigurationParameter, error) {
	uri := fmt.Sprintf("%s/%s", ConfigurationParametersURI, configurationParameterID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return ConfigurationParameter{}, err
	}

	var result struct {
		ConfigurationParameter ConfigurationParameter `json:"configuration-parameter"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ConfigurationParameter{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.ConfigurationParameter, nil
}
