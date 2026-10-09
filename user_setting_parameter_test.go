package dbaas

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const userSettingParameterID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testUserSettingParameterNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "usersettingparameter %s not found."
	}
}`

const testUserSettingParametersResponse = `{
	"user-setting-parameters": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"name": "statement_timeout",
			"type": "int",
			"choices": null,
			"min": 0,
			"max": 2147483647,
			"default_value": "0",
			"unit": "ms",
			"apply_mechanism": "guc",
			"is_available_for_customer": true,
			"is_changeable": true,
			"can_be_empty": true
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"name": "default_transaction_isolation",
			"type": "str",
			"choices": [
				"serializable",
				"repeatable read",
				"read committed",
				"read uncommitted"
			],
			"min": null,
			"max": null,
			"default_value": "read committed",
			"unit": "",
			"apply_mechanism": "guc",
			"is_available_for_customer": true,
			"is_changeable": true,
			"can_be_empty": true
		}
	]
}`

const testUserSettingParameterResponse = `{
	"user-setting-parameter": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"name": "login",
		"type": "bool",
		"choices": null,
		"min": null,
		"max": null,
		"default_value": "true",
		"unit": "",
		"apply_mechanism": "role_attr",
		"is_available_for_customer": true,
		"is_changeable": true,
		"can_be_empty": true
	}
}`

func TestUserSettingParameters(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+UserSettingParametersURI,
		httpmock.NewStringResponder(200, testUserSettingParametersResponse))

	choices := []any{"serializable", "repeatable read", "read committed", "read uncommitted"}
	expected := []UserSettingParameter{
		{
			ID:                     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Name:                   "statement_timeout",
			Type:                   "int",
			Unit:                   "ms",
			ApplyMechanism:         "guc",
			Min:                    0.0,
			Max:                    2147483647.0,
			DefaultValue:           "0",
			Choices:                nil,
			IsAvailableForCustomer: true,
			IsChangeable:           true,
			CanBeEmpty:             true,
		},
		{
			ID:                     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			Name:                   "default_transaction_isolation",
			Type:                   "str",
			Unit:                   "",
			ApplyMechanism:         "guc",
			Min:                    nil,
			Max:                    nil,
			DefaultValue:           "read committed",
			Choices:                choices,
			IsAvailableForCustomer: true,
			IsChangeable:           true,
			CanBeEmpty:             true,
		},
	}

	actual, err := testClient.UserSettingParameters(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUserSettingParametersWithQuery(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "=~/user-setting-parameters",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4", req.URL.Query().Get("datastore_type_id"))
			assert.Equal(t, "statement_timeout", req.URL.Query().Get("name"))
			return httpmock.NewStringResponse(200, testUserSettingParametersResponse), nil
		})

	params := &UserSettingParameterQueryParams{
		DatastoreTypeID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:            "statement_timeout",
	}
	actual, err := testClient.UserSettingParameters(context.Background(), params)

	require.NoError(t, err)
	assert.Len(t, actual, 2)
}

func TestUserSettingParameter(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+UserSettingParametersURI+"/"+userSettingParameterID,
		httpmock.NewStringResponder(200, testUserSettingParameterResponse))

	expected := UserSettingParameter{
		ID:                     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Name:                   "login",
		Type:                   "bool",
		Unit:                   "",
		ApplyMechanism:         "role_attr",
		Min:                    nil,
		Max:                    nil,
		DefaultValue:           "true",
		Choices:                nil,
		IsAvailableForCustomer: true,
		IsChangeable:           true,
		CanBeEmpty:             true,
	}

	actual, err := testClient.UserSettingParameter(context.Background(), userSettingParameterID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUserSettingParameterNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testUserSettingParameterNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+UserSettingParametersURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("usersettingparameter %s not found.", NotFoundEntityID)

	_, err := testClient.UserSettingParameter(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}
