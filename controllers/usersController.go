package controllers

import (
	"database/sql"
	"fmt"
	"main/models"
	"main/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body parameters",
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{
		Email:    body.Email,
		Password: string(hashedPassword),
		FirstName: sql.NullString{
			String: body.FirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: body.LastName,
			Valid:  true,
		},
	}

	_, err = services.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create User: %s", err),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	authorizationCookie, _ := c.Request.Cookie("Authorization")
	if authorizationCookie != nil {
		c.Redirect(http.StatusFound, "/validate")
		return
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body parameters",
		})

		return
	}

	user, err := services.FindUser(models.User{Email: body.Email, Password: body.Password})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not get User",
		})

		return
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if bcryptErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if bcryptErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not generate JWT Token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*15, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "authenticated",
		"user":    user.(*models.User).Email,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
