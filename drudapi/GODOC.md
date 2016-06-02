# drudapi
--
    import "github.com/drud/drud-go/drudapi"


## Usage

#### type Application

```go
type Application struct {
	AppID        string   `json:"app_id"`
	Client       string   `json:"client"`
	Deploys      []Deploy `json:"deploys"`
	GithubHookID int      `json:"github_hook_id"`
	RepoOrg      string   `json:"repo_org"`
	Name         string   `json:"name"`
	Repo         string   `json:"repo"`
	Created      string   `json:"_created,omitempty"`
	Etag         string   `json:"_etag,omitempty"`
	ID           string   `json:"_id,omitempty"`
	Updated      string   `json:"_updated,omitempty"`
}
```

Application ...

#### type Applications

```go
type Applications struct {
	Name  string
	Items []Application `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}
```

Applications entity

#### type Client

```go
type Client struct {
	Created string `json:"_created,omitempty"`
	Etag    string `json:"_etag,omitempty"`
	ID      string `json:"_id,omitempty"`
	Updated string `json:"_updated,omitempty"`
	Email   string `json:"email"`
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone"`
}
```

Client ...

#### func (Client) ETAG

```go
func (c Client) ETAG() string
```
ETAG ...

#### func (Client) JSON

```go
func (c Client) JSON() []byte
```
JSON ...

#### func (Client) PatchJSON

```go
func (c Client) PatchJSON() []byte
```
PatchJSON ...

#### func (Client) Path

```go
func (c Client) Path(method string) string
```
Path ...

#### func (*Client) Unmarshal

```go
func (c *Client) Unmarshal(data []byte) error
```
Unmarshal ...

#### type ClientList

```go
type ClientList struct {
	Items []Client `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}
```

ClientList ...

#### func (ClientList) ETAG

```go
func (c ClientList) ETAG() string
```
ETAG ...

#### func (ClientList) JSON

```go
func (c ClientList) JSON() []byte
```
JSON ...

#### func (ClientList) PatchJSON

```go
func (c ClientList) PatchJSON() []byte
```
PatchJSON ...

#### func (ClientList) Path

```go
func (c ClientList) Path(method string) string
```
Path ...

#### func (*ClientList) Unmarshal

```go
func (c *ClientList) Unmarshal(data []byte) error
```
Unmarshal ...

#### type Credentials

```go
type Credentials struct {
	Username   string `json:"username"`
	Password   string
	Token      string `json:"auth_token"`
	AdminToken string `json:"admin_token"`
}
```

Credentials gets passed around to functions for authenticating with the api

#### type Deploy

```go
type Deploy struct {
	Name          string `json:"name,omitempty"`
	Controller    string `json:"controller,omitempty"`
	Branch        string `json:"branch,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	BasicAuthUser string `json:"basicauth_user,omitempty"`
	BasicAuthPass string `json:"basicauth_pass,omitempty"`
}
```

Deploy ...

#### type Entity

```go
type Entity interface {
	Path(method string) string   // returns the path that must be added to host to get the entity
	Unmarshal(data []byte) error // unmarshal json into entity's fields
	JSON() []byte                //returns the entity's json representation
	PatchJSON() []byte           //returns the entity's json repr minus id field
	ETAG() string                // returns etag
}
```

Entity interface represents eve entities in some functinos

#### type Provider

```go
type Provider struct {
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}
```

Provider ...

#### type Region

```go
type Region struct {
	Name string `json:"name"`
}
```

Region ...

#### type Request

```go
type Request struct {
	Host  string // base path of the api  e.g. https://drudapi.genesis.drud.io/v0.1
	Query string // optional query params e.g. `where={"name":"fred"}``
	Auth  *Credentials
}
```

Request type used for building requests

#### func (*Request) Delete

```go
func (r *Request) Delete(entity Entity) error
```
Patch ...

#### func (*Request) Get

```go
func (r *Request) Get(entity Entity) error
```
Get ...

#### func (*Request) Patch

```go
func (r *Request) Patch(entity Entity) error
```
Patch ...

#### func (*Request) Post

```go
func (r *Request) Post(entity Entity) error
```
Post ...

#### type User

```go
type User struct {
	Username string      `json:"username"`
	Hashpw   string      `json:"hashpw"`
	Token    string      `json:"auth_token,omitempty"`
	Created  string      `json:"_created,omitempty"`
	Etag     string      `json:"_etag,omitempty"`
	ID       string      `json:"_id,omitempty"`
	Updated  string      `json:"_updated,omitempty"`
	Auth     Credentials `json:"-"`
}
```

User represents a user entity from the api

#### type Users

```go
type Users struct {
	Name  string
	Items []User `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}
```

Users entity
