package auth // import "github.com/nildev/auth"
import (
	log "github.com/Sirupsen/logrus"
	"github.com/nildev/auth/domain"
	"github.com/nildev/auth/oauth"
)

// Auth user and return access token
// @method GET
// @path /auth/{provider}/{authCode}
func Auth(provider string, authCode string) (token string, err error) {
	reader, err := oauth.MakeReader(provider)
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	t, err := reader.Token(authCode)
	if t == nil {
		log.Errorf("%s", err)
		return "", err
	}
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	email, err := reader.Email()
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	session, err := domain.MakeSession(*email)
	return session.Token, err
}
