package drudapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/howeyc/gopass"
)

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

// Application ...
type Application struct {
	AppID        string      `json:"app_id"`
	Client       string      `json:"client"`
	Deploys      []Deploy    `json:"deploys"`
	GithubHookID int         `json:"github_hook_id"`
	RepoOrg      string      `json:"repo_org"`
	Name         string      `json:"name"`
	Repo         string      `json:"repo"`
	Created      string      `json:"_created,omitempty"`
	Etag         string      `json:"_etag,omitempty"`
	ID           string      `json:"_id,omitempty"`
	Updated      string      `json:"_updated,omitempty"`
	Auth         Credentials `json:"-"`
}

// Get an app
func (app *Application) Get(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}

	// set aside auth because it will be written over
	auth := app.Auth

	q := u.Query()
	u.Path = path.Join(u.Path, "application")
	if app.AppID != "" {
		u.Path = path.Join(u.Path, app.AppID)
	}
	if app.Name != "" && app.Client != "" {
		q.Set("where", fmt.Sprintf(`{"name":"%s", "client":"%s"}`, app.Name, app.Client))
	}
	q.Set("embedded", `{"client":0,"deploys.controller":0}`)

	u.RawQuery = q.Encode()

	r := &request{
		Type: "GET",
		Path: u.String(),
		Auth: app.Auth,
	}
	bodyText, err := r.send()
	//@todo handle errors here better
	if err != nil {
		return err
	}

	if strings.Contains(string(bodyText), `"_items":`) {
		var apps Applications
		if err := json.Unmarshal(bodyText, &apps); err != nil {
			return err
		}
		if len(apps.Items) == 0 {
			return errors.New("No App found")
		}
		*app = apps.Items[0]
	} else {
		if err := json.Unmarshal(bodyText, &app); err != nil {
			return err
		}
	}
	// reset the auth
	app.Auth = auth

	return nil
}

// Update an application.
func (app *Application) Update(evehost string, etag string) error {
	payload, err := json.Marshal(app)
	if err != nil {
		return err
	}

	r := &request{
		Host:    evehost,
		Type:    "PATCH",
		Path:    path.Join("application", app.AppID),
		Auth:    app.Auth,
		Etag:    etag,
		Payload: payload,
	}

	bodyText, err := r.send()
	if err != nil {
		return fmt.Errorf("Patch failed %s", bodyText)
	}
	return err
}

// Delete an app
func (app *Application) Delete(evehost string) error {

	// auth := Credentials{
	// 	Token: cfg.Token,
	// }

	r := &request{
		Host: evehost,
		Etag: app.Etag,
		Type: "DELETE",
		Path: path.Join("application", app.AppID),
		Auth: app.Auth,
	}

	_, err := r.send()
	if err != nil {
		return err
	}

	return nil
}

// Applications entity
type Applications struct {
	Name  string
	Items []Application `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
	Auth Credentials `json:"-"`
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
	Auth Credentials `json:"-"`
}

// Get ...
func (usr *User) Get(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}
	q := u.Query()
	u.Path = path.Join(u.Path, "users")
	if usr.ID != "" {
		u.Path += usr.ID
	} else if usr.Username != "" {
		q.Set("where", fmt.Sprintf(`{"username":"%s"}`, usr.Username))
	}
	u.RawQuery = q.Encode()

	r := &request{
		Type: "GET",
		Path: u.String(),
		Auth: usr.Auth,
	}
	bodyText, err := r.send()
	//@todo handle errors here better
	if err != nil {
		return err
	}

	if strings.Contains(string(bodyText), `"_items":`) {
		var us Users
		if err := json.Unmarshal(bodyText, &us); err != nil {
			return err
		}
		if len(us.Items) == 0 {
			return errors.New("No Client found")
		}
		*usr = us.Items[0]
	} else {
		if err := json.Unmarshal(bodyText, &usr); err != nil {
			return err
		}
	}

	return nil
}

// Delete entity
func (usr *User) Delete(host string) error {

	r := &request{
		Host: host,
		Etag: usr.Etag,
		Type: "DELETE",
		Path: path.Join("users", usr.ID),
		Auth: usr.Auth,
	}

	_, err := r.send()
	if err != nil {
		return err
	}

	return nil
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

// Client ...
type Client struct {
	Created string      `json:"_created,omitempty"`
	Etag    string      `json:"_etag,omitempty"`
	ID      string      `json:"_id,omitempty"`
	Updated string      `json:"_updated,omitempty"`
	Email   string      `json:"email"`
	Name    string      `json:"name"`
	Phone   string      `json:"phone"`
	Auth    Credentials `json:"-"`
	EveHost string      `json:"-"`
}

// Get ...
func (c *Client) Get(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}
	q := u.Query()
	u.Path = path.Join(u.Path, "client", c.Name)
	u.RawQuery = q.Encode()

	r := &request{
		Type: "GET",
		Path: u.String(),
		Auth: c.Auth,
	}
	bodyText, err := r.send()
	//@todo handle errors here better
	if err != nil {
		return err
	}

	if strings.Contains(string(bodyText), `"_items":`) {
		var cs Clients
		if err := json.Unmarshal(bodyText, &cs); err != nil {
			return err
		}
		if len(cs.Items) == 0 {
			return errors.New("No Client found")
		}
		*c = cs.Items[0]
	} else {
		if err := json.Unmarshal(bodyText, &c); err != nil {
			return err
		}
	}

	return nil
}

// Delete entity
func (c *Client) Delete() error {

	r := &request{
		Host: c.EveHost,
		Etag: c.Etag,
		Type: "DELETE",
		Path: path.Join("client", c.Name),
		Auth: c.Auth,
	}

	_, err := r.send()
	if err != nil {
		return err
	}

	return nil
}

// Clients ...
type Clients struct {
	Name  string
	Items []Client `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
	Auth Credentials `json:"-"`
}

