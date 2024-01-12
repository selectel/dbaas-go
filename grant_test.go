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

const grantID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testGrantNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "grant %s not found."
	}
}`

const testGrantResponse = `
{
	"grant": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"status": "ACTIVE"
	}
}`

const testGrantsResponse = `
{
	"grants": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"status": "ACTIVE"
		}
	]
}`

func TestGrants(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+GrantsURI,
		httpmock.NewStringResponder(200, testGrantsResponse))

	expected := []Grant{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			Status:      StatusActive,
		},
	}

	actual, err := testClient.Grants(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGrant(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+GrantsURI+"/"+grantID,
		httpmock.NewStringResponder(200, testGrantResponse))

	expected := Grant{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusActive,
	}

	actual, err := testClient.Grant(context.Background(), grantID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGrantNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testGrantNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+GrantsURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("grant %s not found.", NotFoundEntityID)

	_, err := testClient.Grant(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateGrant(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+GrantsURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&GrantCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			grants := make(map[string]Grant)
			grants["grant"] = Grant{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Status:      StatusPendingCreate,
			}

			resp, err := httpmock.NewJsonResponse(200, grants)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createGrantOpts := GrantCreateOpts{
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	expected := Grant{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusPendingCreate,
	}

	actual, err := testClient.CreateGrant(context.Background(), createGrantOpts)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
