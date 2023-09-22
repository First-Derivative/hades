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

func FindUser(user models.User) (*models.User, error) {

	query := fmt.Sprintf("SELECT * FROM `users` WHERE email = \"%s\"", user.Email)

	res, err := initializers.DB.Query(query)
	if err != nil {
		return nil, err
	}

	users := models.User{}
	for res.Next() {
		var user models.User
		err := res.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName)
		if err != nil {
			log.Fatal("(GetProducts) res.Scan", err)
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
		err := res.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName)
		if err != nil {
			log.Fatal("(GetProducts) res.Scan", err)
		}
		users = user
	}

	return &users, nil

}
