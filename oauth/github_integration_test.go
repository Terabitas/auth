package oauth

import (
	"os"
	"testing"

	. "github.com/nildev/auth/Godeps/_workspace/src/gopkg.in/check.v1"
)

type AccountIntegrationSuite struct{}

var _ = Suite(&AccountIntegrationSuite{})

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func (s *AccountIntegrationSuite) TestIfTokenIsReturned(c *C) {
	c.Skip("have to")
	// get code https://github.com/login/oauth/authorize?client_id=e976860773d6ab411aa3
	code := "d00bb6c4164c2226ea32"
	reader := makeGitHubProviderUsingEnv(nil)
	token, err := reader.Token(code)

	c.Assert(err, IsNil)
	c.Assert(token, NotNil)
}
