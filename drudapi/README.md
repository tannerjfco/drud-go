# drudapi


Get Method

```go
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {
	fmt.Println("tes test")

	c := &drudapi.Client{
		Name: "1fee",
	}

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "githubtoken",
		},
	}

	err := r.Get(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}

}
```

Post

```go
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {

	c := &drudapi.Client{
		Name: "turtle",
	}

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "githubtoken",
		},
	}

	err := r.Post(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}

}

```

Patch

```go
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {

	c := &drudapi.Client{
		Name: "turtle",
	}

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "githubtoken",
		},
	}

	err := r.Get(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
		fmt.Println("------------------")
	}

	c.Phone = "123-123-1235"
	c.Email = "my@email.com"

	err = r.Patch(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}
}
```

Delete
```go
package main

import (
	"fmt"

	"github.com/drud/drud-go/drudapi"
)

func main() {

	c := &drudapi.Client{
		Name: "turtle",
	}

	r := drudapi.Request{
		Host: "https://drudapi.genesis.drud.io/v0.1",
		Auth: &drudapi.Credentials{
			AdminToken: "githubtoken",
		},
	}

	err := r.Get(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}

	err = r.Delete(c)
	if err != nil {
		fmt.Println(err)
	}

}
```
