package services

import (
	"database/sql"
	"fmt"
	"log"
	"main/initializers"
	"main/models"
	"time"
)

func CreateAuthToken(authToken models.AuthToken) (*sql.Rows, error) {
	var query string

	refreshTokenExpiryUnixTime := time.Unix(authToken.RefreshTokenExpiry, 0)
	refreshTokenExpiry := refreshTokenExpiryUnixTime.Format("2006-01-02 15:04:05")
	query = fmt.Sprintf("INSERT INTO `auth_tokens` (user_id, access_token, refresh_token, invalidated, refresh_expiry) VALUES (\"%d\", \"%s\", \"%s\", false, \"%s\")", authToken.UserID, authToken.AccessToken, authToken.RefreshToken, refreshTokenExpiry)

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func FindAuthToken(refresh_token string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM `auth_tokens` WHERE refresh_token= \"%s\" AND invalidated != true", refresh_token)

	res, err := initializers.DB.Query(query)
	if err != nil {
		return false, err
	}

	token := models.AuthToken{}
	var refreshTokenExpiry []uint8

	for res.Next() {
		var authToken models.AuthToken
		err := res.Scan(&authToken.ID, &authToken.UserID, &authToken.AccessToken, &authToken.RefreshToken, &refreshTokenExpiry, &authToken.Invalidated, &authToken.CreatedAt, &authToken.UpdatedAt)
		if err != nil {
			log.Fatal("(`GetToken`) res.Scan", err)
		}
		token = authToken
	}

	if token.ID == 0 {
		return false, nil
	}

	return true, nil
}

func InvalidateAuthTokens(user_id int) error {

	transaction, err := initializers.DB.Begin()

	if err != nil {
		return err
	}

	defer transaction.Rollback()

	{
		_, err = transaction.Exec("UPDATE auth_tokens SET invalidated = true WHERE user_id = ?", user_id)
		if err != nil {
			return err
		}
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func InvalidateAndResignAuthTokens(user_id int, newAuthToken models.AuthToken) error {
	transaction, err := initializers.DB.Begin()

	if err != nil {
		return err
	}

	defer transaction.Rollback()

	{
		_, err = transaction.Exec("UPDATE auth_tokens SET invalidated = true WHERE user_id = ?", user_id)
		if err != nil {
			return err
		}
	}

	{
		refreshTokenExpiryUnixTime := time.Unix(newAuthToken.RefreshTokenExpiry, 0)
		refreshTokenExpiry := refreshTokenExpiryUnixTime.Format("2006-01-02 15:04:05")
		_, err = transaction.Exec("INSERT INTO auth_tokens (user_id, access_token, refresh_token, refresh_expiry) VALUES (?, ?, ?, ?)",
			newAuthToken.UserID, newAuthToken.AccessToken, newAuthToken.RefreshToken, refreshTokenExpiry)
		if err != nil {
			return err
		}
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}
