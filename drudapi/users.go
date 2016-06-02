package drudapi

import "encoding/json"

// User represents a user entity from the api
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

// Path ...
func (u User) Path(method string) string {
	var path string

	if method == "POST" {
		path = "users"
	} else {
		path = "users/" + u.ID
	}
	return path
}

// Unmarshal ...
func (u *User) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, u)
	return err
}

// JSON ...
func (u User) JSON() []byte {
	u.ID = ""
	u.Etag = ""
	u.Created = ""
	u.Updated = ""

	jbytes, _ := json.Marshal(u)
	return jbytes
}

// PatchJSON ...
func (u User) PatchJSON() []byte {
	u.ID = ""
	u.Etag = ""
	u.Created = ""
	u.Updated = ""

	jbytes, _ := json.Marshal(u)
	return jbytes
}

// ETAG ...
func (u User) ETAG() string {
	return u.Etag
}

// UserList entity
type UserList struct {
	Items []User `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}

// Path ...
func (u UserList) Path(method string) string {
	return "users"
}

// Unmarshal ...
func (u *UserList) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &u)
	return err
}
