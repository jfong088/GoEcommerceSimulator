package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
			fmt.Print("Product ID: ")
			id := readInput()
			fmt.Print("New Stock: ")
			stock := readInput()
			sendCommand("UPDATE_STOCK")
			sendCommand(id)
			sendCommand(stock)
			fmt.Println("Server:", readResponse())
		case "3":
			fmt.Print("Product ID: ")
			id := readInput()
			fmt.Print("New Price: ")
			price := readInput()
			sendCommand("UPDATE_PRICE")
			sendCommand(id)
			sendCommand(price)
			fmt.Println("Server:", readResponse())
		case "4":
			sendCommand("ORDER_HISTORY")
			response := readResponse()
			fmt.Println("Server:\n" + strings.ReplaceAll(response, "|", "\n"))
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
<<<<<<< HEAD
			sendCommand("LIST_PRODUCTS")
			response := readResponse()
			fmt.Println("Server:\n" + strings.ReplaceAll(response, "|", "\n"))
=======
			sendCommand("LIST")
>>>>>>> 4c29e995bd4066fdd3cc08d602078e744f09f378
		case "2":
		case "3":
		case "4":
		case "5":
			sendCommand("LOGOUT")
			fmt.Println("Logged out")
			return
		default:
			fmt.Println("Invalid option...")
		}
	}
}
