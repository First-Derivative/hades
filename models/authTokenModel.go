package models

import (
	"database/sql"
)

type AuthToken struct {
	ID                 int
	User               int
	AccessToken        string
	RefreshTokenID     sql.NullString
	Invalidated        bool
	RefreshTokenExpiry []byte
	CreatedAt          []byte
	UpdatedAt          []byte
}

/*
Need to implement (inherit) custom struct based off: jwt.map_claims

type AccessTokenClaims struct {
	*jwt.MapClaims
}

type RefreshTokenClaims struct {
	*jwt.MapClaims
}

*/
