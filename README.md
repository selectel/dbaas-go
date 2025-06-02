# dbaas-go: Go SDK for Selectel DBaaS
[![Go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/selectel/dbaas-go/)
[![Go Report Card](https://goreportcard.com/badge/github.com/selectel/dbaas-go)](https://goreportcard.com/report/github.com/selectel/dbaas-go)

Package dbaas-go provides Go SDK to work with Selectel DBaaS

## Documentation

The Go library documentation is available at [go.dev](https://pkg.go.dev/github.com/selectel/dbaas-go).

## What this library is capable of

You can use this library to work with the following objects of the Selectel Managed Databases Service:

* acl
* available extension
* configuration parameter
* database
* datastore
* datastore type
* extension
* flavor
* grant
* logical replication slots
* prometheus metrics tokens
* topic
* user

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

Selectel Managed Databases Service currently has the following API endpoint:

| URL                               | Region |
|-----------------------------------|--------|
| https://ru-1.dbaas.selcloud.ru/v1 | ru-1   |
| https://ru-2.dbaas.selcloud.ru/v1 | ru-2   |
| https://ru-3.dbaas.selcloud.ru/v1 | ru-3   |
| https://ru-7.dbaas.selcloud.ru/v1 | ru-7   |
| https://ru-8.dbaas.selcloud.ru/v1 | ru-8   |
| https://ru-9.dbaas.selcloud.ru/v1 | ru-9   |
| https://uz-1.dbaas.selcloud.ru/v1 | uz-1   |
| https://kz-1.dbaas.selcloud.ru/v1 | kz-1   |

You can also retrieve all available API endpoints from the Identity
catalog.

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
        ProjectID: <PORJECT_ID>
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
