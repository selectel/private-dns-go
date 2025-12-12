## Getting started

### Installation

You can install needed `private-dns-go` packages via `go get` command:

```bash
go get github.com/selectel/private-dns-go
```

### Authentication

To work with the Selectel Cloud Private DNS API you first need to:

* Create a Selectel account: [registration page](https://my.selectel.ru/registration).
* Create a project in Selectel Cloud Platform [projects](https://my.selectel.ru/vpc/projects).
* Retrieve a token for your project via API or [go-selvpcclient](https://github.com/selectel/go-selvpcclient).

### Endpoints

You can find available endpoints [here](https://docs.selectel.ru/en/api/urls/).

### Usage example

```go
package main

import (
	"context"
	"fmt"
	"log"

	privatedns "github.com/selectel/private-dns-go/pkg/v1"
)

func main() {

	// Create the client.
    cfg := &privatedns.Config{
        // Token to work with Selectel Cloud project.
        AuthToken: "..."
        // Cloud private dns endpoint to work with.
        URL: "https://ru-3.cloud.api.selcloud.ru/private-dns/"
    }
	client := privatedns.NewPrivateDNSClient(cfg)

	// Get zones for project.
	zones, err := client.ListZones(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print the zones.
	for idx, zone := range zones {
		fmt.Printf("Zone %d: %+v", idx, zone)
	}
}

```
