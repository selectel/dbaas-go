package dbaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDatastoreNotFoundResponse = `{
	"error": {
		"code": 404,
		"title": "Not Found",
		"message": "database %s not found."
	}
}`

const testCreateDatastoreInvalidDatastoreIDResponse = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message":
			"Validation failure: {'datastore.type_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}"
	}
}`

const testUpdateDatastoreInvalidName = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": "Validation failure: {'datastore.name': \"'' is too short\"}"
	}
}`

const testResizeDatastoreInvalidNodeCount = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": "Validation failure: {'resize.node_count': '0 is less than the minimum of 1'}"
	}
}`

//nolint:lll //The line is 130 characters long, which exceeds the maximum of 120 characters.
const testResizeDatastoreWithDiskType = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message": "Validation failure: {'resize.flavor': \"Additional properties are not allowed ('disk_type' was unexpected)\"}"
	}
}`

const testPoolerDatastoreInvalidMode = `{
	"error": {
		"code": 400,
		"title": "Bad Request",
		"message":
			"Validation failure: {'pooler.mode': \"'abc' is not one of ['session', 'transaction', 'statement']\"}"
	}
}`

const testDatastoresResponse = `{
	"datastores": [
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "123e4567e89b12d3a456426655440000",
			"name": "Name",
			"status": "ACTIVE",
			"enabled": true,
			"type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"subnet_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"node_count": 1,
			"is_maintenance": false,
			"is_protected": false,
			"backup_retention_days" : 7,
			"connection": {
				"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
				"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org"
			},
			"flavor": {
				"vcpus": 2,
				"ram": 2048,
				"disk": 32,
				"disk_type": "local" 
			},
			"instances": [
				{
				"id": "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				"ip": "127.0.0.1",
				"floating_ip": "None",
				"role": "MASTER",
				"status": "ACTIVE",
				"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org"
				}
			],
			"pooler": {
				"size": 30,
				"mode": "session"
			},
			"firewall": [
				{
					"ip": "127.0.0.1"
				}
			],
			"config": {}
		},
		{
			"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
			"created_at": "1970-01-01T00:00:00",
			"updated_at": "1970-01-01T00:00:00",
			"project_id": "123e4567e89b12d3a456426655440000",
			"name": "AnotherName",
			"status": "ACTIVE",
			"enabled": true,
			"type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"subnet_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"node_count": 1,
			"is_maintenance": false,
			"is_protected": false,
			"backup_retention_days" : 7,
			"connection": {
				"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5.c.dbaas.selcloud.org",
				"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5.c.dbaas.selcloud.org"
			},
			"flavor": {
				"vcpus": 2,
				"ram": 2048,
				"disk": 32,
				"disk_type": "network-ultra"
			},
			"instances": [
				{
					"id": "40d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
					"ip": "127.0.0.1",
					"floating_ip": "None",
					"role": "MASTER",
					"status": "ACTIVE",
					"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org"
				},
				{
					"id": "50d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
					"ip": "127.0.0.2",
					"floating_ip": "None",
					"role": "REPLICA",
					"status": "ACTIVE",
					"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c56.ru-1.c.dbaas.selcloud.org"
				}
			],
			"pooler": {
				"size": 30,
				"mode": "session"
			},
			"firewall": [
				{
					"ip": "127.0.0.1"
				}
			],
			"config": {}
		}
	]
}
`

const testDatastoreResponse = `{
	"datastore": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "123e4567e89b12d3a456426655440000",
		"name": "Name",
		"status": "ACTIVE",
		"enabled": true,
		"type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"subnet_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"node_count": 1,
		"is_maintenance": false,
		"is_protected": false,
		"backup_retention_days" : 7,
		"connection": {
			"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org"
		},
		"flavor": {
			"vcpus": 2,
			"ram": 2048,
			"disk": 32,
			"disk_type": "local"
		},
		"instances": [
			{
			"id": "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"ip": "127.0.0.1",
			"floating_ip": "192.168.1.1",
			"role": "MASTER",
			"status": "ACTIVE",
			"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org"
			}
		],
		"pooler": {
			"size": 30,
			"mode": "session"
		},
		"firewall": [
			{
				"ip": "127.0.0.1"
			}
		],
		"config": {}
	}
}`

