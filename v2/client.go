package v2

import (
	"fmt"
	"net/http"
	"time"

	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/clickhouse"
	"github.com/selectel/dbaas-go/v2/common"
)

// API is the main manager for DBAAS v2 resources.
type API struct {
	Flavor        *common.FlavorService
	DatastoreType *common.DatastoreTypeService

	ClickHouse *clickhouse.API
	// Opensearch
}

// Public retry config.
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func newAPIWithClient(client transport.Client) *API {
	return &API{
		Flavor:        common.NewFlavorService(client),
		DatastoreType: common.NewDatastoreTypeService(client),
		ClickHouse:    clickhouse.NewAPI(client),
	}
}

func NewAPI(token, endpoint string) (*API, error) {
	client, err := transport.NewHTTPClient(http.DefaultClient, token, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transport client: %w", err)
	}

	return newAPIWithClient(client), nil
}

func NewAPIWithRetry(token, endpoint string, retryConfig RetryConfig) (*API, error) {
	client, err := transport.NewHTTPClient(
		http.DefaultClient,
		token,
		endpoint,
		transport.WithRetry(
			transport.RetryConfig{
				MaxRetries:     retryConfig.MaxRetries,
				InitialBackoff: retryConfig.InitialBackoff,
				MaxBackoff:     retryConfig.MaxBackoff,
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transport client: %w", err)
	}

	return newAPIWithClient(client), nil
}
