package internal

import (
	"strings"

	"github.com/selectel/dbaas-go/internal/transport"
)

type EngineService struct {
	client transport.Client
	engine string
}

func NewEngineService(c transport.Client, engine string) EngineService {

	if engine == "" {
		panic("engine must not be empty")
	}

	return EngineService{
		client: c,
		engine: engine,
	}
}

func (s EngineService) Path(parts ...string) string {

	path := []string{
		"datastores",
		s.engine,
	}

	path = append(path, parts...)

	return "/" + strings.Join(path, "/")
}
