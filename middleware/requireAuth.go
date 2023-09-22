package middleware

import (
	"fmt"
	"main/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "unauthenticated",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {

			user_id := claims["sub"].(int)

			user, err := services.FindUserById(user_id)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status": "unauthenticated",
				})
			}

			c.Set("user", user)
		}

		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "unauthenticated",
		})
	}

}
