package main

import (
	"fmt"
	"server/database"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := "INSERT INTO usuarios (mail, pass, role) VALUES (?, ?, ?)"

	result, err := db.Exec(query, "test@mail.com", "123456", "admin")
	if err != nil {
		panic(err)
	}

	id, _ := result.LastInsertId()

	fmt.Println("User inserted with ID:", id)

	rows, err := db.Query("SELECT id, mail, role FROM usuarios")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var mail string
		var role string

		rows.Scan(&id, &mail, &role)

		fmt.Println(id, mail, role)
	}
}
