# drudapi


Client CRUD example

```go
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "gittoken",
		},
	}

	c := &drudapi.Client{
		Name: "turtle",
	}

	fmt.Println("POsting")
	err := r.Post(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", c)
	}

	c.Phone = "123-123-1235"
	c.Email = "my@email.com"

	fmt.Println("Patching")
	err = r.Patch(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", c)
	}

	fmt.Println("Deleting")
	err = r.Delete(c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Getting")
	err = r.Get(c)
	if err != nil {
		fmt.Println(err)
	}

}

```
