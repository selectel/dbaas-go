package clickhouse

import (
	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
	"github.com/selectel/dbaas-go/v2/internal"
)

type API struct {
	Datastore *DatastoreService
	// NodeGroup
	// ShardGroup
}

func NewAPI(c transport.Client) *API {
	engine := internal.NewEngineService(c, common.EngineClickHouse)

	datastoreSvc := &DatastoreService{engine}

	return &API{
		Datastore: datastoreSvc,
	}
}
