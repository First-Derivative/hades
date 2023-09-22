package initializers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectToDB() {

	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}

	DB = db
}
