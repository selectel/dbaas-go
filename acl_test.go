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

const aclID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

const testACLNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "acl %s not found."
	}
}`

const testCreateACLInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": 
			"Validation failure: {'acl.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testUpdateACLInvalidResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message":
		"Validation failure: At least one of these fields (allow_read, allow_write) must be true"
	}
}`

const testACLResponse = `{
	"acl": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"pattern": "topic1",
		"pattern_type": "literal",
		"allow_read": true,
		"allow_write": true,
		"status": "ACTIVE"
	}
}`

const testACLsResponse = `
{
	"acls": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"pattern": "topic1",
			"pattern_type": "literal",
			"allow_read": true,
			"allow_write": true,
			"status": "ACTIVE"
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"datastore_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"user_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"pattern": "topic2",
			"pattern_type": "literal",
			"allow_read": true,
			"allow_write": true,
			"status": "ACTIVE"
		}
	]
}`

var ACLResponse ACL = ACL{ //nolint
	ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:   "1970-01-01T00:00:00",
	UpdatedAt:   "1970-01-01T00:00:00",
	ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	Pattern:     "topic1",
	PatternType: "literal",
	AllowRead:   true,
	AllowWrite:  true,
	Status:      StatusActive,
}

var ACLExpected ACL = ACL{ //nolint
	ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:   "1970-01-01T00:00:00",
	UpdatedAt:   "1970-01-01T00:00:00",
	ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	Pattern:     "topic1",
	PatternType: "literal",
	AllowRead:   true,
	AllowWrite:  true,
	Status:      StatusActive,
}

func TestACLs(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+ACLsURI,
		httpmock.NewStringResponder(200, testACLsResponse))

	expected := []ACL{
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Pattern:     "topic1",
			PatternType: "literal",
			AllowRead:   true,
			AllowWrite:  true,
			Status:      StatusActive,
		},
		{
			ID:          "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			CreatedAt:   "1970-01-01T00:00:00",
			UpdatedAt:   "1970-01-01T00:00:00",
			ProjectID:   "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			Pattern:     "topic2",
			PatternType: "literal",
			AllowRead:   true,
			AllowWrite:  true,
			Status:      StatusActive,
		},
	}

	actual, err := testClient.ACLs(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestACL(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+ACLsURI+"/"+aclID,
		httpmock.NewStringResponder(200, testACLResponse))

	actual, err := testClient.ACL(context.Background(), aclID)

	require.NoError(t, err)
	assert.Equal(t, ACLExpected, actual)
}

func TestACLNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testACLNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+ACLsURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("acl %s not found.", NotFoundEntityID)

	_, err := testClient.ACL(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateACL(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+ACLsURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&ACLCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			acls := make(map[string]ACL)
			ACLCreateResponse := ACLResponse
			ACLCreateResponse.Status = StatusPendingCreate
			acls["acl"] = ACLCreateResponse

			resp, err := httpmock.NewJsonResponse(200, acls)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	createACLOpts := ACLCreateOpts{
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Pattern:     "topic1",
		PatternType: "literal",
		AllowRead:   true,
		AllowWrite:  true,
	}

	actual, err := testClient.CreateACL(context.Background(), createACLOpts)

	ACLCreateExpected := ACLExpected
	ACLCreateExpected.Status = StatusPendingCreate
	require.NoError(t, err)
	assert.Equal(t, ACLCreateExpected, actual)
}

func TestCreateACLInvalidDatastoreID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+ACLsURI,
		httpmock.NewStringResponder(400, testCreateACLInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		{'acl.datastore_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createACLOpts := ACLCreateOpts{
		DatastoreID: "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		UserID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Pattern:     "topic1",
		PatternType: "literal",
		AllowRead:   true,
		AllowWrite:  true,
	}

	_, err := testClient.CreateACL(context.Background(), createACLOpts)

	require.ErrorAs(t, err, &expected)
}

func TestUpdateACL(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+ACLsURI+"/"+aclID,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&ACLUpdateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			acls := make(map[string]ACL)
			ACLUpdateResponse := ACLResponse
			ACLUpdateResponse.Status = StatusPendingUpdate
			acls["acl"] = ACLUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, acls)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}

			return resp, nil
		})

	updateACLOpts := ACLUpdateOpts{
		AllowRead:  true,
		AllowWrite: false,
	}

	actual, err := testClient.UpdateACL(context.Background(), aclID, updateACLOpts)

	ACLUpdateExpexted := ACLExpected
	ACLUpdateExpexted.Status = StatusPendingUpdate
	require.NoError(t, err)
	assert.Equal(t, ACLUpdateExpexted, actual)
}

func TestUpdateACLInvalidResponse(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+ACLsURI+"/"+aclID,
		httpmock.NewStringResponder(400, testUpdateACLInvalidResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: 
		At least one of these fields (allow_read, allow_write) must be true`

	updateACLOpts := ACLUpdateOpts{
		AllowRead:  false,
		AllowWrite: false,
	}

	_, err := testClient.UpdateACL(context.Background(), aclID, updateACLOpts)

	require.ErrorAs(t, err, &expected)
}
