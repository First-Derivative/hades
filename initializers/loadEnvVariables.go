package initializers

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		if arg == "release" || arg == "--release" {
			gin.SetMode(gin.ReleaseMode)
		}
	}
}
