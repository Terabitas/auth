// +build integration

package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/nildev/auth/Godeps/_workspace/src/golang.org/x/oauth2/github"

	"github.com/nildev/auth/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"
	"github.com/nildev/auth/Godeps/_workspace/src/github.com/jarcoal/httpmock"
	. "github.com/nildev/auth/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/nildev/auth/oauth"
)

type AuthIntegrationSuite struct{}

var _ = Suite(&AuthIntegrationSuite{})

func TestMain(m *testing.M) {
	os.Setenv("ND_SIGN_KEY", "signing-key")

	code := m.Run()
	os.Exit(code)
}

func (s *AccountIntegrationSuite) TestIfAccountWhichDoesNotExistsIsBeingRegisteredReal(c *C) {
	rez, err := Auth(oauth.PROVIDER_GITHUB, "a97b59a9df69b56bbba1")

	token, err := jwt.Parse(rez, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ND_SIGN_KEY")), nil
	})

	c.Assert(token.Valid, Equals, true)
	c.Assert(err, IsNil)
}

func (s *AccountIntegrationSuite) TestIfJWTIsReturned(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", github.Endpoint.TokenURL,
		httpmock.NewStringResponder(200, `{"access_token":"e72e16c7e42f292c6912e7710c838347ae178b4a", "scope":"repo,gist", "token_type":"bearer"}`))

	httpmock.RegisterResponder("GET", "https://api.github.com/user",
		httpmock.NewStringResponder(200, `{
  "login": "octocat",
  "id": 1,
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa octocat",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "octocat@github.com",
  "hireable": false,
  "bio": "There once was...",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "private_gists": 81,
  "disk_usage": 10000,
  "collaborators": 8,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}`))

	rez, err := Auth(oauth.PROVIDER_GITHUB, "xxxx")

	token, err := jwt.Parse(rez, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ND_SIGN_KEY")), nil
	})

	c.Assert(token.Valid, Equals, true)
	c.Assert(err, IsNil)
}
