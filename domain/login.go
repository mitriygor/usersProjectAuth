package domain

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Login struct {
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Email     string `db:"email"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	return l.claimsForUser()
}

func (l Login) claimsForUser() AccessTokenClaims {
	return AccessTokenClaims{
		FirstName: l.FirstName,
		LastName:  l.LastName,
		Email:     l.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
