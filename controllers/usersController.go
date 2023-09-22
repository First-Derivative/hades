package controllers

import (
	"fmt"
	"main/models"
	"main/services"
	"net/http"

	"github.com/gin-gonic/gin"
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
		Email:     body.Email,
		Password:  string(hashedPassword),
		FirstName: body.FirstName,
		LastName:  body.LastName,
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
