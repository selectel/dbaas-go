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
			"datastore_group_name": "postgresql",
			"name": "statement_timeout",
			"type": "int",
			"choices": null,
			"min": 0,
			"max": 2147483647,
			"unit": "ms",
			"apply_mechanism": "guc",
			"is_available_for_customer": true,
			"is_changeable": true
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"datastore_group_name": "postgresql",
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
			"unit": "",
			"apply_mechanism": "guc",
			"is_available_for_customer": true,
			"is_changeable": true
		}
	]
}`

const testUserSettingParameterResponse = `{
	"user-setting-parameter": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_group_name": "postgresql",
		"name": "login",
		"type": "bool",
		"choices": null,
		"min": null,
		"max": null,
		"unit": "",
		"apply_mechanism": "role_attr",
		"is_available_for_customer": true,
		"is_changeable": true
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
			DatastoreGroupName:     "postgresql",
			Name:                   "statement_timeout",
			Type:                   "int",
			Unit:                   "ms",
			ApplyMechanism:         "guc",
			Min:                    0.0,
			Max:                    2147483647.0,
			Choices:                nil,
			IsAvailableForCustomer: true,
			IsChangeable:           true,
		},
		{
			ID:                     "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			DatastoreGroupName:     "postgresql",
			Name:                   "default_transaction_isolation",
			Type:                   "str",
			Unit:                   "",
			ApplyMechanism:         "guc",
			Min:                    nil,
			Max:                    nil,
			Choices:                choices,
			IsAvailableForCustomer: true,
			IsChangeable:           true,
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
			assert.Equal(t, "postgresql", req.URL.Query().Get("datastore_group_name"))
			assert.Equal(t, "statement_timeout", req.URL.Query().Get("name"))
			return httpmock.NewStringResponse(200, testUserSettingParametersResponse), nil
		})

	params := &UserSettingParameterQueryParams{
		DatastoreGroupName: "postgresql",
		Name:               "statement_timeout",
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
		DatastoreGroupName:     "postgresql",
		Name:                   "login",
		Type:                   "bool",
		Unit:                   "",
		ApplyMechanism:         "role_attr",
		Min:                    nil,
		Max:                    nil,
		Choices:                nil,
		IsAvailableForCustomer: true,
		IsChangeable:           true,
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
