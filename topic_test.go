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

const topicID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testTopicNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "topic %s not found."
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

	httpmock.RegisterResponder("GET", testClient.Endpoint+TopicsURI,
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

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+TopicsURI+"/"+topicID,
		httpmock.NewStringResponder(200, testTopicResponse))

	actual, err := testClient.Topic(context.Background(), topicID)

	require.NoError(t, err)
	assert.Equal(t, TopicExpected, actual)
}

func TestTopicNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testTopicNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+TopicsURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("topic %s not found.", NotFoundEntityID)

	_, err := testClient.Topic(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateToopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+TopicsURI,
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
	require.NoError(t, err)
	assert.Equal(t, TopicExpected, actual)
}

func TestCreateTopicInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+TopicsURI,
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

	require.ErrorAs(t, err, &expected)
}

func TestUpdateTopic(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+TopicsURI+"/"+topicID,
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
	require.NoError(t, err)
	assert.Equal(t, TopicUpdateExpected, actual)
}

func TestUpdateTopicInvalidPartitions(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+TopicsURI+"/"+topicID,
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

	require.ErrorAs(t, err, &expected)
}
