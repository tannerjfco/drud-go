drudapi
-------

A Go API for interacting with Drud API.

## Get a Client

```
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {
	fmt.Println("tes test")

	client := &drudapi.Client{
		Name: "1fee",
	}

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "github_token",
		},
	}

	err := r.Get(client)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(client)
	}

}

```
