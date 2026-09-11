package opensearch

import (
	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
)

// API is api client to manage opensearch resources.
type API struct {
	*DatastoreService
	*NodeGroupService
	// ShardGroup
}

func NewAPI(c transport.Client) *API {
	engine := common.NewEngineService(c, common.EngineOpenSearch)

	datastoreSvc := &DatastoreService{engine}
	nodeGroupSvc := &NodeGroupService{engine}

	return &API{
		datastoreSvc,
		nodeGroupSvc,
	}
}
