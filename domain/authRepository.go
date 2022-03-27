package domain

import (
	"database/sql"
	"github.com/mitriygor/usersProjectLib/errors"
	"github.com/mitriygor/usersProjectLib/logger"
)

type AuthRepository interface {
	FindBy(email string, password string) (*Login, *errors.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errors.AppError)
	RefreshTokenExists(refreshToken string) *errors.AppError
}

type AuthRepositoryDb struct {
	client *sql.DB
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errors.AppError {
	sqlSelect := "select token from refresh_tokens where token = ?"
	err := d.client.QueryRow(sqlSelect, refreshToken)
	if err != nil {
		logger.Error("Unexpected DB error: " + err.Err().Error())
		return errors.UnexpectedError("Unexpected DB error")
	}
	return nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errors.AppError) {
	var appErr *errors.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}
	sqlInsert := "INSERT INTO refresh_tokens (token) VALUES ($1)"
	_, err := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("Unexpected DB error: " + err.Error())
		return "", errors.UnexpectedError("Unexpected DB error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) FindBy(email, password string) (*Login, *errors.AppError) {
	var login Login
	var firstName string
	var lastName string

	err := d.client.QueryRow("SELECT firstname, lastname from users WHERE email=$1 AND password=$2", email, password).Scan(&firstName, &lastName)

	if err != nil {
		logger.Error("Error during request from DB: " + err.Error())
		return nil, errors.UnexpectedError("Unexpected DB error")
	}

	login.FirstName = firstName
	login.LastName = lastName
	login.Email = email
	return &login, nil
}

func NewAuthRepository(client *sql.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
