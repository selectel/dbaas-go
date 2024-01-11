package dbaas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const topicID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testTopicNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "topic 123 not found."
	}
}`

const testCreateTopicInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": 
			"Validation failure: {'topic.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testUpdateTopicInvalidPartitionsResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message":
		"Validation failure: {'topic.partitions': \"'4001 is greater than the maximum of 4000'\"}"
	}
}`

const testTopicResponse = `{
	"topic": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "topic1",
		"partitions": 1,
		"status": "ACTIVE"
	}
}`

const testTopicsResponse = `
{
	"topics": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "topic1",
			"partitions": 1,
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "topic2",
			"partitions": 2,
			"status": "ACTIVE"
		}
	]
}`

var TopicResponse = Topic{ //nolint
	ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:   "1970-01-01T00:00:00",
	UpdatedAt:   "1970-01-01T00:00:00",
	ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	Name:        "topic1",
	Partitions:  1,
	Status:      StatusActive,
}

var TopicExpected Topic = Topic{ //nolint
	ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:   "1970-01-01T00:00:00",
	UpdatedAt:   "1970-01-01T00:00:00",
	ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	Name:        "topic1",
	Partitions:  1,
	Status:      StatusActive,
}

func TestTopics(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/topics",
		httpmock.NewStringResponder(200, testTopicsResponse))

	expected := []Topic{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:        "topic1",
			Partitions:  1,
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:        "topic2",
			Partitions:  2,
			Status:      StatusActive,
		},
	}

	actual, err := testClient.Topics(context.Background(), nil)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestTopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/topics/"+topicID,
		httpmock.NewStringResponder(200, testTopicResponse))

	actual, err := testClient.Topic(context.Background(), topicID)

	if assert.NoError(t, err) {
		assert.Equal(t, TopicExpected, actual)
	}
}

func TestTopicNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/topics/123",
		httpmock.NewStringResponder(404, testTopicNotFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "topic 123 not found."

	_, err := testClient.Topic(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}

func TestCreateToopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/topics",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&TopicCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			topics := make(map[string]Topic)
			TopicCreateResponse := TopicResponse
			TopicCreateResponse.Status = StatusPendingCreate
			topics["topic"] = TopicResponse

			resp, err := httpmock.NewJsonResponse(200, topics)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createTopicOpts := TopicCreateOpts{
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "topic1",
		Partitions:  1,
	}

	actual, err := testClient.CreateTopic(context.Background(), createTopicOpts)

	TopicCreateExpexted := TopicExpected
	TopicCreateExpexted.Status = StatusPendingCreate
	if assert.NoError(t, err) {
		assert.Equal(t, TopicExpected, actual)
	}
}

func TestCreateTopicInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+"/topics",
		httpmock.NewStringResponder(400, testCreateTopicInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'topic.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createTopicOpts := TopicCreateOpts{
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:        "topic1",
		Partitions:  1,
	}

	_, err := testClient.CreateTopic(context.Background(), createTopicOpts)

	assert.ErrorAs(t, err, &expected)
}

func TestUpdateTopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+"/topics/"+topicID,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&TopicUpdateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			topics := make(map[string]Topic)
			TopicUpdateResponse := TopicResponse
			TopicUpdateResponse.Status = StatusPendingUpdate
			topics["topic"] = TopicUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, topics)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	updateTopicOpts := TopicUpdateOpts{
		Partitions: 2,
	}

	actual, err := testClient.UpdateTopic(context.Background(), topicID, updateTopicOpts)

	TopicUpdateExpected := TopicExpected
	TopicUpdateExpected.Status = StatusPendingUpdate
	if assert.NoError(t, err) {
		assert.Equal(t, TopicUpdateExpected, actual)
	}
}

func TestUpdateTopicInvalidPartitions(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+"/topics/"+topicID,
		httpmock.NewStringResponder(400, testUpdateTopicInvalidPartitionsResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'topic.partitions': \"'4001 is greater than the maximum of 4000'\"}`

	updateTopicOpts := TopicUpdateOpts{
		Partitions: 4001,
	}

	_, err := testClient.UpdateTopic(context.Background(), topicID, updateTopicOpts)

	assert.ErrorAs(t, err, &expected)
}
