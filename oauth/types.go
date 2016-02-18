package oauth

import (
	"github.com/nildev/auth/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/nildev/auth/Godeps/_workspace/src/github.com/juju/errors"
	"github.com/nildev/auth/Godeps/_workspace/src/golang.org/x/oauth2"
)

const (
	PROVIDER_GITHUB = "github"
)

type (
	// AccountReader type abstracts how data is retrieved from
	// different providers
	AccountReader interface {
		Email() (*string, error)
		Username() (*string, error)
		Name() (*string, error)
		Avatar() (*string, error)
		Data() (interface{}, error)
		Token(authCode string) (*oauth2.Token, error)
	}

	gitHubReader struct {
		clientID     string
		clientSecret string
		token        *oauth2.Token
		user         *github.User
	}
)

// MakeReader creates reader based on requested provider
func MakeReader(provider string) (AccountReader, error) {
	switch provider {
	case PROVIDER_GITHUB:
		return makeGitHubProviderUsingEnv(nil), nil
	}

	return nil, errors.Trace(errors.Errorf("Provider [%s] is not supported", provider))
}
