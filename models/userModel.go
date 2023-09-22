package models

import (
	"database/sql"
)

type User struct {
	ID        int
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
}
