package auth

import (
	"database/sql"
	"fmt"
	"strings"
)

type User struct {
	ID   int
	Mail string
	Role string
}

func Register(mail, pass, role string, db *sql.DB) error {
	if mail == "" || pass == "" || role == "" {
		return fmt.Errorf("all fields are required")
	}

	role = strings.ToLower(role)
	if role != "admin" && role != "client" {
		return fmt.Errorf("role must be 'admin' or 'client'")
	}

	query := "INSERT INTO usuarios (mail, pass, role) VALUES (?, ?, ?)"
	_, err := db.Exec(query, mail, pass, role)
	if err != nil {

		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("email already registered...")
		}
		return fmt.Errorf("could not register user...")
	}

	return nil
}

func Login(mail, pass string, db *sql.DB) (*User, error) {
	if mail == "" || pass == "" {
		return nil, fmt.Errorf("email and password are required...")
	}

	query := "SELECT id, role FROM usuarios WHERE mail=? AND pass=?"
	row := db.QueryRow(query, mail, pass)

	var user User
	user.Mail = mail

	err := row.Scan(&user.ID, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}