// Get ...
func (cs *Clients) Get(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}

	q := u.Query()
	u.Path = path.Join(u.Path, "client")

	if cs.Name != "" {
		q.Set("where", fmt.Sprintf(`{"name":"%s"}`, cs.Name))
	}

	u.RawQuery = q.Encode()
	r := &request{
		Type: "GET",
		Path: u.String(),
		Auth: cs.Auth,
	}
	bodyText, err := r.send()
	//@todo handle errors here better
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyText, cs); err != nil {
		return err
	}

	return nil
}

// Describe prtins the Client list in a human readable way
func (cs *Clients) Describe() {
	fmt.Printf("%v %v found.\n", len(cs.Items), FormatPlural(len(cs.Items), "client", "clients"))
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	defer tabWriter.Flush()
	if len(cs.Items) > 0 {
		fmt.Fprintln(tabWriter, "\nNAME\tEMAIL\tPHONE")
		for _, v := range cs.Items {
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n",
				v.Name,
				v.Email,
				v.Phone,
			)
		}
	}
}

// Credentials gets passed around to functions for authenticating with the api
type Credentials struct {
	Username   string `json:"username"`
	Password   string
	Token      string `json:"auth_token"`
	AdminToken string `json:"admin_token"`
}

// Type used for building requests
type request struct {
	Type    string
	Host    string
	Path    string
	Auth    Credentials
	Payload []byte
	Etag    string
}

// method for constructing url from base parts
func (r *request) url() string {
	// @todo handle host that has no trailing slash
	return r.Host + r.Path
}

