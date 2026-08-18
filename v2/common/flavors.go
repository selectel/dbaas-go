package common

import (
	"context"

	"github.com/selectel/dbaas-go/internal/transport"
)

// FlavorHost stores nested options for FlavorResponse.
type FlavorHost struct {
	Line           string `json:"line"`
	Processor      string `json:"processor"`
	AvailableCount int    `json:"available_count"`
}

// FlavorResponse is the API response for the flavors.
type FlavorResponse struct {
	Host             *FlavorHost `json:"host,omitempty"`
	ID               string      `json:"id"`
	FlSize           string      `json:"fl_size"`
	DiskType         string      `json:"disk_type,omitempty"`
	Type             string      `json:"type,omitempty"`
	DatastoreTypeIDs []string    `json:"datastore_type_ids"`
	AllowedRoles     []string    `json:"allowed_roles,omitempty"`
	VCPUs            int         `json:"vcpus"`
	RAM              int         `json:"ram"`
	Disk             int         `json:"disk"`
	// Name             string      `json:"name,omitempty"`
	// Description      *string     `json:"description,omitempty"`
	// Subtype          string      `json:"subtype"`
}

// DatastoreTypeListResponse is the API response for datastore type list.
type FlavorListResponse struct {
	Errors  string           `json:"errors"`
	Flavors []FlavorResponse `json:"flavors"`
}

type FlavorService struct {
	*baseService
}

// GetDatastoreTypeList returns datastore type list from api.
func (s *FlavorService) GetFlavorList(ctx context.Context) (FlavorListResponse, error) {
	response := FlavorListResponse{}

	err := s.Get(ctx, "", &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func NewFlavorService(client transport.Client) *FlavorService {
	root := "/flavors"

	return &FlavorService{
		baseService: &baseService{
			client:   client,
			rootPath: root,
		},
	}
}
