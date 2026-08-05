package internal

import (
	"context"
	"net/http"

	"github.com/selectel/dbaas-go/internal/transport"
	"github.com/selectel/dbaas-go/v2/common"
)

// BaseService presents base service logic with helpers
type BaseService struct {
	client   transport.Client
	rootPath string
}

func (s *BaseService) Get(ctx context.Context, path string, resp any) error {
	return s.client.Do(ctx, http.MethodGet, s.rootPath+path, nil, resp)
}

func (s *BaseService) Post(ctx context.Context, path string, body any, resp any) error {
	return s.client.Do(ctx, http.MethodPost, s.rootPath+path, body, resp)
}

func (s *BaseService) Put(ctx context.Context, path string, body any, resp any) error {
	return s.client.Do(ctx, http.MethodPut, s.rootPath+path, body, resp)
}

func (s *BaseService) Delete(ctx context.Context, path string) error {
	return s.client.Do(ctx, http.MethodPut, s.rootPath+path, nil, nil)
}

// EngineService presents engine service for specific engine (datastore type)
type EngineService struct {
	*BaseService
	engine common.Engine
}

func NewEngineService(client transport.Client, engine common.Engine) *EngineService {

	root := "/datastores/" + string(engine)

	return &EngineService{
		BaseService: &BaseService{
			client:   client,
			rootPath: root,
		},
		engine: engine,
	}
}
