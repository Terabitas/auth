package oauth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"github.com/digitalocean/godo"
	"strings"
)

// Endpoint is DigitalOcean's OAuth 2.0 endpoint.
var DOEndpoint = oauth2.Endpoint{
	AuthURL:  "https://cloud.digitalocean.com/v1/oauth/authorize",
	TokenURL: "https://cloud.digitalocean.com/v1/oauth/token",
}

func makeDigitalOceanProvider(clientID, clientSecret string, token *oauth2.Token) *digitalOceanReader {
	return &digitalOceanReader{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
	}
}

func makeDigitalOceanProviderUsingEnv(token *oauth2.Token) *digitalOceanReader {
	clientID := os.Getenv("ND_DO_CLIENT_ID")
	clientSecret := os.Getenv("ND_DO_SECRET")
	requiredScope := os.Getenv("ND_DO_SCOPE")
	redirectURL := os.Getenv("ND_DO_REDIRECT_URL")

	return &digitalOceanReader{
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
		redirectURL:  redirectURL,
		scope: 	      strings.Split(requiredScope, ","),
	}
}

// Token reader
func (dor *digitalOceanReader) Token(authCode string) (*oauth2.Token, error) {
	if authCode == "" {
		return nil, errors.Trace(errors.Errorf("authCode is empty!"))
	}

	conf := &oauth2.Config{
		ClientID:     dor.clientID,
		ClientSecret: dor.clientSecret,
		Scopes:       dor.scope,
		Endpoint:     DOEndpoint,
		RedirectURL:  dor.redirectURL,
	}

	token, err := conf.Exchange(oauth2.NoContext, authCode)

	if err != nil {
		log.WithField("authCode", authCode).
			WithField("clientID", dor.clientID).
			WithField("clientSecret", dor.clientSecret).
			WithField("scope", strings.Join(dor.scope, ",")).Errorf("%s", err)
		return nil, errors.Trace(err)
	}

	dor.token = token

	log.WithField("authCode", authCode).Debugf("Got token, will expiry at [%s]", dor.token.Expiry)
	ts := oauth2.StaticTokenSource(dor.token)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := godo.NewClient(tc)
	dor.account, _, err = client.Account.Get()

	if dor.account.Email == ""{
		return nil, errors.Errorf("Could not access user email!scope ?")
	}

	if err != nil {
		log.WithField("authCode", authCode).Errorf("%s", err)
		return nil, errors.Trace(err)
	}
	log.WithField("authCode", authCode).WithField("user", dor.account.Email).Debugf("[%s]", dor.account.String())

	return dor.token, nil
}

// Email reader
func (dor *digitalOceanReader) Email() (*string, error) {

	if dor.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	if dor.account.Email != "" {
		return &dor.account.Email, nil
	}

	return nil, errors.Errorf("No email available!")
}

// Username reader
func (dor *digitalOceanReader) Username() (*string, error) {
	if dor.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return nil, nil
}

// Avatar reader
func (dor *digitalOceanReader) Avatar() (*string, error) {
	if dor.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return nil, nil
}

// Name reader
func (dor *digitalOceanReader) Name() (*string, error) {
	if dor.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return nil, nil
}

// Data reader
func (dor *digitalOceanReader) Data() (interface{}, error) {
	if dor.token == nil {
		return nil, errors.Errorf("Use Token() first to exchange authCode for valid token")
	}

	return dor.account, nil
}
