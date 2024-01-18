package dbaas

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const flavorID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testFlavorNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "flavor %s not found."
	}
}`

const testFlavorsResponse = `{
	"flavors": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "flavor-2",
			"description": "",
			"vcpus": 2,
			"ram": 4096,
			"disk": 10,
			"datastore_type_ids": [
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
			],
			"fl_size": "standard"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"name": "flavor-3",
			"description": "",
			"vcpus": 4,
			"ram": 8192,
			"disk": 100,
			"datastore_type_ids": [
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
			],
			"fl_size": "standard"
		}
	]
}`

const testFlavorResponse = `{
	"flavor": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "flavor-2",
		"description": "",
		"vcpus": 2,
		"ram": 4096,
		"disk": 10,
		"datastore_type_ids": [
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5"
		],
		"fl_size": "standard"
	}
}`

func TestFlavors(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+FlavorsURI,
		httpmock.NewStringResponder(200, testFlavorsResponse))

	expected := []FlavorResponse{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:        "flavor-2",
			Description: "",
			Vcpus:       2,
			RAM:         4096,
			Disk:        10,
			DatastoreTypeIDs: []string{
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			},
			FlSize: "standard",
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			Name:        "flavor-3",
			Description: "",
			Vcpus:       4,
			RAM:         8192,
			Disk:        100,
			DatastoreTypeIDs: []string{
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			},
			FlSize: "standard",
		},
	}

	actual, err := testClient.Flavors(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFlavor(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+FlavorsURI+"/"+flavorID,
		httpmock.NewStringResponder(200, testFlavorResponse))

	expected := FlavorResponse{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "flavor-2",
		Description: "",
		Vcpus:       2,
		RAM:         4096,
		Disk:        10,
		DatastoreTypeIDs: []string{
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
		},
		FlSize: "standard",
	}

	actual, err := testClient.Flavor(context.Background(), flavorID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFlavorNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testFlavorNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+FlavorsURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("flavor %s not found.", NotFoundEntityID)

	_, err := testClient.Flavor(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}
