# Go Web Application with PostgreSQL - Development Setup Guide

This guide will walk you through setting up a development environment for a Go-based web application that uses PostgreSQL as the database. We'll cover installation, environment setup, and database configuration for both Windows and Linux platforms.

## Prerequisites

Before you begin, make sure you have the following prerequisites installed:

- Go (Golang): https://golang.org/dl/
- PostgreSQL: https://www.postgresql.org/download/
- Git: https://git-scm.com/downloads

## Windows Setup

### 1. Install Go

1. Download the latest Go installer for Windows from the official website: https://golang.org/dl/
2. Run the installer and follow the installation instructions.

### 2. Install Git

1. Download the latest Git for Windows from the official website: https://git-scm.com/downloads
2. Run the installer and follow the installation instructions.

### 3. Install PostgreSQL

1. Download the PostgreSQL installer for Windows from the official website: https://www.postgresql.org/download/windows/
2. Run the installer and follow the installation instructions. Remember the password you set for the PostgreSQL superuser during installation.

### 4. Clone the Project

Open a command prompt and run the following command to clone the project repository:

```shell
git clone https://github.com/shakg/sslash.git
```

### 5. Configure PostgreSQL

1. Open the PostgreSQL command prompt.
2. Create a new database for the project:

```sql
CREATE DATABASE your_database;
```

3. Create a new user and grant them access to the database:

```sql
CREATE USER your_username WITH PASSWORD 'your_password';
ALTER ROLE your_username SET client_encoding TO 'utf8';
ALTER ROLE your_username SET default_transaction_isolation TO 'read committed';
ALTER ROLE your_username SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE your_database TO your_username;
```

### 6. Update Connection URL

In your project code, update the PostgreSQL connection URL with your credentials. You can find this in the `main.go` file:

```go
postgresURL := "user=your_username dbname=your_database sslmode=disable"
```

### 7. Run the Application

Navigate to the project directory in the command prompt and run the application:

```shell
go run main.go
```

The application should now be running on `http://localhost:8080`.

## Linux Setup

### 1. Install Go

1. Download the latest Go binary distribution for Linux from the official website: https://golang.org/dl/
2. Extract the downloaded archive to a directory of your choice (e.g., `/usr/local`).

```shell
tar -C /usr/local -xzf go1.XX.X.linux-amd64.tar.gz
```

3. Add Go to your PATH by appending the following line to your `~/.bashrc` or `~/.zshrc` file:

```shell
export PATH=$PATH:/usr/local/go/bin
```

4. Run `source ~/.bashrc` or `source ~/.zshrc` to apply the changes immediately.

### 2. Install Git

Use your distribution's package manager to install Git:

For Ubuntu/Debian:

```shell
sudo apt-get update
sudo apt-get install git
```

For Fedora:

```shell
sudo dnf install git
```

### 3. Install PostgreSQL

Use your distribution's package manager to install PostgreSQL:

For Ubuntu/Debian:

```shell
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
```

For Fedora:

```shell
sudo dnf install postgresql-server
```

### 4. Configure PostgreSQL

1. Start the PostgreSQL service:

For Ubuntu/Debian:

```shell
sudo service postgresql start
```

For Fedora:

```shell
sudo systemctl start postgresql
```

2. Create a new database and user:

```shell
sudo -u postgres psql
```

```sql
CREATE DATABASE your_database;
CREATE USER your_username WITH PASSWORD 'your_password';
ALTER ROLE your_username SET client_encoding TO 'utf8';
ALTER ROLE your_username SET default_transaction_isolation TO 'read committed';
ALTER ROLE your_username SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE your_database TO your_username;
```

Exit the PostgreSQL prompt:

```sql
\q
```

### 5. Clone the Project

Open a terminal and run the following command to clone the project repository:

```shell
git clone https://github.com/yourusername/your-project.git
```

### 6. Update Connection URL

In your project code, update the PostgreSQL connection URL with your credentials. You can find this in the `main.go` file:

```go
postgresURL := "user=your_username dbname=your_database sslmode=disable"
```

### 7. Run the Application

Navigate to the project directory in the terminal and run the application:

```shell
go run main.go
```

The application should now be running on `http://localhost:8080`.

---

You have successfully set up your development environment for the Go web application with PostgreSQL. You can access the application in your web browser at `http://localhost:8080`.