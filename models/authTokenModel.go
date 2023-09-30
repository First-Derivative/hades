package models

type AuthToken struct {
	ID                 int
	UserID             int
	AccessToken        string
	RefreshToken       string
	Invalidated        bool
	RefreshTokenExpiry []byte
	CreatedAt          []byte
	UpdatedAt          []byte
}
