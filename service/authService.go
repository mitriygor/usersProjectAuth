package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/mitriygor/usersProjectAuth/domain"
	"github.com/mitriygor/usersProjectAuth/dto"
	"github.com/mitriygor/usersProjectLib/errors"
	"github.com/mitriygor/usersProjectLib/logger"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errors.AppError)
	Verify(urlParams map[string]string) *errors.AppError
	Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errors.AppError)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (s DefaultAuthService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errors.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {

			var appErr *errors.AppError
			if appErr = s.repo.RefreshTokenExists(request.RefreshToken); appErr != nil {
				return nil, appErr
			}

			var accessToken string
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errors.AuthenticationError("Error invalid token")
	}
	return nil, errors.AuthenticationError("Error cannot generate a new access token until the current one expires")
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errors.AppError) {
	var appErr *errors.AppError
	var login *domain.Login

	if login, appErr = s.repo.FindBy(req.Email, req.Password); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errors.AppError {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errors.AuthorizationError(err.Error())
	} else {
		if jwtToken.Valid {
			return nil
		} else {
			return errors.AuthorizationError("Error invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SECRET), nil
	})
	if err != nil {
		logger.Error("Error during parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo}
}
