# E-commerce Simulator

A simple **Golang + MySQL** project designed to manage users, products, and orders. This guide explains how to install the necessary dependencies and run the project locally using the command line.

---

## 1. Prerequisites

Before getting started, ensure you have the following installed on your system:

* **Go**
* **MySQL Server**
* **Git**

You can verify your installations by running:

```bash
go version
mysql --version
git --version

```

---

## 2. Installing MySQL Server (Windows)

You can easily install **MySQL Server** using the `winget` package manager:

```bash
winget install Oracle.MySQL

```

Verify the installation:

```bash
mysql --version

```

> **Troubleshooting:** If your terminal does not recognize the `mysql` command after installation, you need to add the program's path to your system's **PATH** environment variable. You can do this by searching for "Edit the system environment variables" in the Windows search bar.

Start the MySQL service:

```bash
net start MySQL

```

*Note: If the above command fails, try running:* `net start MySQL84`

---

## 3. Cloning the Repository

Clone the project to your local machine:

```bash
git clone https://github.com/jfong088/GoEcommerceSimulator.git

```

---

## 4. Database Setup

First, log into MySQL to create the database:

```bash
mysql -u root -p

```

Execute the following SQL command to create the database, then exit:

```sql
CREATE DATABASE go_store;
exit

```

Next, from the root directory of the project, import the database schema. This will automatically execute the `database/schema.sql` file:

```bash
mysql -u root -p go_store < src/server/database/schema.sql

```

---

## 5. Installing Go Dependencies

Navigate to the server directory:

```bash
cd src/server

```

Install the MySQL driver for Go:

```bash
go get github.com/go-sql-driver/mysql

```

Clean up and synchronize your dependencies:

```bash
go mod tidy

```

---

## 6. Environment Configuration

Create an environment configuration file named `.env` inside the `src` folder.

Add your database credentials to the file:

```env
DB_USER=root
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=go_store

```

---

## 7. Running the Application

You will need two separate terminal windows to run the server and the client simultaneously.

**Step 1: Start the Server**

```bash
cd src/server
go run main.go

```

**Step 2: Start the Client**

```bash
cd src/client
go run client.go

```

---

## Administrator Features Showcase

* **Add Product:**
* **Update Stock:**
* **Update Price:**
* **Order History:**

## Client Features Showcase

* **Add Product to Cart:**
* **View Cart:**
* **Place Order (Checkout):**

---

## Future Improvements

1. **Product Search:** Instead of displaying the entire catalog, implement a feature allowing users to search for specific products.
2. **Client Purchase History:** Allow clients to view their past orders, including total amount spent, purchase date, and completed status.
3. **User Management:** Grant administrators the ability to ban specific email addresses or clear out users' carts.

## Additional Features

* **Data Persistence:** The application uses a SQL database to reliably store users, products, and orders. Even if the program is closed, all data remains accessible upon restart.
* **User Registration:** The database integration allows for secure user registration and login verification.
