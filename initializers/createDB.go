package initializers

import (
	"fmt"
)

func CreateDB() {
	showTablesQuery := "SHOW TABLES;"
	tables, tablesError := DB.Query(showTablesQuery)

	if tablesError != nil {
		panic(fmt.Sprintf("Database Error: %s", tablesError))
	}

	var header string
	hasUsers := false
	hasTokens := false

	for tables.Next() {
		tables.Scan(&header)
		switch header {
		case "users":
			hasUsers = true
		case "tokens":
			hasTokens = true
		}

		if hasUsers && hasTokens {
			tables.Close()
		}
	}

	if !hasUsers {
		usersQuery := "CREATE TABLE `users` (`id` int PRIMARY KEY AUTO_INCREMENT, `email` varchar(100) NOT NULL UNIQUE, `password` varchar(100) NOT NULL, `firstName` varchar(50), `lastName` varchar(50) );"
		_, usersError := DB.Query(usersQuery)
		if usersError != nil {
			fmt.Println(fmt.Sprintf("Error creating users table: %s", usersError))
		}
	}

	if !hasTokens {
		tokensQuery := "CREATE TABLE `users` (`id` int PRIMARY KEY AUTO_INCREMENT, `email` varchar(100) NOT NULL UNIQUE, `password` varchar(100) NOT NULL, `firstName` varchar(50), `lastName` varchar(50) );"

		_, tokensError := DB.Query(tokensQuery)
		if tokensError != nil {
			fmt.Println(fmt.Sprintf("Error creating tokens table: %s", tokensError))
		}
	}
}
