package domain

import (
	"time"

	"os"

	"github.com/nildev/auth/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"
)

type (
	// Session type
	Session struct {
		Id    string `bson:"_id" json:"id"`
		Token string `bson:"token" json:"token"`
	}
)

// MakeSession constructor
func MakeSession(email string) (*Session, error) {
	signKey := os.Getenv("ND_SIGN_KEY")
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["email"] = email
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(signKey))

	if err != nil {
		return nil, err
	}

	return &Session{
		Id:    email,
		Token: tokenString,
	}, nil
}
