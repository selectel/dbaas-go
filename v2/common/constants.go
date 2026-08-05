package common

// Status represents custom type for various DBaaS v2 objects states.
type State string

const (
	StateRunning  State = "RUNNING"
	StateDegraded State = "DEGRADED"
	StateDown     State = "DOWN"
	StateUnknown  State = "UNKNOWN"
)

// Status represents custom type for various DBaaS v2 objects statuses.
type Status string

const (
	StatusCreating   Status = "CREATING"
	StatusActive     Status = "ACTIVE"
	StatusUpdatin    Status = "UPDATING"
	StatusResizing   Status = "RESIZING"
	StatusError      Status = "ERROR"
	StatusDeleting   Status = "DELETING"
	StatusDeleted    Status = "DELETED"
	StatusLocking    Status = "LOCKING"
	StatusLocked     Status = "LOCKED"
	StatusUnlocking  Status = "UNLOCKING"
	StatusRebuilding Status = "REBUILGING"
	StatusMaintaince Status = "MAINTENANCE"
)

// DiskType represents custom type for various flavor disk types.
type DiskType string

const (
	DiskLocal        DiskType = "local"
	DiskNetworkUltra DiskType = "network-ultra"
)

// Engine represents engine type (datastore type).
type Engine string

const (
	EngineClickHouse Engine = "clickhouse"
	EngineOpenSearch Engine = "opensearch"
)
