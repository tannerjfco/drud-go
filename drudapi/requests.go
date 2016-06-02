package drudapi

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Entity interface represents eve entities in some functinos
type Entity interface {
	Path(method string) string   // returns the path that must be added to host to get the entity
	Unmarshal(data []byte) error // unmarshal json into entity's fields
	JSON() []byte                //returns the entity's json representation
	PatchJSON() []byte           //returns the entity's json repr minus id field
	ETAG() string                // returns etag
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

// Get ...
func (r *Request) Get(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path("GET"))

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

// Post ...
func (r *Request) Post(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path("POST"))

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

// Patch ...
func (r *Request) Patch(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path("PATCH"))

	req, err = http.NewRequest("PATCH", u.String(), bytes.NewBuffer(entity.PatchJSON()))
	if err != nil {
		return errors.New("Error creating NewRequest: " + err.Error())
	}

	req.Header.Set("If-Match", entity.ETAG())
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

// Patch ...
func (r *Request) Delete(entity Entity) error {
	var req *http.Request
	var err error

	u, err := url.Parse(r.Host)
	u.Path = path.Join(u.Path, entity.Path("DELETE"))

	req, err = http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return errors.New("Error creating NewRequest: " + err.Error())
	}

	req.Header.Set("If-Match", entity.ETAG())
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

	return nil
}
