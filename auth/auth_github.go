package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hashicorp/vault/api"
	vaultGithub "github.com/hashicorp/vault/builtin/credential/github"
)

// Github auths against a vault instance using the GitHub auth method.
func Github(gitToken, vaultHost string) (string, error) {

	vaultCFG := *api.DefaultConfig()
	vaultCFG.Address = vaultHost

	vClient, err := api.NewClient(&vaultCFG)
	if err != nil {
		return "", err
	}

	mountInput := map[string]string{
		"mount": "github",
		"token": gitToken,
	}

	// create vault client instance
	cliHandler := vaultGithub.CLIHandler{}
	var cTok string

	cTok, err = cliHandler.Auth(vClient, mountInput)
	if err != nil {
		return "", err
	}
	return cTok, nil
}

// CreateGithubToken creates a github token for a user.
func CreateGithubToken(username, pass, otp, tokenName string) (string, error) {

	auth := github.BasicAuthTransport{
		Username: username,
		Password: pass,
		OTP:      otp,
	}

	client := github.NewClient(auth.Client())

	// scope is "read:org"
	scopes := []github.Scope{
		github.ScopeReadOrg,
		github.ScopeRepo,
	}

	// Create the token name.
	scopeNote := fmt.Sprintf("DRUD Sanctuary Token-%s", tokenName)

	ar := &github.AuthorizationRequest{
		Scopes: scopes,
		Note:   &scopeNote,
	}
	lo := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	for i := 0; i < 2; i++ {
		az, _, err := client.Authorizations.Create(ar)
		if err != nil {
			if i == 1 {
				return "", err
			}
			// if the token exists then delete it and try again
			if strings.Contains(err.Error(), "already_exists") {
				var authID int
				// get list of exiting tokens so we can get the id of the one we wnat
				azs, _, listErr := client.Authorizations.List(lo)
				if listErr != nil {
					return "", err
				}

				for _, authorization := range azs {
					if authorization.Note != nil {
						if *authorization.Note == scopeNote {
							authID = *authorization.ID
						}
					}
				}

				_, err = client.Authorizations.Delete(authID)
				if err != nil {
					return "", err
				}
			}
			// hopefully we haved resolved the issue with creating the authorization so we can retry
			continue
		}
		// set global var with new token
		return *az.Token, nil
	}

	return "", errors.New("Could not generate token.")
}
