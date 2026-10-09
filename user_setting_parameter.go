package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// UserSettingParameter is the API response for the user setting parameters.
type UserSettingParameter struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Unit                   string `json:"unit"`
	ApplyMechanism         string `json:"apply_mechanism"`
	Min                    any    `json:"min"`
	Max                    any    `json:"max"`
	DefaultValue           string `json:"default_value"`
	Choices                []any  `json:"choices"`
	IsAvailableForCustomer bool   `json:"is_available_for_customer"`
	IsChangeable           bool   `json:"is_changeable"`
	CanBeEmpty             bool   `json:"can_be_empty"`
}

// UserSettingParameterQueryParams represents available query parameters
// for listing user setting parameters.
type UserSettingParameterQueryParams struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	DatastoreTypeID string `json:"datastore_type_id,omitempty"`
	DatastoreID     string `json:"datastore_id,omitempty"`
}

const UserSettingParametersURI = "/user-setting-parameters"

// UserSettingParameters returns user setting parameters.
func (api *API) UserSettingParameters(
	ctx context.Context,
	params *UserSettingParameterQueryParams,
) ([]UserSettingParameter, error) {
	uri, err := setQueryParams(UserSettingParametersURI, params)
	if err != nil {
		return []UserSettingParameter{}, err
	}

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []UserSettingParameter{}, err
	}

	var result struct {
		UserSettingParameters []UserSettingParameter `json:"user-setting-parameters"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []UserSettingParameter{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.UserSettingParameters, nil
}

// UserSettingParameter returns a user setting parameter based on the ID.
func (api *API) UserSettingParameter(
	ctx context.Context,
	userSettingParameterID string,
) (UserSettingParameter, error) {
	if err := uuid.Validate(userSettingParameterID); err != nil {
		return UserSettingParameter{}, fmt.Errorf("validate user setting parameter id: %w", err)
	}

	uri := fmt.Sprintf("%s/%s", UserSettingParametersURI, userSettingParameterID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return UserSettingParameter{}, err
	}

	var result struct {
		UserSettingParameter UserSettingParameter `json:"user-setting-parameter"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return UserSettingParameter{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.UserSettingParameter, nil
}
