package domain

import (
	"github.com/golang-jwt/jwt"
	"github.com/mitriygor/usersProjectLib/errors"
	"github.com/mitriygor/usersProjectLib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *errors.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errors.UnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errors.AppError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HMAC_SECRET))
	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errors.UnexpectedError("cannot generate refresh token")
	}
	return signedString, nil
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errors.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SECRET), nil
	})
	if err != nil {
		return "", errors.AuthenticationError("invalid or expired refresh token")
	}
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)

	return authToken.NewAccessToken()
}
