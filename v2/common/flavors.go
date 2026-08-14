package common

// BaseFlavor storers required fields for flavor model in api.
type BaseFlavor struct {
	ID       string         `json:"id"`
	DiskType FlavorDiskType `json:"disk_type"`
	FlSize   string         `json:"fl_size"`
	Disk     int            `json:"disk"`
	RAM      int            `json:"ram"`
	VCPUs    int            `json:"vcpus"`
}

// Aggregating endpoint /flavors has defferent models in response.
//
// Clickhouse
// {
// 	"id": "4882f722-5d20-406a-a119-ca5cedaa4426",
// 	"disk_type": "LOCAL",
// 	"disk": 35,
// 	"fl_size": "STANDARD",
// 	"ram": 4096,
// 	"type": "FIXED",
// 	"vcpus": 2,
// 	"subtype": "STANDARD",
// 	"allowed_roles": [
// 		"DATA",
// 		"KEEPER"
// 	],
// 	"datastore_type_ids": [
// 		"00000000-0000-0000-0000-000000000001",
// 		"00000000-0000-0000-0000-000000000002"
// 	]
// }
//
// Opensearch
// {
// 	"id": "ac3284cd-8730-4a51-9537-383d49556ced",
// 	"disk_type": "LOCAL",
// 	"disk": 32,
// 	"fl_size": "STANDARD",
// 	"ram": 8192,
// 	"type": "FIXED",
// 	"vcpus": 2,
// 	"subtype": "MEMORY",
// 	"allowed_roles": [
// 		"MANAGER",
// 		"DATA"
// 	],
// 	"datastore_type_ids": [
// 		"770379e9-3008-43bf-9e62-7fbb6bd4fa08"
// 	]
// },
//
// Evort
// {
// 	"id": "b72e03b8-8aad-4157-8dd0-bf4dfe5da814",
// 	"name": "2-8192-64-redis-high-freq",
// 	"description": null,
// 	"vcpus": 2,
// 	"ram": 8192,
// 	"disk": 64,
// 	"fl_size": "high_freq",
// 	"subtype": null,
// 	"datastore_type_ids": [
// 		"2c8d1fff-20b2-4fe0-8055-df9ee127929f"
// 	],
// 	"host": null
// },
