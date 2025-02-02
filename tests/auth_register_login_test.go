package tests

import (
	"testing"
	"time"

	ssov1 "github.com/Woland-prj/microtasks_protos/gen/go/sso"
	"github.com/Woland-prj/microtasks_sso/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId       = 0
	appId            = 1
	appAuthSecret    = "test_app_auth_secret"
	appRefreshSecret = "test_app_refresh_secret"

	passDefaultLen = 10
	loginTimeDelta = 1
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUid())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appId,
	})

	require.NoError(t, err)

	loginTime := time.Now()

	authTokenParsed, err := parseToken(respLogin.GetAuthToken(), appAuthSecret)
	require.NotEmpty(t, authTokenParsed)
	require.NoError(t, err)

	authClaims, ok := authTokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUid(), int64(authClaims["id"].(float64)))
	assert.Equal(t, email, authClaims["email"].(string))
	assert.Equal(t, int64(appId), int64(authClaims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL.Auth).Unix(), authClaims["exp"].(float64), loginTimeDelta)

	refreshTokenParsed, err := parseToken(respLogin.GetRefreshToken(), appRefreshSecret)
	require.NotEmpty(t, refreshTokenParsed)
	require.NoError(t, err)

	refreshClaims, ok := refreshTokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUid(), int64(refreshClaims["id"].(float64)))
	assert.Equal(t, email, refreshClaims["email"].(string))
	assert.Equal(t, int64(appId), int64(refreshClaims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL.Refresh).Unix(), refreshClaims["exp"].(float64), loginTimeDelta)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func parseToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
