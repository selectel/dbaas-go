package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTransport(serverURL string) *HTTPClient {
	return NewHTTPClient(
		http.DefaultClient,
		"test-token",
		serverURL,
		"dbaas-go/test",
	)
}

func TestDo_GET(t *testing.T) {
	var (
		gotMethod string
		gotPath   string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write([]byte(`{
			"id":"123",
			"name":"cluster"
		}`))
	}))
	defer server.Close()

	client := newTransport(server.URL)

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/datastores/clickhouse",
		nil,
		&resp,
	)

	require.NoError(t, err)

	require.Equal(t, http.MethodGet, gotMethod)
	require.Equal(t, "/datastores/clickhouse", gotPath)

	require.Equal(t, "123", resp.ID)
	require.Equal(t, "cluster", resp.Name)
}

func TestDo_POSTBody(t *testing.T) {
	var (
		gotMethod string
		gotPath   string
		body      []byte
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)

		gotMethod = r.Method
		gotPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTransport(server.URL)

	req := struct {
		Name string `json:"name"`
	}{
		Name: "cluster",
	}

	err := client.Do(
		context.Background(),
		http.MethodPost,
		"/datastores/clickhouse",
		req,
		nil,
	)

	require.NoError(t, err)

	require.JSONEq(t, `{"name":"cluster"}`, string(body))

	require.Equal(t, http.MethodPost, gotMethod)
	require.Equal(t, "/datastores/clickhouse", gotPath)
}

func TestDo_Headers(t *testing.T) {
	var (
		token     string
		userAgent string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token = r.Header.Get("X-Auth-Token")
		userAgent = r.Header.Get("User-Agent")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTransport(server.URL)

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/datastores/clickhouse",
		nil,
		nil,
	)

	require.NoError(t, err)

	require.Equal(t, "test-token", token)
	require.Equal(t, "dbaas-go/test", userAgent)
}

func TestDo_APIError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte(`
		{
			"error": {
				"message": "validation failed",
				"title": "invalid_request",
				"code": 400
			}
		}
		`))
	}))
	defer server.Close()

	client := newTransport(server.URL)

	err := client.Do(
		context.Background(),
		http.MethodGet,
		"/datastores/clickhouse",
		nil,
		nil,
	)

	require.Error(t, err)

	var apiErr *DBaaSAPIError
	require.ErrorAs(t, err, &apiErr)

	require.Equal(t, "validation failed", apiErr.APIError.Message)
}
