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
	authToken, err := newToken(user, app.ID, app.AuthSecret, authDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := newToken(user, app.ID, app.RefreshSecret, refreshDuration)
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
	appId int64,
	secret string,
	ttl time.Duration,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.UID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = appId

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
