package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Instances represents datastore's instances.
type Instances struct {
	ID               string `json:"id"`
	IP               string `json:"ip"`
	FloatingIP       string `json:"floating_ip"`
	Role             string `json:"role"`
	RoleName         string `json:"role_name"`
	Status           Status `json:"status"`
	Hostname         string `json:"hostname"`
	AvailabilityZone string `json:"availability_zone"`
}

// Flavor represents datastore's flavor.
type Flavor struct {
	DiskType DiskType `json:"disk_type,omitempty"`
	Vcpus    int      `json:"vcpus"`
	RAM      int      `json:"ram"`
	Disk     int      `json:"disk"`
}

// Restore represents restore parameters for datastore.
type Restore struct {
	DatastoreID string `json:"datastore_id,omitempty"`
	TargetTime  string `json:"target_time,omitempty"`
}

// Pooler represents pooler parameters for datastore.
type Pooler struct {
	Mode string `json:"mode,omitempty"`
	Size int    `json:"size,omitempty"`
}

// Firewall represents firewall rules parameters for datastore.
type Firewall struct {
	IP string `json:"ip"`
}

// FloatingIPs represents floating IPs creation schema.
type FloatingIPs struct {
	Master  int `json:"master"`
	Replica int `json:"replica"`
}

// Datastore is the API response for the datastores.
type Datastore struct {
	ID                  string            `json:"id"`
	CreatedAt           string            `json:"created_at"`
	UpdatedAt           string            `json:"updated_at"`
	CreationFinishedAt  string            `json:"creation_finished_at"`
	ProjectID           string            `json:"project_id"`
	Name                string            `json:"name"`
	TypeID              string            `json:"type_id"`
	SubnetID            string            `json:"subnet_id"`
	FlavorID            string            `json:"flavor_id"`
	Status              Status            `json:"status"`
	Connection          map[string]string `json:"connection"`
	Firewall            []Firewall        `json:"firewall"`
	Instances           []Instances       `json:"instances"`
	Config              map[string]any    `json:"config"`
	Pooler              Pooler            `json:"pooler"`
	Flavor              Flavor            `json:"flavor"`
	NodeCount           int               `json:"node_count"`
	Enabled             bool              `json:"enabled"`
	AllowRestore        bool              `json:"allow_restore"`
	IsMaintenance       bool              `json:"is_maintenance"`
	IsProtected         bool              `json:"is_protected"`
	BackupRetentionDays int               `json:"backup_retention_days"`
	DatabasesCount      int               `json:"databases_count"`
	TopicsCount         int               `json:"topics_count"`
	DiskUsed            int               `json:"disk_used"`
}

// Disk represents disk parameters for a get/create datastore ops.
type Disk struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

// ResizeDisk represents disk parameters for a datastore resize operation.
type ResizeDisk struct {
	Size int `json:"size"`
}

// DatastoreCreateOpts represents options for the datastore Create request.
type DatastoreCreateOpts struct {
	Flavor              *Flavor        `json:"flavor,omitempty"`
	Restore             *Restore       `json:"restore,omitempty"`
	Pooler              *Pooler        `json:"pooler,omitempty"`
	FloatingIPs         *FloatingIPs   `json:"floating_ips,omitempty"`
	Config              map[string]any `json:"config,omitempty"`
	Disk                *Disk          `json:"disk,omitempty"`
	TypeID              string         `json:"type_id"`
	SubnetID            string         `json:"subnet_id"`
	FlavorID            string         `json:"flavor_id,omitempty"`
	ProjectID           string         `json:"project_id"`
	RedisPassword       string         `json:"redis_password,omitempty"`
	Name                string         `json:"name"`
	NodeCount           int            `json:"node_count"`
	BackupRetentionDays int            `json:"backup_retention_days,omitempty"`
}

// DatastoreUpdateOpts represents options for the datastore Update request.
type DatastoreUpdateOpts struct {
	Name string `json:"name"`
}

// DatastoreResizeOpts represents options for the datastore Resize request.
type DatastoreResizeOpts struct {
	Flavor    *Flavor     `json:"flavor,omitempty"`
	Disk      *ResizeDisk `json:"disk,omitempty"`
	FlavorID  string      `json:"flavor_id,omitempty"`
	NodeCount int         `json:"node_count,omitempty"`
}

// DatastorePoolerOpts represents options for the datastore's pooler Update request.
type DatastorePoolerOpts struct {
	Mode string `json:"mode,omitempty"`
	Size int    `json:"size,omitempty"`
}

// DatastoreFirewallOpts represents options for the datastore's firewall rules Update request.
type DatastoreFirewallOpts struct {
	IPs []string `json:"ips"`
}

// DatastoreConfigOpts represents options for the datastore's configuration parameters Update request.
type DatastoreConfigOpts struct {
	Config map[string]any `json:"config"`
}

// DatastorePasswordOpts represents options for the Redis datastore's password Update request.
type DatastorePasswordOpts struct {
	RedisPassword string `json:"redis_password"`
}

