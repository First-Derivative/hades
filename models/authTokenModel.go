package models

type AuthToken struct {
	ID                 int
	UserID             int
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry int64
	Invalidated        bool
	CreatedAt          []byte
	UpdatedAt          []byte
}
