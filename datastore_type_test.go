package dbaas

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const datastoreTypeID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testDatastoreTypeNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "datastoretype %s not found."
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

	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoreTypesURI,
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

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestDatastoreType(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoreTypesURI+"/"+datastoreTypeID,
		httpmock.NewStringResponder(200, testDatastoreTypeResponse))

	expected := DatastoreType{
		ID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Engine:  "mysql",
		Version: "8",
	}

	actual, err := testClient.DatastoreType(context.Background(), datastoreTypeID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestDatastoreTypeNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testDatastoreTypeNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoreTypesURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("datastoretype %s not found.", NotFoundEntityID)

	_, err := testClient.DatastoreType(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}
