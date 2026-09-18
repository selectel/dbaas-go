# dbaas-go: Go SDK for Selectel DBaaS
[![Go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/selectel/dbaas-go/)
[![Go Report Card](https://goreportcard.com/badge/github.com/selectel/dbaas-go)](https://goreportcard.com/report/github.com/selectel/dbaas-go)

Package dbaas-go provides Go SDK to work with Selectel DBaaS

## Documentation

The Go library documentation is available at [go.dev](https://pkg.go.dev/github.com/selectel/dbaas-go).

## What this library is capable of

The SDK provides access to Selectel Managed Databases Service resources through DBaaS API v1 and v2 clients.

### DBaaS API v1

The v1 client provides access to the existing DBaaS API v1 resources, including:

* acls
* available extensions
* configuration parameters
* databases
* datastores
* datastore types
* extensions
* flavors
* grants
* logical replication slots
* prometheus metrics tokens
* topics
* users

### DBaaS API v2

The v2 client currently provides:

* datastore types
* flavors
* ClickHouse datastores
* ClickHouse node groups
* ClickHouse shard groups
* ClickHouse datastore configuration parameters

## Getting started

### Instalation

You can install `dbaas-go` package via `go get` command:

```bash
go get github.com/selectel/dbaas-go
```

### Authentication

To work with the Selectel Managed Databases Service API you first need to:

* Create a Selectel account: [registration page](https://my.selectel.ru/registration).
* Create a project in Selectel Cloud Platform [projects](https://my.selectel.ru/vpc/projects).
* Retrieve a token for your project via API or [go-selvpcclient](https://github.com/selectel/go-selvpcclient).

### Endpoints

Selectel Managed Databases Service API endpoints are available for different regions.

Specify the API version in the endpoint according to the client you use by appending `/v1` or `/v2`.

For example:

```text
https://ru-2.dbaas.selcloud.ru/v1  # DBaaS API v1
https://ru-2.dbaas.selcloud.ru/v2  # DBaaS API v2
```

The currently available DBaaS API endpoints are:

| URL                                   | Region |
| ------------------------------------- | ------ |
| https://ru-1.dbaas.selcloud.ru        | ru-1   |
| https://ru-2.dbaas.selcloud.ru        | ru-2   |
| https://ru-3.dbaas.selcloud.ru        | ru-3   |
| https://ru-6.dbaas.selcloud.ru        | ru-6   |
| https://ru-7.dbaas.selcloud.ru        | ru-7   |
| https://ru-8.dbaas.selcloud.ru        | ru-8   |
| https://ru-9.dbaas.selcloud.ru        | ru-9   |
| https://ke-1.dbaas.api.servercore.com | ke-1   |
| https://kz-1.dbaas.api.servercore.com | kz-1   |
| https://uz-1.dbaas.api.servercore.com | uz-1   |
| https://uz-2.dbaas.api.servercore.com | uz-2   |

You can also retrieve all available API endpoints from the Identity catalog.


### Docs
You can use Godoc to view methods and signatures
```shell
go install golang.org/x/tools/cmd/godoc@latest
export PATH="$PATH:$HOME/go/bin/"
godoc -http=localhost:6060
```
View
http://localhost:6060/pkg/github.com/selectel/dbaas-go/

### Usage example

#### DBaaS API v1
```go
package main

import (
    "context"
    "log"
    "fmt"

    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
    "github.com/selectel/dbaas-go"
)

func main() {
    // Token to work with Selectel Cloud project.
    token := "TOKEN"

    // DBaaS endpoint to work with.
    endpoint := "https://ru-2.dbaas.selcloud.ru/v1"

    openstackEndpoint := "https://api.selvpc.ru/identity/v3/"

    openstackRegion := "ru-2"

    // Initialize the DBaaS v1 client.
    dbaasClient, err := dbaas.NewDBAASClient(token, endpoint)
    if err != nil {
        log.Fatal(err)
    }

    // Prepare empty context.
    ctx := context.Background()

    // Get available datastore types.
    datastoreTypes, err := dbaasClient.DatastoreTypes(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Auth options for openstack to get all subnets.
    devopts := gophercloud.AuthOptions{
        IdentityEndpoint: openstackEndpoint,
        TokenID:          token,
    }

    provider, err := openstack.AuthenticatedClient(devopts)
    if err != nil {
        log.Fatal(err)
    }

    // Create a new network client.
    networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{Region: openstackRegion})
    if err != nil {
        log.Fatal(err)
    }

    // Get a list of available subnets.
    listOpts := subnets.ListOpts{
        IPVersion: 4,
    }
    allPages, err := subnets.List(networkClient, listOpts).AllPages()
    if err != nil {
        log.Fatal(err)
    }
    allSubnets, err := subnets.ExtractSubnets(allPages)
    if err != nil {
        log.Fatal(err)
    }

    // Create options for a new datastore.
    datastoreCreateOpts := dbaas.DatastoreCreateOpts{
        Name:      "go_cluster",
        TypeID:    datastoreTypes[0].ID,
        NodeCount: 1,
        SubnetID:  allSubnets[0].ID,
        Flavor:    &dbaas.Flavor{Vcpus: 2, RAM: 4096, Disk: 32},
    }

    // Create a new datastore.
    newDatastore, err := dbaasClient.CreateDatastore(ctx, datastoreCreateOpts)
    if err != nil {
        log.Fatal(err)
    }

    // Print datastores fields.
    fmt.Printf("Created datastore: %+v\n", newDatastore)
}

```

#### DBaaS API v2
```go
package main

import (
	"context"
	"fmt"
	"log"

	dbaas "github.com/selectel/dbaas-go/v2"
	"github.com/selectel/dbaas-go/v2/clickhouse"
	"github.com/selectel/dbaas-go/v2/common"
)

func main() {
	// Token to work with Selectel Cloud project.
	token := "TOKEN"

	// DBaaS v2 endpoint to work with.
	endpoint := "https://ru-2.dbaas.selcloud.ru/v2"

	// Initialize the DBaaS v2 client.
	client, err := dbaas.NewAPI(token, endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare empty context.
	ctx := context.Background()

	// Create a new ClickHouse datastore.
    shardOneWeight := 100
	shardOneHasPublicIPs := false
	datastore, err := client.ClickHouse.CreateDatastore(
		ctx,
		clickhouse.DatastoreCreateRequest{
			Name:     "go_cluster",
			TypeID:   "DATASTORE_TYPE_ID",
			SubnetID: "SUBNET_ID",
			Password: "PASSWORD",
			NodeGroups: []clickhouse.NodeGroupCreateRequest{
				{
					Name:      "shard1",
					Role:      clickhouse.NodeGroupRoleData,
					NodeCount: 1,
					Flavor: clickhouse.FlavorForNodeGroupRequest{
						Type: common.FlavorTypeFIXED,
						ID:   "FLAVOR_ID",
						DiskType: common.FlavorDiskLocal,
					},
					Weight:       &shardOneWeight,
					HasPublicIPs: &shardOneHasPublicIPs,
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print datastore fields.
	fmt.Printf("Created datastore: %+v\n", datastore)
}
```

The `TypeID`, `SubnetID` and `Flavor.ID` values must be replaced with IDs available in your Selectel project.
