package dto

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/mitriygor/usersProjectAuth/domain"
)

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"token"`
}

func (r RefreshTokenRequest) IsAccessTokenValid() *jwt.ValidationError {

	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SECRET), nil
	})
	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			return vErr
		}
	}
	return nil
}
