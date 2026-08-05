package v2

import (
	"net/http"
	"time"

	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/clickhouse"
)

const (
	// appName specifies an application name.
	appName = "dbaas-go"

	// appVersion specifies an application version.
	appVersion = "0.2.0"

	// userAgent contains a basic user agent that will be used in queries.
	userAgent = appName + "/" + appVersion
)

// API is the main manager for DBAAS v2 resources
type API struct {
	// Flavors
	// DatastoreTypes

	ClickHouse *clickhouse.API
	// Opensearch

}

// Public retry config
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func newAPIWithClient(client transport.Client) *API {
	return &API{
		ClickHouse: clickhouse.NewAPI(client),
	}
}

func NewAPI(token, endpoint string) (*API, error) {

	client := transport.NewHTTPClient(http.DefaultClient, token, endpoint, userAgent)

	return newAPIWithClient(client), nil
}

func NewAPIWithRetry(token, endpoint string, retryConfig RetryConfig) (*API, error) {

	client := transport.NewHTTPClient(
		http.DefaultClient,
		token,
		endpoint,
		userAgent,
		transport.WithRetry(
			transport.RetryConfig{
				MaxRetries:     retryConfig.MaxRetries,
				InitialBackoff: retryConfig.InitialBackoff,
				MaxBackoff:     retryConfig.MaxBackoff,
			},
		),
	)

	return newAPIWithClient(client), nil
}
