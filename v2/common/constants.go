package common

// DatastoreState represents custom type for DBaaS v2 datastore states.
type DatastoreState string

const (
	DatastoreStateRunning  DatastoreState = "RUNNING"
	DatastoreStateDegraded DatastoreState = "DEGRADED"
	DatastoreStateDown     DatastoreState = "DOWN"
	DatastoreStateUnknown  DatastoreState = "UNKNOWN"
)

// DatastoreStatus represents custom type for DBaaS v2 datastore statuses.
type DatastoreStatus string

const (
	DatastoreStatusCreating   DatastoreStatus = "CREATING"
	DatastoreStatusActive     DatastoreStatus = "ACTIVE"
	DatastoreStatusUpdatin    DatastoreStatus = "UPDATING"
	DatastoreStatusResizing   DatastoreStatus = "RESIZING"
	DatastoreStatusError      DatastoreStatus = "ERROR"
	DatastoreStatusDeleting   DatastoreStatus = "DELETING"
	DatastoreStatusDeleted    DatastoreStatus = "DELETED"
	DatastoreStatusLocking    DatastoreStatus = "LOCKING"
	DatastoreStatusLocked     DatastoreStatus = "LOCKED"
	DatastoreStatusUnlocking  DatastoreStatus = "UNLOCKING"
	DatastoreStatusRebuilding DatastoreStatus = "REBUILGING"
	DatastoreStatusMaintaince DatastoreStatus = "MAINTENANCE"
)

// FlavorDiskType represents custom type for flavor disk types.
type FlavorDiskType string

const (
	FlavorDiskLocal        FlavorDiskType = "local"
	FlavorDiskNetworkUltra FlavorDiskType = "network-ultra"
)

// Engine represents engine type (datastore type).
type Engine string

const (
	EngineClickHouse Engine = "clickhouse"
	EngineOpenSearch Engine = "opensearch"
)