// Method for handling GET, POST, PATCH requests to the api
func (r *request) send() ([]byte, error) {
	var req *http.Request
	var err error
	if string(r.Payload) == "" {
		req, err = http.NewRequest(r.Type, r.url(), nil)
	} else {
		req, err = http.NewRequest(r.Type, r.url(), bytes.NewBuffer(r.Payload))
	}
	if err != nil {
		return []byte(""), errors.New("Error creating NewRequest: " + err.Error())
	}

	if r.Type == "PATCH" || r.Type == "DELETE" {
		req.Header.Set("If-Match", r.Etag)
	}

	// check for admin token, then auth token, then user Credentials
	if r.Auth.AdminToken != "" {
		req.Header.Set("Authorization", "token "+r.Auth.AdminToken)
	} else if r.Auth.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Auth.Token)
	} else {
		req.SetBasicAuth(r.Auth.Username, r.Auth.Password)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), errors.New("Failed to connect to server")
	}
	// Handle different status codes
	if resp.StatusCode == 406 {
		return []byte(""), errors.New("Token expired: run `drud Login` and retry your previous task")
		//401 is unauthorized
	} else if resp.StatusCode == 401 {
		if r.Auth.Token == "" {
			return []byte(""), errors.New("Bad user/pass")
		}
		return []byte(""), errors.New("Bad token")
	} else if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		return []byte(""), errors.New(string(resp.Status) + ": " + string(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

// GetResource sends a GET request to the Drud api /endpoint and returns
// the entity or the entity in ["_items"][0]
// endpoint is an eve querystring like `users?where={"username":"bilbo"}`
// return one item and fails if there are more or none.
func GetResource(url string, auth *Credentials) ([]interface{}, error) {
	//send an Authenticated GET request to the eve host + endpoint
	r := &request{
		Type: "GET",
		Path: url,
		Auth: *auth,
	}
	bodyText, err := r.send()
	//@todo handle errors here better
	if err != nil {
		return nil, err
	}
	//unmarshal the returned json string into a map of string => interface
	var dat map[string]interface{}
	if err := json.Unmarshal(bodyText, &dat); err != nil {
		panic(err)
	}

	var items []interface{}
	if _, exists := dat["_items"]; exists == true {
		//type assertion to array (this allows to use indexes)
		items = dat["_items"].([]interface{})
		//return error if there query returns more than one or no items
		if len(items) == 0 {
			return nil, errors.New("returned 0 results")
		}
	} else {
		items = append(items, dat)
	}

	//return map of resource (note: values will need to be type asserted)
	return items, nil
}

// UpdateResource sends a patch with the json in payload to the endpoint
func UpdateResource(evehost string, endpoint string, query string, auth *Credentials, payload []byte) (map[string]interface{}, error) {
	var querypath string
	if strings.Contains(query, "where") {
		querypath = endpoint + query
	} else {
		querypath = endpoint + "/" + query
	}
	//fmt.Println(evehost + querypath)
	resources, err := GetResource(evehost+querypath, auth)
	if err != nil {
		return nil, err
	}
	resource := resources[0].(map[string]interface{})
	objectID := resource["_id"].(string)
	etag := resource["_etag"].(string)
	path := endpoint + "/" + objectID
	r := &request{
		Host:    evehost,
		Type:    "PATCH",
		Path:    path,
		Auth:    *auth,
		Etag:    etag,
		Payload: payload,
	}
	body, err := r.send()
	if err != nil {
		return nil, errors.New("Patch failed")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	return dat, nil
}

// CreateResource ...
func CreateResource(host string, endpoint string, payload []byte, auth *Credentials) (map[string]interface{}, error) {
	r := &request{
		Host:    host,
		Type:    "POST",
		Path:    endpoint,
		Auth:    *auth,
		Payload: payload,
	}
	resp, err := r.send()
	if err != nil {
		return nil, err
	}

	var dat map[string]interface{}
	if err = json.Unmarshal(resp, &dat); err != nil {
		panic(err)
	}

	ID := dat["_id"].(string)
	if endpoint == "application" {
		ID = dat["app_id"].(string)
	} else if endpoint == "client" {
		ID = dat["name"].(string)
	}

	// Now get the user so we can print the token
	// @todo get post to users to return token so i dont have to do this
	resources, err := GetResource(fmt.Sprintf("%s%s/%s", host, endpoint, ID), auth)
	if err != nil {
		return nil, err
	}
	resource := resources[0].(map[string]interface{})

	return resource, nil
}

// IsTokenValid makes a GET request to EVE Host (base URL) using the passed in token
// if the request returns somethign other than 200 then false is returned
// @todo as we implement better role-based authentication this may need to test
// at a different url due to permissins changes.
func IsTokenValid(token string, eveHost string) bool {
	auth := Credentials{
		Token: token,
	}
	r := &request{
		Type: "GET",
		Host: eveHost,
		Auth: auth,
	}
	_, err := r.send()

	if err != nil {
		return false
	}

	return true
}

// RenewUserToken sends a Patch to the user with an empty auth_token and returns the new token
func RenewUserToken(auth *Credentials, eveHost string) (string, error) {
	endpoint := `users`
	query := `?where={"username":"` + auth.Username + `"}`
	payload := []byte(`{"auth_token":""}`)

	entity, err := UpdateResource(eveHost, endpoint, query, auth, payload)
	if err != nil {
		return "", err
	}

	return entity["auth_token"].(string), nil
}

// FormatPlural is a simple wrapper which returns different strings based on the count value.
func FormatPlural(count int, single string, plural string) string {
	if count == 1 {
		return single
	}
	return plural
}

// GetMaskedInput gets a string from prompt to user but masks the input from the terminal
func GetMaskedInput() string {
	pass, _ := gopass.GetPasswdMasked()
	return string(pass)
}

// GetCredentials prompts user for credentials
func GetCredentials() (string, string) {
	var username string
	fmt.Println("Username:")
	//read value from terminal into username var
	fmt.Scanf("%s", &username)
	fmt.Printf("Password:\n")
	//get user's password without echoing to terminal
	pass := GetMaskedInput()

	return username, string(pass)
}