// DatastoreQueryParams represents available query parameters for datastore.
type DatastoreQueryParams struct {
	ID            string `json:"id,omitempty"`
	ProjectID     string `json:"project_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Status        Status `json:"status,omitempty"`
	Enabled       string `json:"enabled,omitempty"`
	TypeID        string `json:"type_id,omitempty"`
	FlavorID      string `json:"flavor_id,omitempty"`
	SubnetID      string `json:"subnet_id,omitempty"`
	AllowRestore  bool   `json:"allow_restore,omitempty"`
	IsMaintenance bool   `json:"is_maintenance,omitempty"`
	IsProtected   bool   `json:"is_protected,omitempty"`
	Deleted       bool   `json:"deleted,omitempty"`
}

// DatastoreBackupsOpts represents update options for the Datastore backups.
type DatastoreBackupsOpts struct {
	BackupRetentionDays int `json:"backup_retention_days"`
}

const DatastoresURI = "/datastores"

// Datastores returns all datastores.
func (api *API) Datastores(ctx context.Context, params *DatastoreQueryParams) ([]Datastore, error) {
	uri, err := setQueryParams(DatastoresURI, params)
	if err != nil {
		return []Datastore{}, err
	}

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Datastore{}, err
	}

	var result struct {
		Datastores []Datastore `json:"datastores"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return []Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastores, nil
}

// Datastore returns a datastore based on the ID.
func (api *API) Datastore(ctx context.Context, datastoreID string) (Datastore, error) {
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s", DatastoresURI, datastoreID)

	resp, err := api.makeRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// CreateDatastore creates a new datastore.
func (api *API) CreateDatastore(ctx context.Context, opts DatastoreCreateOpts) (Datastore, error) {
	config := convertConfigValues(opts.Config)
	createDatastoreOpts := struct {
		Datastore DatastoreCreateOpts `json:"datastore"`
	}{
		Datastore: opts,
	}
	createDatastoreOpts.Datastore.Config = config
	requestBody, err := json.Marshal(createDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, DatastoresURI, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// UpdateDatastore updates an existing datastore.
func (api *API) UpdateDatastore(ctx context.Context, datastoreID string, opts DatastoreUpdateOpts) (Datastore, error) {
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s", DatastoresURI, datastoreID)
	updateDatastoreOpts := struct {
		Datastore DatastoreUpdateOpts `json:"datastore"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(updateDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// DeleteDatastore deletes an existing datastore.
func (api *API) DeleteDatastore(ctx context.Context, datastoreID string) error {
	if err := uuid.Validate(datastoreID); err != nil {
		return fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s", DatastoresURI, datastoreID)

	_, err := api.makeRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}

// ResizeDatastore resizes an existing datastore.
func (api *API) ResizeDatastore(ctx context.Context, datastoreID string, opts DatastoreResizeOpts) (Datastore, error) {
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/resize", DatastoresURI, datastoreID)
	resizeDatastoreOpts := struct {
		Datastore DatastoreResizeOpts `json:"resize"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(resizeDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPost, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// PoolerDatastore updates pooler parameters of an existing datastore.
func (api *API) PoolerDatastore(ctx context.Context, datastoreID string, opts DatastorePoolerOpts) (Datastore, error) {
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/pooler", DatastoresURI, datastoreID)
	poolerDatastoreOpts := struct {
		Datastore DatastorePoolerOpts `json:"pooler"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(poolerDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// FirewallDatastore updates firewall rules of an existing datastore.
func (api *API) FirewallDatastore(ctx context.Context, datastoreID string, opts DatastoreFirewallOpts) (Datastore, error) { //nolint
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/firewall", DatastoresURI, datastoreID)
	firewallDatastoreOpts := struct {
		Datastore DatastoreFirewallOpts `json:"firewall"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(firewallDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// ConfigDatastore updates configuration parameters rules of an existing datastore.
func (api *API) ConfigDatastore(ctx context.Context, datastoreID string, opts DatastoreConfigOpts) (Datastore, error) { //nolint
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/config", DatastoresURI, datastoreID)
	opts.Config = convertConfigValues(opts.Config)
	requestBody, err := json.Marshal(opts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// PasswordDatastore updates password of an existing Redis datastore.
func (api *API) PasswordDatastore(ctx context.Context, datastoreID string, opts DatastorePasswordOpts) (Datastore, error) { //nolint
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/password", DatastoresURI, datastoreID)
	passwordDatastoreOpts := struct {
		Datastore DatastorePasswordOpts `json:"password"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(passwordDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}

// BackupsDatastore updates backups parameters of an existing datastore.
func (api *API) BackupsDatastore(ctx context.Context, datastoreID string, opts DatastoreBackupsOpts) (Datastore, error) { //nolint
	if err := uuid.Validate(datastoreID); err != nil {
		return Datastore{}, fmt.Errorf("Error during datastoreID validate, %w", err)
	}

	uri := fmt.Sprintf("%s/%s/backups", DatastoresURI, datastoreID)
	backupsDatastoreOpts := struct {
		Datastore DatastoreBackupsOpts `json:"backups"`
	}{
		Datastore: opts,
	}
	requestBody, err := json.Marshal(backupsDatastoreOpts)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error marshalling params to JSON, %w", err)
	}

	resp, err := api.makeRequest(ctx, http.MethodPut, uri, requestBody)
	if err != nil {
		return Datastore{}, err
	}

	var result struct {
		Datastore Datastore `json:"datastore"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Datastore{}, fmt.Errorf("Error during Unmarshal, %w", err)
	}

	return result.Datastore, nil
}
