package dbaas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const databaseID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testDatabaseNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "database 123 not found."
	}
}`

const testCreateDatabaseInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": 
			"Validation failure: {'database.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testUpdateDatabaseInvalidOwnerIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message":
			"Validation failure: {'database.owner_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testDatabaseResponse = `{
	"database": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"name": "db",
		"owner_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"lc_collate": "C",
		"lc_ctype": "C",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"status": "ACTIVE"
	}
}`

const testDatabasesResponse = `
{
	"databases": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"name": "db",
			"owner_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"lc_collate": "C",
			"lc_ctype": "C",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"name": "db123",
			"owner_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"lc_collate": "ru_RU.utf8",
			"lc_ctype": "ru_RU.utf8",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		}
	]
}`

func TestDatabases(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/databases",
		httpmock.NewStringResponder(200, testDatabasesResponse))

	expected := []Database{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			Name:        "db",
			OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			LcCollate:   "C",
			LcCtype:     "C",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			Name:        "db123",
			OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			LcCollate:   "ru_RU.utf8",
			LcCtype:     "ru_RU.utf8",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:      StatusActive,
		},
	}

	actual, err := testClient.Databases(context.Background(), nil)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestDatabase(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/databases/"+databaseID,
		httpmock.NewStringResponder(200, testDatabaseResponse))

	expected := Database{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		Name:        "db",
		OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		LcCollate:   "C",
		LcCtype:     "C",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusActive,
	}

	actual, err := testClient.Database(context.Background(), databaseID)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestDatabaseNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/databases/123",
		httpmock.NewStringResponder(404, testDatabaseNotFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "database 123 not found."

	_, err := testClient.Database(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}

func TestCreateDatabase(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/databases",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatabaseCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			databases := make(map[string]Database)
			databases["database"] = Database{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				Name:        "db",
				OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				LcCollate:   "C",
				LcCtype:     "C",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Status:      StatusPendingCreate,
			}

			resp, err := httpmock.NewJsonResponse(200, databases)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createDatabaseOpts := DatabaseCreateOpts{
		Name:        "db",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	expected := Database{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		Name:        "db",
		OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		LcCollate:   "C",
		LcCtype:     "C",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusPendingCreate,
	}

	actual, err := testClient.CreateDatabase(context.Background(), createDatabaseOpts)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestCreateDatabaseInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/databases",
		httpmock.NewStringResponder(400, testCreateDatabaseInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'database.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createDatabaseOpts := DatabaseCreateOpts{
		Name:        "db",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f",
		OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	_, err := testClient.CreateDatabase(context.Background(), createDatabaseOpts)

	assert.ErrorAs(t, err, &expected)
}

func TestUpdateDatabase(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+"/databases/"+databaseID,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatabaseUpdateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			databases := make(map[string]Database)
			databases["database"] = Database{
				ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				CreatedAt:   "1970-01-01T00:00:00",
				UpdatedAt:   "1970-01-01T00:00:00",
				Name:        "db",
				OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				LcCollate:   "C",
				LcCtype:     "C",
				DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				Status:      StatusPendingUpdate,
			}

			resp, err := httpmock.NewJsonResponse(200, databases)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	updateDatabaseOpts := DatabaseUpdateOpts{
		OwnerID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
	}

	expected := Database{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		Name:        "db",
		OwnerID:     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		LcCollate:   "C",
		LcCtype:     "C",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusPendingUpdate,
	}

	actual, err := testClient.UpdateDatabase(context.Background(), databaseID, updateDatabaseOpts)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestUpdateDatabaseInvalidOwnerID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+"/databases/"+databaseID,
		httpmock.NewStringResponder(400, testUpdateDatabaseInvalidOwnerIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'database.owner_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	updateDatabaseOpts := DatabaseUpdateOpts{
		OwnerID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f",
	}

	_, err := testClient.UpdateDatabase(context.Background(), databaseID, updateDatabaseOpts)

	assert.ErrorAs(t, err, &expected)
}
