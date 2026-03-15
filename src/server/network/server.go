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
<<<<<<< HEAD
			return
		case "UPDATE_STOCK":
			handleUpdateStock(reader, conn, db)
		case "UPDATE_PRICE":
			handleUpdatePrice(reader, conn, db)
		case "ORDER_HISTORY":
			handleOrderHistory(conn, db)
		case "LIST_PRODUCTS":
			handleListProducts(conn, db)
=======

		case "ADD":
			handleAddProduct(reader, conn, db)

>>>>>>> 4c29e995bd4066fdd3cc08d602078e744f09f378
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

func handleUpdateStock(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	id := readLine(reader)
	stock := readLine(reader)

	_, err := db.Exec("UPDATE productos SET cantidad = ? WHERE id = ?", stock, id)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not update stock")
		return
	}
	fmt.Fprintln(conn, "OK Stock updated successfully")
}

func handleUpdatePrice(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	id := readLine(reader)
	price := readLine(reader)

	_, err := db.Exec("UPDATE productos SET precio = ? WHERE id = ?", price, id)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not update price")
		return
	}
	fmt.Fprintln(conn, "OK Price updated successfully")
}

func handleOrderHistory(conn net.Conn, db *sql.DB) {
	query := `
		SELECT o.id, u.mail, p.nombre, o.cantidad, o.order_status 
		FROM ordenes o 
		JOIN usuarios u ON o.id_usuario = u.id 
		JOIN productos p ON o.id_producto = p.id
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not fetch order history")
		return
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id, cantidad int
		var mail, nombre, status string
		rows.Scan(&id, &mail, &nombre, &cantidad, &status)
		result = append(result, fmt.Sprintf("Order ID: %d - User: %s - Product: %s - Qty: %d - Status: %s", id, mail, nombre, cantidad, status))
	}

	if len(result) == 0 {
		fmt.Fprintln(conn, "OK No orders found")
		return
	}
	fmt.Fprintln(conn, "OK |"+strings.Join(result, "|"))
}

func handleListProducts(conn net.Conn, db *sql.DB) {
	rows, err := db.Query("SELECT id, nombre, cantidad, precio FROM productos")
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not fetch products")
		return
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id, cantidad int
		var name string
		var precio float64
		rows.Scan(&id, &name, &cantidad, &precio)
		result = append(result, fmt.Sprintf("ID: %d | Product: %s | Stock: %d | Price: $%.2f", id, name, cantidad, precio))
	}

	if len(result) == 0 {
		fmt.Fprintln(conn, "OK No products found")
		return
	}
	fmt.Fprintln(conn, "OK |"+strings.Join(result, "|"))
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
