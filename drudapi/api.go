package drudapi

// -------- Types for representing Drud API Resources --------

// Deploy ...
type Deploy struct {
	Name          string `json:"name,omitempty"`
	Controller    string `json:"controller,omitempty"`
	Branch        string `json:"branch,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	BasicAuthUser string `json:"basicauth_user,omitempty"`
	BasicAuthPass string `json:"basicauth_pass,omitempty"`
}

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

// Users entity
type Users struct {
	Name  string
	Items []User `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}

// Region ...
type Region struct {
	Name string `json:"name"`
}

// Provider ...
type Provider struct {
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}
