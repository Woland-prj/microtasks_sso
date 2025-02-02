package jwt

import (
	"time"

	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

const (
	_tokenTypeAuth    = "auth"
	_tokenTypeRefresh = "refresh"
)

func NewTokenPair(
	user *entities.User,
	app *entities.App,
	authDuration time.Duration,
	refreshDuration time.Duration,
) (*entities.JwtTokenPair, error) {
	authToken, err := newToken(user, app, authDuration, _tokenTypeAuth)
	if err != nil {
		return nil, err
	}

	refreshToken, err := newToken(user, app, authDuration, _tokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return &entities.JwtTokenPair{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func newToken(
	user *entities.User,
	app *entities.App,
	ttl time.Duration,
	tokenType string,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.UID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	var tokenString string
	var err error
	if tokenType == _tokenTypeRefresh {
		tokenString, err = token.SignedString([]byte(app.RefreshSecret))
	} else {
		tokenString, err = token.SignedString([]byte(app.AuthSecret))
	}

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
