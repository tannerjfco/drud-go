package drudapi

import "encoding/json"

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

// Application ...
type Application struct {
	AppID        string   `json:"app_id"`
	Client       Client   `json:"client"`
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

// Path ...
func (a Application) Path(method string) string {
	var path string

	if method == "POST" {
		path = "application"
	} else {
		path = "application/" + a.AppID
	}
	return path
}

// Unmarshal ...
func (a *Application) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, a)
	return err
}

// JSON ...
func (a Application) JSON() []byte {
	a.ID = ""
	a.Etag = ""
	a.Created = ""
	a.Updated = ""

	jbytes, _ := json.Marshal(a)
	return jbytes
}

// PatchJSON ...
func (a Application) PatchJSON() []byte {
	a.ID = ""
	a.Etag = ""
	a.Created = ""
	a.Updated = ""
	// removing name because it has been setup as the id param in drudapi and cannot be  patched
	a.AppID = ""

	jbytes, _ := json.Marshal(a)
	return jbytes
}

// ETAG ...
func (a Application) ETAG() string {
	return a.Etag
}

// ApplicationList entity
type ApplicationList struct {
	Name  string
	Items []Application `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}

// Path ...
func (a ApplicationList) Path(method string) string {
	return "application"
}

// Unmarshal ...
func (a *ApplicationList) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &a)
	return err
}
