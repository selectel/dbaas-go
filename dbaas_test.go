package dbaas

func SetupTestClient() *API {
	testClient, _ := NewDBAASClient("test-token", "http://localhost/v1")
	return testClient
}