const testMultiNodeDatastoreResponse = `{
	"datastore": {
		"id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"created_at": "1970-01-01T00:00:00",
		"updated_at": "1970-01-01T00:00:00",
		"project_id": "123e4567e89b12d3a456426655440000",
		"name": "Name",
		"status": "ACTIVE",
		"enabled": true,
		"type_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"subnet_id": "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		"node_count": 3,
		"is_maintenance": false,
		"is_protected": false,
		"backup_retention_days" : 7,
		"connection": {
			"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"replica-1": "replica-1.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"replica-2": "replica-2.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org"
		},
		"flavor": {
			"vcpus": 2,
			"ram": 2048,
			"disk": 32
		},
		"instances": [
			{
			"id": "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"ip": "127.0.0.1",
			"floating_ip": "192.168.1.1",
			"role": "MASTER",
			"status": "ACTIVE",
			"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org"
			},
			{
			"id": "40d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"ip": "127.0.0.2",
			"floating_ip": "192.168.1.2",
			"role": "REPLICA",
			"status": "ACTIVE",
			"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c56.ru-1.c.dbaas.selcloud.org"
			},
			{
			"id": "50d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			"floating_ip": "None",
			"ip": "127.0.0.3",
			"role": "REPLICA",
			"status": "ACTIVE",
			"hostname": "9c387698-42a9-4555-9a8c-46eee7dc8c57.ru-1.c.dbaas.selcloud.org"
			}
		],
		"pooler": {
			"size": 30,
			"mode": "session"
		},
		"firewall": [
			{
				"ip": "127.0.0.1"
			}
		],
		"config": {}
	}
}`

const datastoreID = "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4"

var datastoreListExpected []Datastore = []Datastore{ //nolint
	{
		ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:           "1970-01-01T00:00:00",
		UpdatedAt:           "1970-01-01T00:00:00",
		ProjectID:           "123e4567e89b12d3a456426655440000",
		Name:                "Name",
		Status:              "ACTIVE",
		Enabled:             true,
		TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount:           1,
		IsMaintenance:       false,
		IsProtected:         false,
		BackupRetentionDays: 7,
		Connection: map[string]string{
			"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		},
		Flavor: Flavor{
			Vcpus:    2,
			RAM:      2048,
			Disk:     32,
			DiskType: "local",
		},
		Instances: []Instances{{
			ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			IP:         "127.0.0.1",
			FloatingIP: "None",
			Role:       "MASTER",
			Status:     "ACTIVE",
			Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
		}},
		Pooler: Pooler{
			Size: 30,
			Mode: "session",
		},
		Firewall: []Firewall{{
			IP: "127.0.0.1",
		}},
		Config: map[string]interface{}{},
	},
	{
		ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5",
		CreatedAt:           "1970-01-01T00:00:00",
		UpdatedAt:           "1970-01-01T00:00:00",
		ProjectID:           "123e4567e89b12d3a456426655440000",
		Name:                "AnotherName",
		Status:              "ACTIVE",
		Enabled:             true,
		TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount:           1,
		IsMaintenance:       false,
		IsProtected:         false,
		BackupRetentionDays: 7,
		Connection: map[string]string{
			"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5.c.dbaas.selcloud.org",
			"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f5.c.dbaas.selcloud.org",
		},
		Flavor: Flavor{
			Vcpus:    2,
			RAM:      2048,
			Disk:     32,
			DiskType: "network-ultra",
		},
		Instances: []Instances{
			{
				ID:         "40d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				IP:         "127.0.0.1",
				FloatingIP: "None",
				Role:       "MASTER",
				Status:     "ACTIVE",
				Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
			},
			{
				ID:         "50d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				IP:         "127.0.0.2",
				FloatingIP: "None",
				Role:       "REPLICA",
				Status:     "ACTIVE",
				Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c56.ru-1.c.dbaas.selcloud.org",
			},
		},
		Pooler: Pooler{
			Size: 30,
			Mode: "session",
		},
		Firewall: []Firewall{{
			IP: "127.0.0.1",
		}},
		Config: map[string]interface{}{},
	},
}

var datastoreCreateResponse Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingCreate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreCreateExpected Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingCreate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreUpdateResponse Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingUpdate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreUpdateExpected Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingUpdate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreResizeResponse Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusResizing,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreResizeExpected Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusResizing,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config:   map[string]interface{}{},
}

var datastoreUpdateConfigResponse Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingUpdate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config: map[string]interface{}{
		"work_mem":                 256.0,
		"session_replication_role": "replica",
	},
}

