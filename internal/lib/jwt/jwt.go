package jwt

import (
	"time"

	"github.com/Woland-prj/microtasks_sso/internal/domain/cerrors"
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

func ValidateToken(token string, secret string) (int64, error) {
	parsedToken, err := parseToken(token, secret)
	if err != nil {
		return 0, cerrors.NewInvalidTokenError(cerrors.TokenBadFormat)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if int64(claims["exp"].(float64)) <= time.Now().Unix() {
			return 0, cerrors.NewInvalidTokenError(cerrors.TokenExpired)
		}

		return int64(claims["id"].(float64)), nil
	}

	return 0, cerrors.NewInvalidTokenError(cerrors.TokenBadFormat)
}

func parseToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}