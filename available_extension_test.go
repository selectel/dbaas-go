package dbaas

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const availableExtensionID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testAvailableExtensionNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "availableextension 123 not found."
	}
}`

const testAvailableExtensionsResponse = `{
	"available-extensions": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "xml2",
			"datastore_type_ids": [
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
			],
			"dependency_ids": []
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"name": "prefix",
			"datastore_type_ids": [
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
			],
			"dependency_ids": []
		}
	]
}`

const testAvailableExtensionResponse = `{
	"available-extension": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "xml2",
		"datastore_type_ids": [
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
		],
		"dependency_ids": []
	}
}`

func TestAvailableExtensions(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/available-extensions",
		httpmock.NewStringResponder(200, testAvailableExtensionsResponse))

	expected := []AvailableExtension{
		{
			ID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name: "xml2",
			DatastoreTypeIDs: []string{
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			},
			DependencyIDs: []string{},
		},
		{
			ID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			Name: "prefix",
			DatastoreTypeIDs: []string{
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			},
			DependencyIDs: []string{},
		},
	}

	actual, err := testClient.AvailableExtensions(context.Background())

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestAvailableExtension(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/available-extensions/"+datastoreTypeID,
		httpmock.NewStringResponder(200, testAvailableExtensionResponse))

	expected := AvailableExtension{
		ID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name: "xml2",
		DatastoreTypeIDs: []string{
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
		},
		DependencyIDs: []string{},
	}

	actual, err := testClient.AvailableExtension(context.Background(), availableExtensionID)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestAvailableExtensionNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/available-extensions/123",
		httpmock.NewStringResponder(404, testAvailableExtensionNotFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "availableextension 123 not found."

	_, err := testClient.AvailableExtension(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}
