package clickhouse

import (
	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
)

// API is api client to manage clickhouse resources.
type API struct {
	*DatastoreService
	*NodeGroupService
	*ShardGroupService
}

func NewAPI(c transport.Client) *API {
	engine := common.NewEngineService(c, common.EngineClickHouse)

	datastoreSvc := &DatastoreService{engine}
	nodeGroupSvc := &NodeGroupService{engine}
	shardGroupSvc := &ShardGroupService{engine}

	return &API{
		datastoreSvc,
		nodeGroupSvc,
		shardGroupSvc,
	}
}
