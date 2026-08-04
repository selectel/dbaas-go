package dbaas

import (
	"errors"

	"github.com/selectel/dbaas-go/internal/transport"
)

// Error titles.
const (
	ErrorNotFoundTitle   = "Not Found"
	ErrorBadRequestTitle = "Bad Request"
)

var ErrorEndpointVersionMismatch = errors.New("endpoint version mismatch")

// DBaaSAPIError alias is a type of an error raised by API calls made by this library.
type DBaaSAPIError = transport.DBaaSAPIError
