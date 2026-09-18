package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/selectel/dbaas-go/internal/transport"
)

// baseService presents base service logic with helpers.
type baseService struct {
	client   transport.Client
	rootPath string
}

func (s *baseService) Get(ctx context.Context, path string, resp any) error {
	err := s.client.Do(ctx, http.MethodGet, s.rootPath+path, nil, resp)
	if err != nil {
		return fmt.Errorf("failed to execute GET request to %s: %w", path, err)
	}
	return nil
}

func (s *baseService) Post(ctx context.Context, path string, body any, resp any) error {
	err := s.client.Do(ctx, http.MethodPost, s.rootPath+path, body, resp)
	if err != nil {
		return fmt.Errorf("failed to execute POST request to %s: %w", path, err)
	}
	return nil
}

func (s *baseService) Put(ctx context.Context, path string, body any, resp any) error {
	err := s.client.Do(ctx, http.MethodPut, s.rootPath+path, body, resp)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request to %s: %w", path, err)
	}
	return nil
}

func (s *baseService) Patch(ctx context.Context, path string, body any, resp any) error {
	err := s.client.Do(ctx, http.MethodPatch, s.rootPath+path, body, resp)
	if err != nil {
		return fmt.Errorf("failed to execute PATCH request to %s: %w", path, err)
	}
	return nil
}

func (s *baseService) Delete(ctx context.Context, path string) error {
	err := s.client.Do(ctx, http.MethodDelete, s.rootPath+path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE request to %s: %w", path, err)
	}
	return nil
}

// EngineService presents engine service for specific engine (datastore type).
type EngineService struct {
	*baseService
	engine Engine
}

func (s *EngineService) DatastorePath(datastoreID string, parts ...string) string {
	path := []string{datastoreID}

	path = append(path, parts...)

	return "/" + strings.Join(path, "/")
}

func NewEngineService(client transport.Client, engine Engine) *EngineService {
	root := "/datastores/" + string(engine)

	return &EngineService{
		baseService: &baseService{
			client:   client,
			rootPath: root,
		},
		engine: engine,
	}
}
