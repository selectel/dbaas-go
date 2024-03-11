package dbaas

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const testInstanceNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "instance %s not found."
	}
}`

const instanceID = "7d959e48-4d42-41cd-8515-aae0a8482d8d"

func TestCreateFloatingIP(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+FloatingIPsURI,
		httpmock.NewStringResponder(200, ""))

	createFloatingIPOpts := FloatingIPsOpts{
		InstanceID: instanceID,
	}

	err := testClient.CreateFloatingIP(context.Background(), createFloatingIPOpts)

	require.NoError(t, err)
}

func TestDeleteFloatingIP(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", testClient.Endpoint+FloatingIPsURI,
		httpmock.NewStringResponder(200, ""))

	deleteFloatingIPOpts := FloatingIPsOpts{
		InstanceID: instanceID,
	}

	err := testClient.DeleteFloatingIP(context.Background(), deleteFloatingIPOpts)

	require.NoError(t, err)
}

func TestCreateFloatingIPInstanceNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", testClient.Endpoint+FloatingIPsURI,
		httpmock.NewStringResponder(404, testInstanceNotFoundResponse))

	createFloatingIPOpts := FloatingIPsOpts{
		InstanceID: instanceID,
	}

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("instance %s not found.", instanceID)

	err := testClient.DeleteFloatingIP(context.Background(), createFloatingIPOpts)

	require.ErrorAs(t, err, &expected)
}
