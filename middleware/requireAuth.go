package middleware

import (
	"fmt"
	"main/models"
	"main/services"
	"main/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {

	accessTokenString, acessTokenStringError := c.Cookie("Authorization")

	if acessTokenStringError != nil {

		refreshTokenString, refreshTokenStringError := c.Cookie("Authorization-Refresh")
		if refreshTokenStringError != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthenticated",
			})

			return
		}

		refreshToken, _ := jwt.Parse(refreshTokenString, func(refreshToken *jwt.Token) (interface{}, error) {
			if _, ok := refreshToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", refreshToken.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {

			hasAuthToken, findAuthTokenErr := services.FindAuthToken(refreshTokenString)
			if findAuthTokenErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": findAuthTokenErr,
				})
				return
			}

			user_id := int(refreshClaims["sub"].(float64))
			expiryClaim := refreshClaims["exp"].(float64)

			newAccessTokenString, _, webTokenError := utils.GenerateWebToken(user_id, true)
			if webTokenError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Could not generate JWT Token",
				})
				return
			}

			if hasAuthToken != true {
				utils.ClearAuthCookies(c)

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status": "unauthenticated",
				})

				return
			}

			user, findUserError := services.FindUserById(user_id)

			if findUserError != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status": "unauthenticated",
				})

				return
			}

			var authToken = models.AuthToken{
				UserID:             user_id,
				AccessToken:        newAccessTokenString,
				RefreshToken:       refreshTokenString,
				RefreshTokenExpiry: int64(expiryClaim),
			}

			services.InvalidateAndResignAuthTokens(user_id, authToken)

			c.Set("user", user)
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", newAccessTokenString, 3600, "/", "", gin.ReleaseMode == "relase", true)
			c.Next()
		}
	} else {
		acessToken, _ := jwt.Parse(accessTokenString, func(acessToken *jwt.Token) (interface{}, error) {
			if _, ok := acessToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", acessToken.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if acessClaims, ok := acessToken.Claims.(jwt.MapClaims); ok && acessToken.Valid {
			if float64(time.Now().Unix()) < acessClaims["exp"].(float64) {
				user_id := acessClaims["sub"].(float64)

				user, err := services.FindUserById(int(user_id))

				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"status": "unauthenticated",
					})

					return
				}

				c.Set("user", user)
				c.Next()
			}

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthenticated",
			})
		}
	}
}
