package storeclient

import (
	"database/sql"
	"fmt"
	"sync"
)

var mutex sync.Mutex

func AddToCart(idUser, name string, cantidad int, db *sql.DB) (string, error) {

	mutex.Lock()
	defer mutex.Unlock()

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
	_, err = db.Exec("UPDATE products SET amount = amount - ? WHERE id = ?", cantidad, idProduct)
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

	var total float64
	err := db.QueryRow(`
        SELECT SUM(p.price * o.cantidad) 
        FROM orders o 
        JOIN products p ON o.id_product = p.id 
        WHERE o.id_user = ? AND o.order_status = 'in cart'
    `, idUser).Scan(&total)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec("UPDATE orders SET order_status = 'completed' WHERE id_user = ? AND order_status = 'in cart'", idUser)
	if err != nil {
		return 0, err
	}

	return total, nil
}
