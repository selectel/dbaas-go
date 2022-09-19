package dbaas

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const datastoreTypeID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testDatastoreTypeNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "datastoretype 123 not found."
	}
}`

const testDatastoreTypesResponse = `{
	"datastore-types": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"engine": "mysql",
			"version": "8"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"engine": "postgresql",
			"version": "12"
		}
	]
}`

const testDatastoreTypeResponse = `{
	"datastore-type": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"engine": "mysql",
		"version": "8"
	}
}`

func TestDatastoreTypes(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/datastore-types",
		httpmock.NewStringResponder(200, testDatastoreTypesResponse))

	expected := []DatastoreType{
		{
			ID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Engine:  "mysql",
			Version: "8",
		},
		{
			ID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			Engine:  "postgresql",
			Version: "12",
		},
	}

	actual, err := testClient.DatastoreTypes(context.Background())

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestDatastoreType(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/datastore-types/"+datastoreTypeID,
		httpmock.NewStringResponder(200, testDatastoreTypeResponse))

	expected := DatastoreType{
		ID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Engine:  "mysql",
		Version: "8",
	}

	actual, err := testClient.DatastoreType(context.Background(), datastoreTypeID)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestDatastoreTypeNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/datastore-types/123",
		httpmock.NewStringResponder(404, testDatastoreTypeNotFoundResponse))

	expected := &APIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "datastoretype 123 not found."

	_, err := testClient.DatastoreType(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}
