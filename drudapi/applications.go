package drudapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Deploy ...
type Deploy struct {
	Name          string `json:"name,omitempty"`
	Template      string `json:"template,omitempty"`
	Branch        string `json:"branch,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	BasicAuthUser string `json:"basicauth_user,omitempty"`
	BasicAuthPass string `json:"basicauth_pass,omitempty"`
	AutoManaged   bool   `json:"auto_managed,omitempty"`
	MigrateFrom   string `json:"migrate_from,omitempty"`
}

// Application ...
type Application struct {
	AppID        string   `json:"app_id"`
	Client       Client   `json:"client"`
	Deploys      []Deploy `json:"deploys,omitempty"`
	GithubHookID int      `json:"github_hook_id,omitempty"`
	RepoOrg      string   `json:"repo_org,omitempty"`
	Name         string   `json:"name"`
	Repo         string   `json:"repo,omitempty"`
	RepoDetails  *struct {
		Host   string `json:"host,omitempty"`
		Name   string `json:"name,omitempty"`
		Org    string `json:"org,omitempty"`
		Branch string `json:"branch,omitempty"`
		Dest   string `json:"dest,omitempty"`
	} `json:"repo_details,omitempty"`
	Created string `json:"_created,omitempty"`
	Etag    string `json:"_etag,omitempty"`
	ID      string `json:"_id,omitempty"`
	Updated string `json:"_updated,omitempty"`
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
	a.RepoDetails = nil

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

// RepoURL ...
func (a *Application) RepoURL(token string) string {
	var url string
	if a.RepoDetails != nil {
		if token == "" {
			url = fmt.Sprintf("https://%s@%s/%s/%s.git",
				token,
				a.RepoDetails.Host,
				a.RepoDetails.Org,
				a.Name,
			)
		} else {
			url = fmt.Sprintf("https://%s/%s/%s.git", a.RepoDetails.Host, a.RepoDetails.Org, a.Name)
		}
	}

	return url
}

// Describe an application..mostly used for displaying deploys
func (a *Application) Describe() {
	fmt.Println("App:", a.Name, "Client:", a.Client.Name)
	fmt.Printf("\n%v %v found.\n", len(a.Deploys), FormatPlural(len(a.Deploys), "deploy", "deploys"))
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	defer tabWriter.Flush()

	fmt.Fprintln(tabWriter, "\nNAME\tTEMPLATE\tBRANCH\tHOSTNAME\tBASICAUTH USERNAME\tBASICAUTH PASSWORD\tPROTOCOL\tAUTO MANAGED")
	for _, dep := range a.Deploys {
		var managed string

		if dep.AutoManaged == true {
			managed = "     ✓"
		}

		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			dep.Name,
			dep.Template,
			dep.Branch,
			dep.Hostname,
			dep.BasicAuthUser,
			dep.BasicAuthPass,
			dep.Protocol,
			managed,
		)
	}
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

// Describe pretty prints the entity
func (a *ApplicationList) Describe() {
	fmt.Printf("%v %v found.\n", len(a.Items), FormatPlural(len(a.Items), "application", "applications"))
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	defer tabWriter.Flush()

	fmt.Fprintln(tabWriter, "\nNAME\tCLIENT\tAPP(s)\tCREATED\tUPDATED")
	for _, app := range a.Items {
		// gather list of deploys by name
		var appNames []string
		for _, dep := range app.Deploys {
			appNames = append(appNames, dep.Name)
		}

		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\n",
			app.Name,
			app.Client.Name,
			strings.Join(appNames, ","),
			app.Created,
			app.Updated,
		)
	}
}
