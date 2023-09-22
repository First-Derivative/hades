package initializers

import "fmt"

func CreateDB() {
	query := "CREATE TABLE `users` (`id` int PRIMARY KEY AUTO_INCREMENT, `email` varchar(100) NOT NULL UNIQUE, `password` varchar(100) NOT NULL, `firstName` varchar(50), `lastName` varchar(50) )"

	_, err := DB.Query(query)
	if err != nil {
		fmt.Print("CreateDB returning null because `users` table already exists ")
	}
}
