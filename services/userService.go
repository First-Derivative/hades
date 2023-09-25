package services

import (
	"database/sql"
	"fmt"
	"log"
	"main/initializers"
	"main/models"
)

func CreateUser(user models.User) (*sql.Rows, error) {
	var query string

	hasFirstName := user.FirstName.Valid
	hasLirstName := user.LastName.Valid

	if hasFirstName && hasLirstName {
		query = fmt.Sprintf("INSERT INTO `users` (email, password, firstName, lastName) VALUES (\"%s\", \"%s\", \"%s\", \"%s\")", user.Email, user.Password, user.FirstName.String, user.LastName.String)
	} else {
		query = fmt.Sprintf("INSERT INTO `users` (email, password) VALUES (\"%s\", \"%s\")", user.Email, user.Password)
	}

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateUserLogin(id int) (*sql.Rows, error) {
	query := fmt.Sprintf("UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id=\"%d\";", id)

	res, err := initializers.DB.Query(query)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func FindUser(email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("User email is empty")
	}

	query := fmt.Sprintf("SELECT * FROM `users` WHERE email = \"%s\"", email)

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	users := models.User{}
	for res.Next() {
		var user models.User
		err := res.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
		if err != nil {
			log.Fatal("(GetUser) res.Scan", err)
		}
		users = user
	}

	return &users, nil

}

func FindUserById(id int) (*models.User, error) {

	query := fmt.Sprintf("SELECT * FROM `users` WHERE id = \"%d\"", id)

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	users := models.User{}
	for res.Next() {
		var user models.User
		err := res.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
		if err != nil {
			log.Fatal("(GetUser) res.Scan", err)
		}
		users = user
	}

	return &users, nil

}
