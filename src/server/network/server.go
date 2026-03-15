package network

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"server/admin"
	"server/auth"
	"strconv"
	"strings"
)

func HandleClient(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			return
		}

		command := strings.TrimSpace(line)

		switch command {
		case "REGISTER":
			handleRegister(reader, conn, db)

		case "LOGIN":
			handleLogin(reader, conn, db)

		case "LOGOUT":
			fmt.Printf("Client logged out: %s\n", conn.RemoteAddr())

		case "ADD":
			handleAddProduct(reader, conn, db)

		default:
			fmt.Fprintln(conn, "ERROR Unknown command")
		}
	}
}

func handleRegister(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	mail := readLine(reader)
	pass := readLine(reader)
	role := readLine(reader)

	err := auth.Register(mail, pass, role, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}
	fmt.Printf("Client number -> ( %s ) registered with email ( %s ) and is now a ( %s )\n", conn.RemoteAddr(), mail, role)

	fmt.Fprintln(conn, "OK User registered successfully")
}

func handleLogin(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	mail := readLine(reader)
	pass := readLine(reader)

	user, err := auth.Login(mail, pass, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Invalid credentials")
		return
	}
	fmt.Printf("Client number ->( %s ) logged in as ( %s ) with email ( %s )\n", conn.RemoteAddr(), user.Role, user.Mail)

	fmt.Fprintln(conn, "OK "+user.Role)
	fmt.Fprintln(conn, user.ID)
}

func handleAddProduct(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	idUser := readLine(reader)
	name := readLine(reader)
	priceStr := readLine(reader)
	amountStr := readLine(reader)

	err := admin.ValidateAdmin(idUser, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Fprintln(conn, "ERROR price must be a number")
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Fprintln(conn, "ERROR amount must be a number")
		return
	}

	err = admin.AddProduct(name, price, amount, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}

	fmt.Fprintln(conn, "product added by user "+idUser)

}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
