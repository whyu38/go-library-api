# Go Library API 📚

A robust RESTful API for a Library Management System built with **Go (Golang) v1.22+** and **MySQL**. This project demonstrates clean architecture, database transactions, and raw SQL queries without ORM.

## 🚀 Tech Stack

* **Language:** Go (Golang) 1.22+
* **Database:** MySQL
* **Routing:** `net/http` (Standard Library with Go 1.22 `ServeMux`)
* **Driver:** `go-sql-driver/mysql`
* **Architecture:** MVC (Model, View/JSON, Controller)

## ✨ Features

* **Books Management:** Create, Read, Update, Delete (CRUD) books.
* **Customer Management:** Manage library members.
* **Borrowing System (Transactions):**
    * Borrow a book (Auto-decrement stock).
    * Return a book (Auto-increment stock).
    * **Atomic Transactions:** Ensures data consistency between `loans` and `books` tables.
* **CORS Enabled:** Ready to be consumed by Frontend (React/Vue).

## 📂 Project Structure

```bash
go-library-api/
├── config/             # Database connection logic
├── controllers/        # Request handlers & business logic
├── models/             # Data structures (Structs)
├── routes/             # API Route definitions
├── .env                # Environment variables (DB credentials)
├── main.go             # Entry point
└── go.mod              # Dependencies
