package services

import (
	"database/sql"
	"fmt"
	"main/initializers"
	"main/models"
)

func CreateUser(user models.User) (*sql.Rows, error) {
	var query string

	hasFirstName := user.FirstName != ""
	hasLirstName := user.LastName != ""

	if hasFirstName && hasLirstName {
		query = fmt.Sprintf("INSERT INTO `users` (email, password, firstName, lastName) VALUES (\"%s\", \"%s\", \"%s\", \"%s\")", user.Email, user.Password, user.FirstName, user.LastName)
	} else {
		query = fmt.Sprintf("INSERT INTO `users` (email, password) VALUES (\"%s\", \"%s\")", user.Email, user.Password)
	}

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	return res, nil
}
