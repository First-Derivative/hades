package main

import (
	"fmt"
	"main/controllers"
	"main/initializers"
	"main/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT")))
}
