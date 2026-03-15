package admin

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

func ValidateAdmin(id string, db *sql.DB) error {
	query := "SELECT role FROM users WHERE id=?"
	row := db.QueryRow(query, id)

	var role string

	err := row.Scan(&role)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	if role != "admin" {
		return fmt.Errorf("User not admin")
	}
	return nil
}

func AddProduct(name string, price float64, amount int, db *sql.DB) error {

	query := "INSERT INTO products (name, price, amount) VALUES (?, ?, ?)"

	_, err := db.Exec(query, name, price, amount)
	if err != nil {
		return err
	}

	return nil
}
