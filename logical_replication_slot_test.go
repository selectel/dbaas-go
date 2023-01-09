package dbaas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const slotID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testLogicalReplicationSlotNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "logicalreplicationslot 123 not found."
	}
}`

const testLogicalReplicationSlotResponse = `
{
	"logical-replication-slot": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"name": "test_slot",
		"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"status": "ACTIVE"
	}
}`

const testLogicalReplicationSlotsResponse = `
{
	"logical-replication-slots": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"name": "test_slot",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"name": "test_slot-123",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"database_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"status": "ACTIVE"
		}
	]
}`

func TestLogicalReplicationSlots(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/logical-replication-slots",
		httpmock.NewStringResponder(200, testLogicalReplicationSlotsResponse))

	expected := []LogicalReplicationSlot{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			Name:        "test_slot",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			Name:        "test_slot-123",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Status:      StatusActive,
		},
	}

	actual, err := testClient.LogicalReplicationSlots(context.Background(), nil)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestLogicalReplicationSlot(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/logical-replication-slots/"+slotID,
		httpmock.NewStringResponder(200, testLogicalReplicationSlotResponse))

	expected := LogicalReplicationSlot{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		Name:        "test_slot",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusActive,
	}

	actual, err := testClient.LogicalReplicationSlot(context.Background(), slotID)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestLogicalReplicationSlotNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/logical-replication-slots/123",
		httpmock.NewStringResponder(404, testLogicalReplicationSlotNotFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "logicalreplicationslot 123 not found."

	_, err := testClient.LogicalReplicationSlot(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}

func TestCreateLogicalReplicationSlot(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	expected := LogicalReplicationSlot{
		ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:   "1970-01-01T00:00:00",
		UpdatedAt:   "1970-01-01T00:00:00",
		Name:        "test_slot",
		ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Status:      StatusPendingCreate,
	}

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/logical-replication-slots",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&LogicalReplicationSlotCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			slots := make(map[string]LogicalReplicationSlot)
			slots["logical-replication-slot"] = expected

			resp, err := httpmock.NewJsonResponse(200, slots)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createLogicalReplicationSlotOpts := LogicalReplicationSlotCreateOpts{
		Name:        "test_slot",
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatabaseID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	}

	actual, err := testClient.CreateLogicalReplicationSlot(context.Background(), createLogicalReplicationSlotOpts)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}
