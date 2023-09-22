package main

import (
	"fmt"
	"main/controllers"
	"main/initializers"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/users", controllers.Signup)

	r.Run(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT")))

}
