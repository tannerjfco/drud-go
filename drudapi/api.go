package drudapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Entity interface represents eve entities in some functinos
type Entity interface {
	Path() string                // returns the path that must be added to host to get the entity
	Unmarshal(data []byte) error // unmarshal json into entity's fields
	JSON() []byte                //returns the entity's json representation
}

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

// Applications entity
type Applications struct {
	Name  string
	Items []Application `json:"_items"`
	Meta  struct {
		MaxResults int `json:"max_results"`
		Page       int `json:"page"`
		Total      int `json:"total"`
	} `json:"_meta"`
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

// Client ...
type Client struct {
	Created string `json:"_created,omitempty"`
	Etag    string `json:"_etag,omitempty"`
	ID      string `json:"_id,omitempty"`
	Updated string `json:"_updated,omitempty"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

// Path ...
func (c Client) Path() string {
	return "client/" + c.Name
}

// Unmarshal ...
func (c *Client) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &c)
	return err
}

// JSON ...
func (c Client) JSON() ([]byte, error) {
	jbytes, err := json.Marshal(c)
	return jbytes, err
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
}

// Credentials gets passed around to functions for authenticating with the api
type Credentials struct {
	Username   string `json:"username"`
	Password   string
	Token      string `json:"auth_token"`
	AdminToken string `json:"admin_token"`
}

// Request type used for building requests
type Request struct {
	Host string
	Auth *Credentials
}

/*
Get Method for handling GET, POST, PATCH requests to the api

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
			AdminToken: "fecb013f3c7138cd044568ec9bec074d8a0ab8e4",
		},
	}

	err := r.Get(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}

}

*/
func (r *Request) Get(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path())

	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("Error making GET request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if r.Auth != nil {
		// check for admin token, then auth token, then user Credentials
		if r.Auth.AdminToken != "" {
			req.Header.Set("Authorization", "token "+r.Auth.AdminToken)
		} else if r.Auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+r.Auth.Token)
		} else {
			req.SetBasicAuth(r.Auth.Username, r.Auth.Password)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Handle different status codes
	if resp.StatusCode-200 > 100 {
		return fmt.Errorf("%s: %d", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = entity.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) Post(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path())

	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(entity.JSON()))
	if err != nil {
		return errors.New("Error creating NewRequest: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	if r.Auth != nil {
		// check for admin token, then auth token, then user Credentials
		if r.Auth.AdminToken != "" {
			req.Header.Set("Authorization", "token "+r.Auth.AdminToken)
		} else if r.Auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+r.Auth.Token)
		} else {
			req.SetBasicAuth(r.Auth.Username, r.Auth.Password)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Handle different status codes
	if resp.StatusCode-200 > 100 {
		return fmt.Errorf("%s: %d", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = entity.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
