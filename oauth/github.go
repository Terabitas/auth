package oauth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"golang.org/x/oauth2"
	oauthGithub "golang.org/x/oauth2/github"
	"strings"
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
	requiredScope := os.Getenv("ND_GITHUB_SCOPE")

	return &gitHubReader{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
		scope: 	      strings.Split(requiredScope, ","),
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
		Scopes:       gr.scope,
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

	gr.emails, _, err = client.Users.ListEmails(nil)
	if err != nil {
		log.WithField("authCode", authCode).Errorf("%s", err)
		return nil, errors.Trace(err)
	}

	if len(gr.emails) == 0 && gr.user.Email == nil{
		return nil, errors.Errorf("Could not access user email!scope ?")
	}

	if err != nil {
		log.WithField("authCode", authCode).Errorf("%s", err)
		return nil, errors.Trace(err)
	}
	log.WithField("authCode", authCode).WithField("user", gr.user.Email).Debugf("[%s]", gr.user.String())
	log.WithField("authCode", authCode).WithField("emails", gr.user.Email).Debugf("[%+v]",gr.emails)

	return gr.token, nil
}

// Email reader
func (gr *gitHubReader) Email() (*string, error) {

	if gr.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	if gr.user.Email != nil {
		return gr.user.Email, nil
	}

	for _, eml := range gr.emails {
		if eml.Verified == nil && eml.Primary == nil {
			continue
		}

		if *eml.Verified && *eml.Primary && eml.Email != nil {
			return eml.Email, nil
		}
	}

	return nil, errors.Errorf("No email available!")
}

// Username reader
func (gr *gitHubReader) Username() (*string, error) {
	if gr.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return gr.user.Login, nil
}

// Avatar reader
func (gr *gitHubReader) Avatar() (*string, error) {
	if gr.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return gr.user.AvatarURL, nil
}

// Name reader
func (gr *gitHubReader) Name() (*string, error) {
	if gr.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return gr.user.Name, nil
}

// Data reader
func (gr *gitHubReader) Data() (interface{}, error) {
	if gr.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return gr.user, nil
}
