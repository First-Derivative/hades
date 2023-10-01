package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateWebToken(user_id int, acess_token bool) (string, int64, error) {
	var tokenExpiry int64
	var ACESS_TOKEN_EXPIRY = time.Now().Add(time.Hour).Unix()
	var REFRESH_TOKEN_EXPIRY = time.Now().Add(time.Hour * 24 * 31).Unix()

	if acess_token {
		tokenExpiry = ACESS_TOKEN_EXPIRY
	} else {
		tokenExpiry = REFRESH_TOKEN_EXPIRY
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user_id,
		"exp": tokenExpiry,
		"iat": time.Now(),
		"iss": "hades",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", -1, err
	}

	return tokenString, tokenExpiry, nil

}

func ClearAuthCookies(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", gin.ReleaseMode == "relase", true)
	c.SetCookie("Authorization-Refresh", "", -1, "/", "", gin.ReleaseMode == "relase", true)
}