var datastoreUpdateConfigExpected Datastore = Datastore{ //nolint
	ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	CreatedAt:           "1970-01-01T00:00:00",
	UpdatedAt:           "1970-01-01T00:00:00",
	ProjectID:           "123e4567e89b12d3a456426655440000",
	Name:                "Name",
	Status:              StatusPendingUpdate,
	Enabled:             true,
	TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
	NodeCount:           1,
	IsMaintenance:       false,
	IsProtected:         false,
	BackupRetentionDays: 7,
	Connection: map[string]string{
		"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
	},
	Flavor: Flavor{
		Vcpus: 2,
		RAM:   2048,
		Disk:  32,
	},
	Instances: []Instances{{
		ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		IP:         "127.0.0.1",
		FloatingIP: "None",
		Role:       "MASTER",
		Status:     "ACTIVE",
		Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
	}},
	Pooler: Pooler{
		Size: 30,
		Mode: "session",
	},
	Firewall: []Firewall{},
	Config: map[string]interface{}{
		"work_mem":                 256.0,
		"session_replication_role": "replica",
	},
}

func TestDatastores(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoresURI,
		httpmock.NewStringResponder(200, testDatastoresResponse))

	actual, err := testClient.Datastores(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, datastoreListExpected, actual)
}

func TestDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoresURI+"/"+datastoreID,
		httpmock.NewStringResponder(200, testDatastoreResponse))

	expected := Datastore{
		ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:           "1970-01-01T00:00:00",
		UpdatedAt:           "1970-01-01T00:00:00",
		ProjectID:           "123e4567e89b12d3a456426655440000",
		Name:                "Name",
		Status:              "ACTIVE",
		Enabled:             true,
		TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount:           1,
		IsMaintenance:       false,
		IsProtected:         false,
		BackupRetentionDays: 7,
		Connection: map[string]string{
			"MASTER": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"master": "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		},
		Flavor: Flavor{
			Vcpus:    2,
			RAM:      2048,
			Disk:     32,
			DiskType: "local",
		},
		Instances: []Instances{{
			ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
			IP:         "127.0.0.1",
			FloatingIP: "192.168.1.1",
			Role:       "MASTER",
			Status:     "ACTIVE",
			Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
		}},
		Pooler: Pooler{
			Size: 30,
			Mode: "session",
		},
		Firewall: []Firewall{{
			IP: "127.0.0.1",
		}},
		Config: map[string]interface{}{},
	}

	actual, err := testClient.Datastore(context.Background(), datastoreID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMultiNodeDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoresURI+"/"+datastoreID,
		httpmock.NewStringResponder(200, testMultiNodeDatastoreResponse))

	expected := Datastore{
		ID:                  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		CreatedAt:           "1970-01-01T00:00:00",
		UpdatedAt:           "1970-01-01T00:00:00",
		ProjectID:           "123e4567e89b12d3a456426655440000",
		Name:                "Name",
		Status:              "ACTIVE",
		Enabled:             true,
		TypeID:              "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		SubnetID:            "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount:           3,
		IsMaintenance:       false,
		IsProtected:         false,
		BackupRetentionDays: 7,
		Connection: map[string]string{
			"MASTER":    "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"master":    "master.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"replica-1": "replica-1.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
			"replica-2": "replica-2.20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4.c.dbaas.selcloud.org",
		},
		Flavor: Flavor{
			Vcpus: 2,
			RAM:   2048,
			Disk:  32,
		},
		Instances: []Instances{
			{
				ID:         "30d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				IP:         "127.0.0.1",
				FloatingIP: "192.168.1.1",
				Role:       "MASTER",
				Status:     "ACTIVE",
				Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c55.ru-1.c.dbaas.selcloud.org",
			},
			{
				ID:         "40d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				IP:         "127.0.0.2",
				FloatingIP: "192.168.1.2",
				Role:       "REPLICA",
				Status:     "ACTIVE",
				Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c56.ru-1.c.dbaas.selcloud.org",
			},
			{
				ID:         "50d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
				IP:         "127.0.0.3",
				FloatingIP: "None",
				Role:       "REPLICA",
				Status:     "ACTIVE",
				Hostname:   "9c387698-42a9-4555-9a8c-46eee7dc8c57.ru-1.c.dbaas.selcloud.org",
			},
		},
		Pooler: Pooler{
			Size: 30,
			Mode: "session",
		},
		Firewall: []Firewall{{
			IP: "127.0.0.1",
		}},
		Config: map[string]interface{}{},
	}

	actual, err := testClient.Datastore(context.Background(), datastoreID)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestDatastoreNotFound(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	notFoundResponse := fmt.Sprintf(testDatastoreNotFoundResponse, NotFoundEntityID)
	httpmock.RegisterResponder("GET", testClient.Endpoint+DatastoresURI+"/"+NotFoundEntityID,
		httpmock.NewStringResponder(404, notFoundResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 404
	expected.APIError.Title = ErrorNotFoundTitle
	expected.APIError.Message = fmt.Sprintf("datastore %s not found.", NotFoundEntityID)

	_, err := testClient.Datastore(context.Background(), NotFoundEntityID)

	require.ErrorAs(t, err, &expected)
}

func TestCreateDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreCreateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	createDatastoreOpts := DatastoreCreateOpts{
		Name:      "Name",
		TypeID:    "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount: 1,
		SubnetID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Flavor:    &Flavor{Vcpus: 2, RAM: 2048, Disk: 32},
	}

	actual, err := testClient.CreateDatastore(context.Background(), createDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreCreateExpected, actual)
}

func TestCreateDatastoreWithFlavorDiskType(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastoreCreateResponseWithFlavorDiskType := datastoreCreateResponse
			datastoreCreateResponseWithFlavorDiskType.Flavor.DiskType = DiskNetworkUltra
			datastores["datastore"] = datastoreCreateResponseWithFlavorDiskType

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	createDatastoreOpts := DatastoreCreateOpts{
		Name:      "Name",
		TypeID:    "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount: 1,
		SubnetID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Flavor:    &Flavor{Vcpus: 2, RAM: 2048, Disk: 32, DiskType: DiskNetworkUltra},
	}

	datastoreCreateExpectedWithDiskType := datastoreCreateExpected
	datastoreCreateExpectedWithDiskType.Flavor.DiskType = DiskNetworkUltra

	actual, err := testClient.CreateDatastore(context.Background(), createDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreCreateExpectedWithDiskType, actual)
}

func TestCreateDatatastoreInvalidTypeID(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI,
		httpmock.NewStringResponder(400, testCreateDatastoreInvalidDatastoreIDResponse))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure:
		{'datastore.type_id': \"'20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f' is not a 'UUID'\"}`

	createDatastoreOpts := DatastoreCreateOpts{
		Name:      "Name",
		TypeID:    "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f",
		NodeCount: 1,
		SubnetID:  "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Flavor:    &Flavor{Vcpus: 2, RAM: 2048, Disk: 32},
	}

	_, err := testClient.CreateDatastore(context.Background(), createDatastoreOpts)

	require.ErrorAs(t, err, &expected)
}

func TestUpdateDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreUpdateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	updateDatastoreOpts := DatastoreUpdateOpts{
		Name: "NewName",
	}

	actual, err := testClient.UpdateDatastore(context.Background(), datastoreID, updateDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreUpdateExpected, actual)
}

func TestUpdateDatatastoreInvalidName(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID,
		httpmock.NewStringResponder(400, testUpdateDatastoreInvalidName))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = "Validation failure: {'datastore.name': \"'' is too short\"}"

	updateDatastoreOpts := DatastoreUpdateOpts{
		Name: "",
	}

	_, err := testClient.UpdateDatastore(context.Background(), datastoreID, updateDatastoreOpts)

	require.ErrorAs(t, err, &expected)
}

func TestResizeDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	//nolint:goconst // string `/resize` has 3 occurrences, make it a constant
	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/resize",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreResizeOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreResizeResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	resizeDatastoreOpts := DatastoreResizeOpts{
		NodeCount: 3,
		Flavor:    &Flavor{Vcpus: 2, RAM: 4096, Disk: 32},
	}

	actual, err := testClient.ResizeDatastore(context.Background(), datastoreID, resizeDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreResizeExpected, actual)
}

func TestResizeDatatastoreInvalidNodeCount(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/resize",
		httpmock.NewStringResponder(400, testResizeDatastoreInvalidNodeCount))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = "Validation failure: {'resize.node_count': '0 is less than the minimum of 1'}"

	resizeDatastoreOpts := DatastoreResizeOpts{
		NodeCount: 0,
	}

	_, err := testClient.ResizeDatastore(context.Background(), datastoreID, resizeDatastoreOpts)

	require.ErrorAs(t, err, &expected)
}

func TestResizeDatatastoreWithDiskType(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/resize",
		httpmock.NewStringResponder(400, testResizeDatastoreWithDiskType))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure: {'resize.flavor': \"Additional properties are not allowed 
	('disk_type' was unexpected)\"}`

	resizeDatastoreOpts := DatastoreResizeOpts{
		Flavor: &Flavor{Vcpus: 2, RAM: 4096, Disk: 32, DiskType: DiskNetworkUltra},
	}

	_, err := testClient.ResizeDatastore(context.Background(), datastoreID, resizeDatastoreOpts)

	require.ErrorAs(t, err, &expected)
}

func TestPoolerDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/pooler",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastorePoolerOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	poolerDatastoreOpts := DatastorePoolerOpts{
		Mode: "transaction",
		Size: 50,
	}

	actual, err := testClient.PoolerDatastore(context.Background(), datastoreID, poolerDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreUpdateExpected, actual)
}

func TestPoolerDatatastoreInvalidMode(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/pooler",
		httpmock.NewStringResponder(400, testPoolerDatastoreInvalidMode))

	expected := &DBaaSAPIError{}
	expected.APIError.Code = 400
	expected.APIError.Title = ErrorBadRequestTitle
	expected.APIError.Message = `Validation failure:
		{'pooler.mode': \"'abc' is not one of ['session', 'transaction', 'statement']\"}`

	poolerDatastoreOpts := DatastorePoolerOpts{
		Mode: "abc",
		Size: 30,
	}

	_, err := testClient.PoolerDatastore(context.Background(), datastoreID, poolerDatastoreOpts)

	require.ErrorAs(t, err, &expected)
}

func TestFirewallDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/firewall",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreFirewallOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	firewallDatastoreOpts := DatastoreFirewallOpts{
		IPs: []string{"127.0.0.1", "127.0.0.2"},
	}

	actual, err := testClient.FirewallDatastore(context.Background(), datastoreID, firewallDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreUpdateExpected, actual)
}

func TestConfigDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/config",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreConfigOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreUpdateConfigResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	configDatastoreOpts := DatastoreConfigOpts{
		Config: map[string]interface{}{
			"work_mem":                 256,
			"session_replication_role": "replica",
		},
	}

	actual, err := testClient.ConfigDatastore(context.Background(), datastoreID, configDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreUpdateConfigExpected, actual)
}

func TestPasswordDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/password",
		func(req *http.Request) (*http.Response, error) {
			tempPasswordOpts := struct {
				Password DatastorePasswordOpts `json:"password"`
			}{
				Password: DatastorePasswordOpts{},
			}
			if err := json.NewDecoder(req.Body).Decode(&tempPasswordOpts); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreUpdateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	passwordDatastoreOpts := DatastorePasswordOpts{
		RedisPassword: "ohf1xi$geiSh7ae8eixush5roo:sho4n",
	}

	actual, err := testClient.PasswordDatastore(context.Background(), datastoreID, passwordDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreUpdateExpected, actual)
}

func TestBackupsDatastore(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()
	expectedResponse := datastoreUpdateExpected
	expectedResponse.BackupRetentionDays = 14

	httpmock.RegisterResponder("PUT", testClient.Endpoint+DatastoresURI+"/"+datastoreID+"/backups",
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreBackupsOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = expectedResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	backupsDatastoreOpts := DatastoreBackupsOpts{
		BackupRetentionDays: 14,
	}

	actual, err := testClient.BackupsDatastore(context.Background(), datastoreID, backupsDatastoreOpts)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, actual)
}

func TestCreateDatastoreWithFloatingIPs(t *testing.T) {
	httpmock.Activate()
	testClient := SetupTestClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", testClient.Endpoint+DatastoresURI,
		func(req *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(req.Body).Decode(&DatastoreCreateOpts{}); err != nil {
				return httpmock.NewStringResponse(400, ""), err
			}

			datastores := make(map[string]Datastore)
			datastores["datastore"] = datastoreCreateResponse

			resp, err := httpmock.NewJsonResponse(200, datastores)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), err
			}
			return resp, nil
		})

	createDatastoreOpts := DatastoreCreateOpts{
		Name:        "Name",
		TypeID:      "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		NodeCount:   1,
		SubnetID:    "20d7bcf4-f8d6-4bf6-b8f6-46cb440a87f4",
		Flavor:      &Flavor{Vcpus: 2, RAM: 2048, Disk: 32},
		FloatingIPs: &FloatingIPs{Master: 1, Replica: 0},
	}

	actual, err := testClient.CreateDatastore(context.Background(), createDatastoreOpts)

	require.NoError(t, err)
	assert.Equal(t, datastoreCreateExpected, actual)
}
