package dbaas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const grantID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testGrantNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "grant 123 not found."
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

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/grants",
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

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestGrant(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/grants/"+grantID,
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

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestGrantNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/grants/123",
		httpmock.NewStringResponder(404, testGrantNotFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "grant 123 not found."

	_, err := testClient.Grant(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}

func TestCreateGrant(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/grants",
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

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}
