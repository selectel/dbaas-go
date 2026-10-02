package dbaas

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRolesResponse = `{
	"roles": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "pg_read_all_data"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "pg_write_all_data"
		}
	]
}`

func TestRoles(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+RolesURI,
		httpmock.NewStringResponder(200, testRolesResponse))

	expected := []Roles{
		{
			ID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreTypeID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:            "pg_read_all_data",
		},
		{
			ID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreTypeID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:            "pg_write_all_data",
		},
	}

	actual, err := testClient.Roles(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
