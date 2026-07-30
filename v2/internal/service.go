package internal

import (
	"strings"

	"github.com/selectel/dbaas-go/internal/transport"
)

type EngineService struct {
	Client transport.Client
	Engine string
}

func NewEngineService(c transport.Client, engine string) EngineService {

	return EngineService{
		Client: c,
		Engine: engine,
	}
}

func (s EngineService) Path(parts ...string) string {

	path := []string{
		"datastores",
		s.Engine,
	}

	path = append(path, parts...)

	return "/" + strings.Join(path, "/")
}
