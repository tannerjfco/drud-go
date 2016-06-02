package drudapi

import (
	"log"
	"testing"
)

func TestGetApplication200(t *testing.T) {
	expectedResp := `
{
    "_updated": "Mon, 23 May 2016 20:23:37 GMT",
    "github_hook_id": 123455,
    "name": "killah",
    "app_id": "mumbojumbo",
    "repo_org": "ourorg",
    "repo": "",
    "client": {
        "_updated": "Mon, 23 May 2016 20:23:34 GMT",
        "name": "somebiz",
        "email": "",
        "phone": "",
        "_created": "Mon, 23 May 2016 20:23:34 GMT",
        "_id": "574366c6e2638a001f430115",
        "_etag": "9906d3a8584f0fabbd96451013b38d20fff5f5d3"
    },
    "deploys": [
        {
            "protocol": "http",
            "basicauth_pass": "drud",
            "basicauth_user": "drud",
            "name": "default",
            "branch": "master"
        }
    ],
    "_created": "Mon, 23 May 2016 20:23:37 GMT",
    "_id": "d546fh5ert456",
    "_etag": "34b4d4c312a2d5cb916fef5330b2a14f53acac4b"
}`
	server := getTestServer(200, expectedResp)
	defer server.Close()

	r := Request{
		Host: server.URL,
		Auth: &Credentials{
			AdminToken: "dfgdfg",
		},
	}

	a := &Application{
		AppID: "mumbojumbo",
	}

	err := r.Get(a)
	if err != nil {
		log.Fatal(err)
	}

	expect(t, a.GithubHookID, 123455)
	expect(t, a.Name, "killah")
	expect(t, a.RepoOrg, "ourorg")
	expect(t, a.Created, "Mon, 23 May 2016 20:23:37 GMT")
	expect(t, a.Etag, "34b4d4c312a2d5cb916fef5330b2a14f53acac4b")
	expect(t, a.Client.Name, "somebiz")

}
