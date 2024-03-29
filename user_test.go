package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const userID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testUserNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "user %s not found."
	}
}`

const testCreateUserInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": 
			"Validation failure: {'user.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testUserResponse = `
{
	"user": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "user",
		"status": "ACTIVE"
	}
}`

const testUsersResponse = `
{
	"users": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "user",
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "user123",
			"status": "ACTIVE"
		}
	]
}`

func TestUsers(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+UsersURI,
		httpmock.NewStringResponder(200, testUsersResponse))

	expected := []User{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:        "user",
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:        "user123",
			Status:      StatusActive,
		},
	}

	actual, err := testClient.Users(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUser(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+UsersURI+"/"+userID,
		httpmock.NewStringResponder(200, testUserResponse))

	expected := User{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "user",
		Status:      StatusActive,
	}

	actual, err := testClient.User(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUserNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testUserNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+UsersURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("user %s not found.", NotFoundEntityID)

	_, err := testClient.User(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateUser(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+UsersURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&UserCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			users := make(map[string]User)
			users["user"] = User{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Name:        "user",
				Status:      StatusPendingCreate,
			}

			resp, err := httpmock.NewJsonResponse(200, users)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createUserOpts := UserCreateOpts{
		Name:        "user",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Password:    "secret",
	}

	expected := User{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "user",
		Status:      StatusPendingCreate,
	}

	actual, err := testClient.CreateUser(context.Background(), createUserOpts)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCreateUserInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+UsersURI,
		httpmock.NewStringResponder(400, testCreateUserInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'user.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createUserOpts := UserCreateOpts{
		Name:        "user",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f",
		Password:    "secret",
	}

	_, err := testClient.CreateUser(context.Background(), createUserOpts)

	require.ErrorAs(t, err, &expected)
}

func TestUpdateUser(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+UsersURI+"/"+userID,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&UserUpdateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			users := make(map[string]User)
			users["user"] = User{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Name:        "user",
				Status:      StatusPendingUpdate,
			}

			resp, err := httpmock.NewJsonResponse(200, users)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	updateUserOpts := UserUpdateOpts{
		Password: "secret",
	}

	expected := User{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "user",
		Status:      StatusPendingUpdate,
	}

	actual, err := testClient.UpdateUser(context.Background(), userID, updateUserOpts)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
