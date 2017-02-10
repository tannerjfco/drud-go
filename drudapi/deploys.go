package drudapi

import (
	"encoding/json"
	"fmt"

	"github.com/gosuri/uitable"
)

// Deploy ...
type Deploy struct {
	DeployID       string      `json:"deploy_id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Application    Application `json:"application,omitempty"`
	Template       string      `json:"template,omitempty"`
	Branch         string      `json:"branch,omitempty"`
	Hostname       string      `json:"hostname,omitempty"`
	Protocol       string      `json:"protocol,omitempty"`
	BasicAuthUser  string      `json:"basicauth_user,omitempty"`
	BasicAuthPass  string      `json:"basicauth_pass,omitempty"`
	AutoManaged    bool        `json:"auto_managed,omitempty"`
	MigrateFrom    string      `json:"migrate_from,omitempty"`
	Url            string      `json:"url,omitempty"`
	AuthKey        string      `json:"auth_key,omitempty"`
	SecureAuthKey  string      `json:"secure_auth_key,omitempty"`
	LoggedInKey    string      `json:"logged_in_key,omitempty"`
	NonceKey       string      `json:"nonce_key,omitempty"`
	AuthSalt       string      `json:"auth_salt,omitempty"`
	SecureAuthSalt string      `json:"secure_auth_salt,omitempty"`
	LoggedInSalt   string      `json:"logged_in_salt,omitempty"`
	NonceSalt      string      `json:"nonce_salt,omitempty"`
	Created        string      `json:"_created,omitempty"`
	Etag           string      `json:"_etag,omitempty"`
	ID             string      `json:"_id,omitempty"`
	Updated        string      `json:"_updated,omitempty"`
}

// Path ...
func (d Deploy) Path(method string) string {
	var path string

	if method == "POST" {
		path = "deploys"
	} else {
		path = "deploys/" + d.DeployID
	}
	return path
}

// Unmarshal ...
func (d *Deploy) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, d)
	return err
}

// JSON ...
func (d Deploy) JSON() []byte {
	d.ID = ""
	d.Etag = ""
	d.Created = ""
	d.Updated = ""

	jbytes, _ := json.Marshal(d)
	return jbytes
}

// PatchJSON ...
func (d Deploy) PatchJSON() []byte {
	d.ID = ""
	d.Etag = ""
	d.Created = ""
	d.Updated = ""
	// removing name because it has been setup as the id param in drudapi and cannot be  patched
	d.DeployID = ""

	jbytes, _ := json.Marshal(d)
	return jbytes
}

// ETAG ...
func (d Deploy) ETAG() string {
	return d.Etag
}

// Describe an application..mostly used for displaying deploys
func (d *Deploy) Describe() {

	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true // wrap columns

	table.AddRow("DEPLOY:", d.DeployID)
	table.AddRow("APP:", d.Application.AppID)
	table.AddRow("CLIENT:", d.Application.Client.Name)
	table.AddRow("CREATED:", d.Created)

	fmt.Println(table)

}

// DeployList entity
type DeployList struct {
	Name  string
	Items []Deploy `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
}

// Path ...
func (d DeployList) Path(method string) string {
	return "deploys"
}

// Unmarshal ...
func (d *DeployList) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &d)
	return err
}

// Describe pretty prints the entity
func (d *DeployList) Describe() {
	fmt.Printf("%v %v found.\n\n", len(d.Items), FormatPlural(len(d.Items), "deploy", "deploys"))

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("NAME", "CLIENT", "CREATED")
	for _, deploy := range d.Items {
		table.AddRow("DEPLOY ID:", deploy.DeployID)
		table.AddRow("NAME:", deploy.Name)
		table.AddRow("URL:", deploy.Url)
		table.AddRow("TEMPLATE:", deploy.Template)
		table.AddRow("BRANCH:", deploy.Branch)
		table.AddRow("AUTH USER:", deploy.BasicAuthUser)
		table.AddRow("AUTH PASS:", deploy.BasicAuthPass)
		table.AddRow("AUTO MANAGED:", deploy.AutoManaged)
		table.AddRow("\n")

	}
	fmt.Println(table)

}
