package oauth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"golang.org/x/oauth2"
	oauthGithub "golang.org/x/oauth2/github"
)

func makeGitHubProvider(clientID, clientSecret string, token *oauth2.Token) *gitHubReader {
	return &gitHubReader{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
	}
}

func makeGitHubProviderUsingEnv(token *oauth2.Token) *gitHubReader {
	clientID := os.Getenv("ND_GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("ND_GITHUB_SECRET")

	return &gitHubReader{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
	}
}

// Token reader
func (gr *gitHubReader) Token(authCode string) (*oauth2.Token, error) {
	if authCode == "" {
		return nil, errors.Trace(errors.Errorf("authCode is empty!"))
	}

	conf := &oauth2.Config{
		ClientID:     gr.clientID,
		ClientSecret: gr.clientSecret,
		Scopes:       []string{},
		Endpoint:     oauthGithub.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, authCode)

	if err != nil {
		log.WithField("authCode", authCode).Errorf("%s", err)
		return nil, errors.Trace(err)
	}

	gr.token = token

	log.WithField("authCode", authCode).Debugf("Got token, will expiry at [%s]", gr.token.Expiry)
	ts := oauth2.StaticTokenSource(gr.token)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	gr.user, _, err = client.Users.Get("")

	if err != nil {
		log.WithField("authCode", authCode).Errorf("%s", err)
		return nil, errors.Trace(err)
	}
	log.WithField("authCode", authCode).WithField("user", gr.user.Email).Debugf("[%s]", gr.user.String())

	return gr.token, nil
}

// Email reader
func (gr *gitHubReader) Email() (*string, error) {

	if gr.token == nil {
		return nil, errors.Trace(errors.Errorf("Use Token() first to exchange authCode for valid token"))
	}

	return gr.user.Email, nil
}

// Username reader
func (gr *gitHubReader) Username() (*string, error) {
	if gr.token == nil {
		return nil, errors.Trace(errors.Errorf("Use Token() first to exchange authCode for valid token"))
	}

	return gr.user.Login, nil
}

// Avatar reader
func (gr *gitHubReader) Avatar() (*string, error) {
	if gr.token == nil {
		return nil, errors.Trace(errors.Errorf("Use Token() first to exchange authCode for valid token"))
	}

	return gr.user.AvatarURL, nil
}

// Name reader
func (gr *gitHubReader) Name() (*string, error) {
	if gr.token == nil {
		return nil, errors.Trace(errors.Errorf("Use Token() first to exchange authCode for valid token"))
	}

	return gr.user.Name, nil
}

// Data reader
func (gr *gitHubReader) Data() (interface{}, error) {
	if gr.token == nil {
		return nil, errors.Trace(errors.Errorf("Use Token() first to exchange authCode for valid token"))
	}

	return gr.user, nil
}
