package secrets

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

// ConfigVault sets globals for vault access
func ConfigVault(tokenFile string, vaultHost string) string {
	// ensure there is a token for use with vault unless this is the auth command
	SetSecretEditor(&editor, "atom -w")

	var err error
	if _, err = os.Stat(tokenFile); os.IsNotExist(err) {
		log.Fatal("No sanctuary token found. Run `drud auth --help`")
	}

	var cTok string
	cTok, err = getSanctuaryToken(tokenFile)
	if err != nil {
		log.Fatalln("Error reading token file", err)
	}

	vPtr, err := NewAuthVault(vaultHost, cTok)
	if err != nil {
		log.Fatalln("Could not create vault api client")
	}
	vault = *vPtr

	return cTok
}

// NewAuthVault returns an authenticated vault
func NewAuthVault(vaultHost string, token string) (*api.Logical, error) {
	vaultCFG := *api.DefaultConfig()
	vaultCFG.Address = vaultHost

	vClient, err := api.NewClient(&vaultCFG)
	if err != nil {
		return nil, err
	}

	vClient.SetToken(token)

	return vClient.Logical(), nil
}

// GetTokenDetails returns a map of the user's token info
func GetTokenDetails() (map[string]interface{}, error) {
	sobj := Secret{
		Path: "/auth/token/lookup-self",
	}

	err := sobj.Read()
	if err != nil {
		return nil, err
	}

	return sobj.Data, nil
}

func getSanctuaryToken(tokenFile string) (string, error) {
	fileBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(fileBytes)), nil
}

func GetSecretEditor() string {
	// allow user to have different editor for secrets
	// fall back to default editor
	editor := os.Getenv("SECRET_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	return editor
}

func SetSecretEditor(varToSet *string, defaultEditor string) {
	*varToSet = GetSecretEditor()
	if *varToSet == "" {
		*varToSet = defaultEditor
	}
}
