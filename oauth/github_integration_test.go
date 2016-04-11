package oauth

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

type AccountIntegrationSuite struct{}

var _ = Suite(&AccountIntegrationSuite{})

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func (s *AccountIntegrationSuite) TestIfTokenIsReturnedUsingGithub(c *C) {
	c.Skip("have to")
	// get code https://github.com/login/oauth/authorize?client_id=e976860773d6ab411aa3
	code := "d00bb6c4164c2226ea32"
	reader := makeGitHubProviderUsingEnv(nil)
	token, err := reader.Token(code)

	c.Assert(err, IsNil)
	c.Assert(token, NotNil)
}

func (s *AccountIntegrationSuite) TestIfTokenIsReturnedUsingDigitalOcean(c *C) {
	c.Skip("have to")
	// get code https://cloud.digitalocean.com/v1/oauth/authorize?client_id=a4f61c8db9c24444a6376810d28b037e9c17cc64d0b8d6cc5f113f331457774d&redirect_uri=http://127.0.0.1/auth/do&response_type=code
	code := "7a6eac4470a39b4bd7cef5601f719ae45f5fbe73fa307420d5df93ae202a8dc9"
	reader := makeDigitalOceanProviderUsingEnv(nil)
	token, err := reader.Token(code)

	c.Assert(err, IsNil)
	c.Assert(token, NotNil)
}
