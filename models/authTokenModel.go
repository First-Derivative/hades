package models

type AuthToken struct {
	ID                 int
	UserID             int
	AccessToken        string
	RefreshToken       string
	Invalidated        bool
	RefreshTokenExpiry int64
	CreatedAt          []byte
	UpdatedAt          []byte
}
