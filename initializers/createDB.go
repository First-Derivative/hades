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
		case "auth_tokens":
			hasTokens = true
		}

		if hasUsers && hasTokens {
			tables.Close()
		}
	}

	if !hasUsers {
		usersQuery := "CREATE TABLE `users` (`id` int PRIMARY KEY AUTO_INCREMENT, `email` varchar(100) NOT NULL UNIQUE, `password` varchar(100) NOT NULL, `firstName` varchar(50), `lastName` varchar(50),  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, `last_login_at` DATETIME);"
		_, usersError := DB.Query(usersQuery)
		if usersError != nil {
			fmt.Println(fmt.Sprintf("Error creating users table: %s", usersError))
		} else {
			fmt.Println("DB: users table created")
		}
	}

	if !hasTokens {
		tokensQuery := "CREATE TABLE `auth_tokens` (`id` int PRIMARY KEY AUTO_INCREMENT, `user_id` INT NOT NULL, `access_token` VARCHAR(384) NOT NULL, `refresh_token` VARCHAR(384) NOT NULL, `refresh_expiry` DATETIME NOT NULL, `invalidated` BOOLEAN NOT NULL DEFAULT FALSE, `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, INDEX `idx_refresh_token` (refresh_token), KEY user_id_idx (user_id));"

		_, tokensError := DB.Query(tokensQuery)
		if tokensError != nil {
			fmt.Println(fmt.Sprintf("Error creating auth_tokens table: %s", tokensError))
		} else {
			fmt.Println("DB: auth_tokens table created")
		}
	}
}
