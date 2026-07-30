package clickhouse

import (
	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/internal"
)

type API struct {
	Datastore *DatastoreService
}

func NewAPI(c transport.Client) *API {
	engine := internal.EngineService{
		Client: c,
		Engine: "clickhouse",
	}
	datastoreSvc := &DatastoreService{engine}

	return &API{
		Datastore: datastoreSvc,
	}
}
