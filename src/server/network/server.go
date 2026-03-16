package network

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"server/admin"
	"server/auth"
	storeclient "server/storeClient"
	"strconv"
	"strings"
	"sync"
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
			return
		case "UPDATE_STOCK":
			handleUpdateStock(reader, conn, db)
		case "UPDATE_PRICE":
			handleUpdatePrice(reader, conn, db)
		case "ORDER_HISTORY":
			handleOrderHistory(conn, db)
		case "LIST_PRODUCTS":
			handleListProducts(conn, db)

		case "ADD":
			handleAddProduct(reader, conn, db)

		case "ADDTOCART":
			handleAddToCart(reader, conn, db)

		case "VIEWMYCART":
			handleViewMyCart(reader, conn, db)

		case "PLACEORDER":
			handlePlaceOrder(reader, conn, db)
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

	log.Printf("REGISTER - email: %s role: %s", mail, role)

	fmt.Fprintln(conn, "OK User registered successfully")
}

func handleLogin(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	mail := readLine(reader)
	pass := readLine(reader)

	user, err := auth.Login(mail, pass, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Invalid credentials")
		fmt.Fprintln(conn, "0")
		return
	}
	log.Printf("LOGIN - email: %s role: %s", user.Mail, user.Role)

	fmt.Printf("Client number ->( %s ) logged in as ( %s ) with email ( %s )\n", conn.RemoteAddr(), user.Role, user.Mail)

	fmt.Fprintln(conn, "OK "+user.Role)
	fmt.Fprintln(conn, user.ID)
}

func handleAddProduct(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	idUser := readLine(reader)
	name := readLine(reader)
	amountStr := readLine(reader)
	priceStr := readLine(reader)

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
	log.Printf("ADD PRODUCT - name: %s price: %.2f amount: %d", name, price, amount)

	fmt.Fprintln(conn, "product added by user "+idUser)

}

func handleAddToCart(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	idUser := readLine(reader)
	name := readLine(reader)
	amountStr := readLine(reader)

	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		fmt.Fprintln(conn, "ERROR invalid number")
		return
	}

	productName, err := storeclient.AddToCart(idUser, name, amount, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}
	log.Printf("ADD TO CART - user: %s product: %s amount: %d", idUser, productName, amount)

	fmt.Fprintln(conn, "OK "+productName+" added to cart")
}

func handleViewMyCart(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	idUser := readLine(reader)

	items, err := storeclient.ViewCart(idUser, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}
	log.Printf("VIEW CART - user: %s items: %d", idUser, len(items))

	fmt.Fprintln(conn, len(items))
	for _, item := range items {

		fmt.Fprintln(conn, fmt.Sprintf("%s|%d|%.2f|%s", item.Name, item.Cantidad, item.Price, item.Status))
	}
}

var mu sync.Mutex

func handlePlaceOrder(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	idUser := readLine(reader)

	mu.Lock()
	defer mu.Unlock()

	total, err := storeclient.PlaceOrder(idUser, db)
	if err != nil {
		fmt.Fprintln(conn, "ERROR "+err.Error())
		return
	}
	log.Printf("PLACE ORDER - user: %s total: $%.2f", idUser, total)

	fmt.Fprintln(conn, fmt.Sprintf("OK Order placed! Total: $%.2f", total))
}
func handleUpdateStock(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	id := readLine(reader)
	stock := readLine(reader)

	result, err := db.Exec("UPDATE products SET amount = ? WHERE name = ?", stock, id)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not update stock")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(conn, "ERROR checking update result")
		return
	}

	if rows == 0 {
		fmt.Fprintln(conn, "ERROR product not found")
		return
	}
	log.Printf("UPDATE STOCK - product: %s new stock: %s", id, stock)

	fmt.Fprintln(conn, "OK Stock updated successfully")
}

func handleUpdatePrice(reader *bufio.Reader, conn net.Conn, db *sql.DB) {
	id := readLine(reader)
	price := readLine(reader)

	result, err := db.Exec("UPDATE products SET price = ? WHERE name = ?", price, id)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not update price")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(conn, "ERROR checking update result")
		return
	}

	if rows == 0 {
		fmt.Fprintln(conn, "ERROR product not found")
		return
	}
	log.Printf("UPDATE PRICE - product: %s new price: %s", id, price)

	fmt.Fprintln(conn, "OK Price updated successfully")
}

func handleOrderHistory(conn net.Conn, db *sql.DB) {
	query := `
		SELECT o.id_user, id_product, p.name, o.cantidad, o.order_status 
		FROM orders o 
		JOIN users u ON o.id_user = u.id 
		JOIN products p ON o.id_product = p.id
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not fetch order history"+err.Error())
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
	log.Printf("ORDER HISTORY - fetched %d orders", len(result))

	fmt.Fprintln(conn, "OK |"+strings.Join(result, "|"))
}

func handleListProducts(conn net.Conn, db *sql.DB) {
	rows, err := db.Query("SELECT name, amount, price FROM products")
	if err != nil {
		fmt.Fprintln(conn, "ERROR Could not fetch products")
		return
	}
	defer rows.Close()

	type Product struct {
		Name   string
		Amount int
		Price  float64
	}

	var products []Product

	for rows.Next() {
		var p Product
		rows.Scan(&p.Name, &p.Amount, &p.Price)
		products = append(products, p)
	}
	log.Printf("LIST PRODUCTS - fetched %d products", len(products))

	fmt.Fprintln(conn, len(products))

	for _, p := range products {
		fmt.Fprintln(conn, fmt.Sprintf("%s|%d|%.2f", p.Name, p.Amount, p.Price))
	}
}
func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
