package main

import (
	"fmt"
	"go-library-api/config"
	"go-library-api/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDatabase()

	// 2. Setup Router
	mux := http.NewServeMux()
	
    // Panggil fungsi dari folder routes
	routes.InitRoutes(mux)

	// 3. Middleware CORS (Agar bisa diakses semua Frontend)
	handler := enableCORS(mux)

	// 4. Jalankan Server
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

	fmt.Printf("Server API Library berjalan di http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// Fungsi CORS Universal (Penting untuk Frontend!)
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Izinkan akses dari mana saja (*)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}