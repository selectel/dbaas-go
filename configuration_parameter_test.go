package dbaas

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const configurationParameterID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testConfigurationParameterNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "configurationparameter 123 not found."
	}
}`

const testConfigurationParametersResponse = `{
	"configuration-parameters": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "temp_file_limit",
			"type": "int",
			"choices": null,
			"min": -1,
			"max": 2147483647,
			"default_value": -1,
			"unit": "kB",
			"is_restart_required": false,
			"is_changeable": true
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "concurrent_insert",
			"type": "str",
			"choices": [
				"NEVER",
				"AUTO",
				"ALWAYS",
				"0",
				"1",
				"2"
			],
			"min": null,
			"max": null,
			"default_value": "AUTO",
			"unit": "",
			"is_restart_required": false,
			"is_changeable": true
		}
	]
}`

const testConfigurationParameterResponse = `{
	"configuration-parameter": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "thread_pool_size",
		"type": "int",
		"choices": null,
		"min": 1,
		"max": 64,
		"default_value": null,
		"unit": "",
		"is_restart_required": true,
		"is_changeable": true
	}
}`

func TestConfigurationParameters(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/configuration-parameters",
		httpmock.NewStringResponder(200, testConfigurationParametersResponse))

	choices := []interface{}{"NEVER", "AUTO", "ALWAYS", "0", "1", "2"}
	expected := []ConfigurationParameter{
		{
			ID:                "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreTypeID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:              "temp_file_limit",
			Type:              "int",
			Unit:              "kB",
			Min:               -1.0,
			Max:               2147483647.0,
			DefaultValue:      -1.0,
			Choices:           nil,
			IsRestartRequired: false,
			IsChangeable:      true,
		},
		{
			ID:                "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreTypeID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:              "concurrent_insert",
			Type:              "str",
			Unit:              "",
			Min:               nil,
			Max:               nil,
			DefaultValue:      "AUTO",
			Choices:           choices,
			IsRestartRequired: false,
			IsChangeable:      true,
		},
	}

	actual, err := testClient.ConfigurationParameters(context.Background())

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestConfigurationParameter(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/configuration-parameters/"+configurationParameterID,
		httpmock.NewStringResponder(200, testConfigurationParameterResponse))

	expected := ConfigurationParameter{
		ID:                "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		DatastoreTypeID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:              "thread_pool_size",
		Type:              "int",
		Unit:              "",
		Min:               1.0,
		Max:               64.0,
		DefaultValue:      nil,
		Choices:           nil,
		IsRestartRequired: true,
		IsChangeable:      true,
	}

	actual, err := testClient.ConfigurationParameter(context.Background(), configurationParameterID)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestConfigurationParameterNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+"/configuration-parameters/123",
		httpmock.NewStringResponder(404, testConfigurationParameterNotFoundResponse))

	expected := &APIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = "configurationparameter 123 not found."

	_, err := testClient.ConfigurationParameter(context.Background(), "123")

	assert.ErrorAs(t, err, &expected)
}
