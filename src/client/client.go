package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	conn         net.Conn
	serverReader *bufio.Reader
	userReader   *bufio.Reader

	idUser string
)

func main() {
	var err error

	conn, err = net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	serverReader = bufio.NewReader(conn)
	userReader = bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n-------- Main Menu -------")
		fmt.Println("1) Login")
		fmt.Println("2) Register")
		fmt.Println("3) Exit")
		fmt.Print("Choose an option: ")

		option := readInput()

		switch option {
		case "1":
			handleLogin()
		case "2":
			handleRegister()
		case "3":
			fmt.Println("Bye :)")
			return
		default:
			fmt.Println("Invalid option...")
		}
	}
}

func sendCommand(cmd string) {
	fmt.Fprintln(conn, cmd)
}

func readResponse() string {
	response, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Println("Connection closed...")
		os.Exit(1)
	}
	return strings.TrimSpace(response)
}

func readInput() string {
	input, _ := userReader.ReadString('\n')
	return strings.TrimSpace(input)
}

func handleRegister() {
	fmt.Print("Email: ")
	mail := readInput()

	fmt.Print("Password: ")
	pass := readInput()

	fmt.Print("Role (admin/client): ")
	role := readInput()

	sendCommand("REGISTER")
	sendCommand(mail)
	sendCommand(pass)
	sendCommand(role)

	response := readResponse()
	fmt.Println("Server:", response)
}

func handleLogin() {
	fmt.Print("Email: ")
	mail := readInput()

	fmt.Print("Password: ")
	pass := readInput()

	sendCommand("LOGIN")
	sendCommand(mail)
	sendCommand(pass)

	response := readResponse()
	fmt.Println("Server:", response)
	idUser = readResponse()

	if strings.HasPrefix(response, "OK admin") {
		adminMenu()
	} else if strings.HasPrefix(response, "OK client") {
		clientMenu()
	}
}

func adminMenu() {
	fmt.Println("\n=== Admin Panel ===")
	for {
		fmt.Println("\n--- Admin Menu ---")
		fmt.Println("1) Add product")
		fmt.Println("2) Update stock")
		fmt.Println("3) Update price")
		fmt.Println("4) View order history")
		fmt.Println("5) List products")
		fmt.Println("6) Logout")
		fmt.Print("Choose an option: ")

		option := readInput()
		switch option {
		case "1":
			handleAddProduct()
		case "2":
		case "3":
		case "4":
		case "5":
			sendCommand("LIST")
		case "6":
			sendCommand("LOGOUT")
			fmt.Println("Logged out...")
			return
		default:
			fmt.Println("Invalid option...")
		}
	}
}

func handleAddProduct() {
	fmt.Print("Product Name: ")
	name := readInput()

	fmt.Print("Product amount: ")
	amount := readInput()

	fmt.Print("Product price: ")
	price := readInput()
	sendCommand("ADD")
	sendCommand(idUser)
	sendCommand(name)
	sendCommand(amount)
	sendCommand(price)

	response := readResponse()
	fmt.Println("Server:", response)

}

func clientMenu() {
	fmt.Println("\n=== Client Panel ===")
	for {
		fmt.Println("\n--- Client Menu ---")
		fmt.Println("1) List products")
		fmt.Println("2) Add to cart")
		fmt.Println("3) View cart")
		fmt.Println("4) Place order")
		fmt.Println("5) Logout")
		fmt.Print("Choose an option: ")

		option := readInput()
		switch option {
		case "1":
			sendCommand("LIST")
		case "2":
			handleAddToCart()
		case "3":
			handleViewMyCart()
		case "4":
			handlePlaceOrder()
		case "5":
			sendCommand("LOGOUT")
			fmt.Println("Logged out")
			return
		default:
			fmt.Println("Invalid option...")
		}
	}
}

func handleAddToCart() {
	//handleListProducts()
	fmt.Print("Product name: ")
	name := readInput()

	fmt.Print("Amount: ")
	amount := readInput()
	sendCommand("ADDTOCART")
	sendCommand(idUser)
	sendCommand(name)
	sendCommand(amount)

	response := readResponse()
	fmt.Println("Server:", response)
}

func handleViewMyCart() {
	sendCommand("VIEWMYCART")
	sendCommand(idUser)

	countStr := readResponse()
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("ERROR:", countStr)
		return
	}

	if count == 0 {
		fmt.Println("Your cart is empty.")
		return
	}

	fmt.Println("\n Product               | Amount | Price    | Status")
	fmt.Println("--------------------------------------------------")
	for i := 0; i < count; i++ {
		line := readResponse()
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			fmt.Printf(" %-22s | %-3s | $%-7s | %s\n", parts[0], parts[1], parts[2], parts[3])
		}
	}

}
func handlePlaceOrder() {
	sendCommand("PLACEORDER")
	sendCommand(idUser)

	response := readResponse()
	fmt.Println("Server:", response)
}
