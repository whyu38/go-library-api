package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDatabase() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Tidak menemukan file .env, pastikan environment variable sudah diset")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal membuka driver database: ", err)
	}

	// Cek koneksi
	if err := DB.Ping(); err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	fmt.Println("Sukses terhubung ke Database!")
};