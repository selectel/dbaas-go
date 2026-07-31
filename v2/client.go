package v2

import (
	"net/http"

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

type API struct {
	// Flavors
	// DatastoreTypes

	ClickHouse *clickhouse.API
	// Opensearch

}

func NewAPI(token, endpoint string) (*API, error) {

	client := transport.NewHTTPClient(http.DefaultClient, token, endpoint, userAgent)

	return NewAPIWithClient(client), nil
}

func NewAPIWithClient(client transport.Client) *API {
	return &API{
		ClickHouse: clickhouse.NewAPI(client),
	}
}
