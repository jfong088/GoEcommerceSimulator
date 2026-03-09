# Go Store

Proyecto simple en **Go + MySQL** para manejar:

* Usuarios
* Productos
* Órdenes

Este README explica cómo **instalar las dependencias y correr el proyecto localmente usando solo comandos**.

---

# 1. Requisitos

Instalar:

* Go
* MySQL Server
* Git

Verificar instalación:

```bash
go version
mysql --version
git --version
```

---

# 2. Instalar MySQL Server

## Windows

Instalar **MySQL Server** usando winget:

```bash
winget install Oracle.MySQL
```

Verificar instalación:

```bash
mysql --version
```

Iniciar el servicio:

```bash
net start MySQL
```

Si no funciona:

```bash
net start MySQL80
```

---

## Linux (Ubuntu / Debian)

```bash
sudo apt update
sudo apt install mysql-server
```

Iniciar servicio:

```bash
sudo systemctl start mysql
```

---

## Mac

Instalar con Homebrew:

```bash
brew install mysql
```

Iniciar servicio:

```bash
brew services start mysql
```

---

# 3. Clonar el repositorio

```bash
git clone https://github.com/jfong088/GoEcommerceSimulator.git
```

---

# 4. Crear las tablas con el schema

Desde la raíz del proyecto ejecutar:

```bash
mysql -u root -p go_store < src/server/database/schema.sql
```

Esto ejecutará automáticamente el archivo:

```
database/schema.sql
```
---

# 5. Instalar dependencias de Go

Instalar el driver de MySQL:

```bash
go get github.com/go-sql-driver/mysql
```

Limpiar dependencias:

```bash
go mod tidy
```

---

# 6. Configurar conexión a la base de datos

Editar el archivo:

```
database/connection.go
```

Configurar tus credenciales:

```go
user := "root"
password := "password"
host := "127.0.0.1"
port := "3306"
dbname := "go_store"
```

