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

const extensionID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testExtensionNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "extension %s not found."
	}
}`

const testCreateExtensionInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": 
			"Validation failure: {'extension.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testExtensionResponse = `{
	"extension": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"available_extension_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"status": "ACTIVE"
	}
}`

const testExtensionsResponse = `
{
	"extensions": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"available_extension_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"available_extension_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		}
	]
}`

func TestExtensions(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+ExtensionsURI,
		httpmock.NewStringResponder(200, testExtensionsResponse))

	expected := []Extension{
		{
			ID:                   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			AvailableExtensionID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:            "1970-01-01T00:00:00",
			UpdatedAt:            "1970-01-01T00:00:00",
			DatastoreID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:           "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:               StatusActive,
		},
		{
			ID:                   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			AvailableExtensionID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:            "1970-01-01T00:00:00",
			UpdatedAt:            "1970-01-01T00:00:00",
			DatastoreID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:           "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:               StatusActive,
		},
	}

	actual, err := testClient.Extensions(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExtension(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+ExtensionsURI+"/"+extensionID,
		httpmock.NewStringResponder(200, testExtensionResponse))

	expected := Extension{
		ID:                   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		AvailableExtensionID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:            "1970-01-01T00:00:00",
		UpdatedAt:            "1970-01-01T00:00:00",
		DatastoreID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:           "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:               StatusActive,
	}

	actual, err := testClient.Extension(context.Background(), extensionID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExtensionNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testExtensionNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+ExtensionsURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("extension %s not found.", NotFoundEntityID)

	_, err := testClient.Extension(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateExtension(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+ExtensionsURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&ExtensionCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			extensions := make(map[string]Extension)
			extensions["extension"] = Extension{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Status:      StatusPendingCreate,
			}

			resp, err := httpmock.NewJsonResponse(200, extensions)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createExtensionOpts := ExtensionCreateOpts{
		AvailableExtensionID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:           "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	expected := Extension{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusPendingCreate,
	}

	actual, err := testClient.CreateExtension(context.Background(), createExtensionOpts)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCreateExtensionInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+ExtensionsURI,
		httpmock.NewStringResponder(400, testCreateExtensionInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'extension.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createExtensionOpts := ExtensionCreateOpts{
		AvailableExtensionID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f",
		DatabaseID:           "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	_, err := testClient.CreateExtension(context.Background(), createExtensionOpts)

	require.ErrorAs(t, err, &expected)
}
