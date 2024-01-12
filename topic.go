package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Topic is the API response for the topics.
type Topic struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ProjectID   string `json:"project_id"`
	DatastoreID string `json:"datastore_id"`
	Name        string `json:"name"`
	Status      Status `json:"status"`
	Partitions  uint16 `json:"partitions"`
}

// TopicCreateOpts represents options for the topic Create request.
type TopicCreateOpts struct {
	DatastoreID string `json:"datastore_id"`
	Name        string `json:"name"`
	Partitions  uint16 `json:"partitions"`
}

// TopicUpdateOpts represents options for the topic Update request.
type TopicUpdateOpts struct {
	Partitions uint16 `json:"partitions"`
}

// TopicQueryParams represents available query parameters for the topic.
type TopicQueryParams struct {
	ID          string `json:"id,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	DatastoreID string `json:"datastore_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      Status `json:"status,omitempty"`
}

const TopicsURI = "/topics"

// Topics returns all topics.
func (api *API) Topics(ctx context.Context, params *TopicQueryParams) ([]Topic, error) {
	uri, err := setQueryParams(TopicsURI, params)
	if err != nil {
		return []Topic{}, err
	}

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Topic{}, err
	}

	var result struct {
		Topics []Topic `json:"topics"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Topic{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Topics, nil
}

// Topic returns a topic based on the ID.
func (api *API) Topic(ctx context.Context, topicID string) (Topic, error) {
	uri := fmt.Sprintf("%s/%s", TopicsURI, topicID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return Topic{}, err
	}

	var result struct {
		Topic Topic `json:"topic"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Topic{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Topic, nil
}

// CreateTopic creates a new topic.
func (api *API) CreateTopic(ctx context.Context, opts TopicCreateOpts) (Topic, error) {
	createTopicOpts := struct {
		Topic TopicCreateOpts `json:"topic"`
	}{
		Topic: opts,
	}
	requestBody, err := json.Marshal(createTopicOpts)
	if err != nil {
		return Topic{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, TopicsURI, requestBody)
	if err != nil {
		return Topic{}, err
	}

	var result struct {
		Topic Topic `json:"topic"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Topic{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Topic, nil
}

// UpdateTopic updates an existing topic.
func (api *API) UpdateTopic(ctx context.Context, topicID string, opts TopicUpdateOpts) (Topic, error) {
	uri := fmt.Sprintf("%s/%s", TopicsURI, topicID)
	updateTopicOpts := struct {
		Topic TopicUpdateOpts `json:"topic"`
	}{
		Topic: opts,
	}
	requestBody, err := json.Marshal(updateTopicOpts)
	if err != nil {
		return Topic{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Topic{}, err
	}

	var result struct {
		Topic Topic `json:"topic"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Topic{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Topic, nil
}

// DeleteTopic deletes an existing topic.
func (api *API) DeleteTopic(ctx context.Context, topicID string) error {
	uri := fmt.Sprintf("%s/%s", TopicsURI, topicID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}
