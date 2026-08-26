package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// UserCreateOpts represents options for the user Create request.
type UserCreateOpts struct {
	Settings    map[string]any `json:"settings,omitempty"`
	Name        string         `json:"name,omitempty"`
	Password    string         `json:"password,omitempty"`
	DatastoreID string         `json:"datastore_id,omitempty"`
	Roles       []string       `json:"roles,omitempty"`
}

// UserUpdateOpts represents options for the user Update request.
type UserUpdateOpts struct {
	Password string `json:"password"`
}

// UserRolesUpdateOpts represents options for the user roles Update request.
type UserRolesUpdateOpts struct {
	Roles []string `json:"roles"`
}

// UserSettingsUpdateOpts represents options for the user settings Update request.
type UserSettingsUpdateOpts struct {
	Settings map[string]any `json:"settings"`
}

// User is the API response for the users.
type User struct {
	Settings    map[string]any `json:"settings,omitempty"`
	ID          string         `json:"id"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	ProjectID   string         `json:"project_id"`
	DatastoreID string         `json:"datastore_id"`
	Name        string         `json:"name"`
	Status      Status         `json:"status"`
	Roles       []string       `json:"roles"`
}

const (
	UsersURI              = "/users"
	UserSettingsURISuffix = "settings"
)

// User returns a user based on the ID.
func (api *API) User(ctx context.Context, userID string) (User, error) {
	uri := fmt.Sprintf("%s/%s", UsersURI, userID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return User{}, err
	}

	var result struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.User, nil
}

// Users returns all users.
func (api *API) Users(ctx context.Context) ([]User, error) {
	resp, err := api.makeRequest(ctx, http.MethodGet, UsersURI, nil)
	if err != nil {
		return []User{}, err
	}

	var result struct {
		Users []User `json:"users"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Users, nil
}

// CreateUser creates a new user.
func (api *API) CreateUser(ctx context.Context, opts UserCreateOpts) (User, error) {
	if opts.Settings != nil {
		opts.Settings = convertConfigValues(opts.Settings)
	}
	createUserOpts := struct {
		User UserCreateOpts `json:"user"`
	}{
		User: opts,
	}
	requestBody, err := json.Marshal(createUserOpts)
	if err != nil {
		return User{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, UsersURI, requestBody)
	if err != nil {
		return User{}, err
	}

	var result struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.User, nil
}

// DeleteUser deletes an existing user.
func (api *API) DeleteUser(ctx context.Context, userID string) error {
	uri := fmt.Sprintf("%s/%s", UsersURI, userID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user.
func (api *API) UpdateUser(ctx context.Context, userID string, opts UserUpdateOpts) (User, error) {
	uri := fmt.Sprintf("%s/%s", UsersURI, userID)
	updateUserOpts := struct {
		User UserUpdateOpts `json:"user"`
	}{
		User: opts,
	}
	requestBody, err := json.Marshal(updateUserOpts)
	if err != nil {
		return User{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return User{}, err
	}

	var result struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.User, nil
}

// UpdateUserRoles updates roles for an existing user.
func (api *API) UpdateUserRoles(ctx context.Context, userID string, opts UserRolesUpdateOpts) (User, error) {
	if err := uuid.Validate(userID); err != nil {
		return User{}, fmt.Errorf("validate user id: %w", err)
	}

	uri := fmt.Sprintf("%s/%s/roles", UsersURI, userID)

	requestBody, err := json.Marshal(opts)
	if err != nil {
		return User{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return User{}, err
	}

	var result struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.User, nil
}

// UpdateUserSettings updates PostgreSQL role settings of an existing user.
// Only provided keys are changed; omitted keys are kept. Send nil as a value to unset a setting.
func (api *API) UpdateUserSettings(ctx context.Context, userID string, opts UserSettingsUpdateOpts) (User, error) {
	uri := fmt.Sprintf("%s/%s/%s", UsersURI, userID, UserSettingsURISuffix)
	if opts.Settings != nil {
		opts.Settings = convertConfigValues(opts.Settings)
	}
	requestBody, err := json.Marshal(opts)
	if err != nil {
		return User{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return User{}, err
	}

	var result struct {
		User User `json:"user"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return User{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.User, nil
}
