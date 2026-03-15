package storeclient

import (
	"database/sql"
	"fmt"
)

func AddToCart(idUser, name string, cantidad int, db *sql.DB) (string, error) {
	var idProduct int
	var stock int

	err := db.QueryRow("SELECT id, amount FROM products WHERE LOWER(name) = LOWER(?)", name).Scan(&idProduct, &stock)
	if err != nil {
		return "", fmt.Errorf("product '%s' not found", name)
	}

	if cantidad <= 0 {
		return "", fmt.Errorf("invalid amount")
	}

	if cantidad > stock {
		return "", fmt.Errorf("insufficient stock, only %d available", stock)
	}

	_, err = db.Exec(
		"INSERT INTO orders (id_user, id_product, cantidad, order_status) VALUES (?, ?, ?, ?)",
		idUser, idProduct, cantidad, "in cart",
	)
	if err != nil {
		return "", err
	}

	return name, nil
}

type CartItem struct {
	Name     string
	Cantidad int
	Price    float64
	Status   string
}

func ViewCart(idUser string, db *sql.DB) ([]CartItem, error) {
	rows, err := db.Query(`
        SELECT p.name, o.cantidad, p.price, o.order_status 
        FROM orders o JOIN products p ON o.id_product = p.id 
        WHERE o.id_user = ? AND o.order_status = 'in cart'
    `, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		rows.Scan(&item.Name, &item.Cantidad, &item.Price, &item.Status)
		items = append(items, item)
	}

	return items, nil
}

func PlaceOrder(idUser string, db *sql.DB) (float64, error) {

	var count int
	db.QueryRow("SELECT COUNT(*) FROM orders WHERE id_user = ? AND order_status = 'in cart'", idUser).Scan(&count)
	if count == 0 {
		return 0, fmt.Errorf("cart is empty")
	}

	rows, err := db.Query(`
        SELECT p.name, p.amount, o.cantidad FROM orders o JOIN products p ON o.id_product = p.id WHERE o.id_user = ? AND o.order_status = 'in cart'`, idUser)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var stock, cantidad int
		rows.Scan(&name, &stock, &cantidad)
		if cantidad > stock {
			return 0, fmt.Errorf("insufficient stock for '%s', only %d available", name, stock)
		}
	}
	var total float64
	db.QueryRow(`
        SELECT SUM(p.price * o.cantidad) FROM orders o JOIN products p ON o.id_product = p.id WHERE o.id_user = ? AND o.order_status = 'in cart'`, idUser).Scan(&total)

	_, err = db.Exec(`
        UPDATE products p JOIN orders o ON o.id_product = p.id SET p.amount = p.amount - o.cantidad WHERE o.id_user = ? AND o.order_status = 'in cart'`, idUser)
	if err != nil {
		return 0, err
	}
	_, err = db.Exec("UPDATE orders SET order_status = 'completed' WHERE id_user = ? AND order_status = 'in cart'", idUser)
	if err != nil {
		return 0, err
	}

	return total, nil
}
