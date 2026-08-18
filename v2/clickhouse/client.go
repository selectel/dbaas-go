package clickhouse

import (
	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
	"github.com/selectel/dbaas-go/v2/internal"
)

// API is api client to manage clickhouse resources.
type API struct {
	Datastore *DatastoreService
	NodeGroup *NodeGroupService
	// ShardGroup
}

func NewAPI(c transport.Client) *API {
	engine := internal.NewEngineService(c, common.EngineClickHouse)

	datastoreSvc := &DatastoreService{engine}
	nodeGroupSvc := &NodeGroupService{engine}

	return &API{
		Datastore: datastoreSvc,
		NodeGroup: nodeGroupSvc,
	}
}
