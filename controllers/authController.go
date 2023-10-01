package controllers

import (
	"database/sql"
	"fmt"
	"main/models"
	"main/services"
	"main/utils"
	"net/http"
	"os"

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

	user, err := services.FindUser(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
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

	_, updateLoginErr := services.UpdateUserLoginAt(user.ID)
	if updateLoginErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Server Error: %s", updateLoginErr),
		})

		return
	}

	accessTokenString, _, refreshTokenErr := utils.GenerateWebToken(user.ID, true)
	if refreshTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate JWT Token",
		})
		return
	}

	refreshTokenString, refreshTokenExpiry, refreshTokenErr := utils.GenerateWebToken(user.ID, false)
	if refreshTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate JWT Token",
		})
		return
	}

	var authToken = models.AuthToken{
		UserID:             user.ID,
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		RefreshTokenExpiry: refreshTokenExpiry,
	}

	_, authTokenEror := services.CreateAuthToken(authToken)
	if authTokenEror != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Server Error: %s", authTokenEror),
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", accessTokenString, 3600, "/", "", gin.ReleaseMode == "relase", true)
	c.SetCookie("Authorization-Refresh", refreshTokenString, 3600*24*31, "/", "", gin.ReleaseMode == "relase", true)

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
	var user_id int

	user, err := c.Get("user")
	if err != true {
		tokenString, tokenStringError := c.Cookie("Authorization")

		if tokenStringError == nil {
			acessToken, _ := jwt.Parse(tokenString, func(acessToken *jwt.Token) (interface{}, error) {
				if _, ok := acessToken.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", acessToken.Header["alg"])
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if acessClaims, ok := acessToken.Claims.(jwt.MapClaims); ok && acessToken.Valid {
				user_id = int(acessClaims["sub"].(float64))
			}
		}

	} else {
		user_id = user.(*models.User).ID
	}

	if user_id > 0 {
		services.InvalidateAuthTokens(user_id)
	}

	utils.ClearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
