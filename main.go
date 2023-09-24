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
	initializers.CreateDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3005"
	}

	address := fmt.Sprintf("127.0.0.1:%s", PORT)
	fmt.Println("Hades server running...")
	r.Run(address)

	defer initializers.DB.Close()
}
