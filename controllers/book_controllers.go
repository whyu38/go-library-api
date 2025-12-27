package controllers

import (
	"encoding/json"
	"go-library-api/config"
	"go-library-api/models"
	"net/http"
	"strconv"
)

// 1. GET: Ambil Semua Buku
func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, title, author, stock FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Stock); err != nil {
			continue
		}
		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// 2. GET: Ambil Satu Buku Detail
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Fitur Go 1.22
	
	var b models.Book
	err := config.DB.QueryRow("SELECT id, title, author, stock FROM books WHERE id = ?", id).Scan(&b.ID, &b.Title, &b.Author, &b.Stock)
	if err != nil {
		http.Error(w, "Buku tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// 3. POST: Buat Buku Baru
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec("INSERT INTO books (title, author, stock) VALUES (?, ?, ?)", b.Title, b.Author, b.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	b.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

// 4. PUT: Edit Buku
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Logic update database
	_, err := config.DB.Exec("UPDATE books SET title = ?, author = ?, stock = ? WHERE id = ?", b.Title, b.Author, b.Stock, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Buku berhasil diupdate", "data": b})
}

// 5. DELETE: Hapus Buku
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// PENTING: Biasanya database akan menolak jika buku ini sedang dipinjam (Foreign Key Constraint).
	// Itu bagus untuk keamanan data.
	_, err := config.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Gagal menghapus (mungkin buku sedang dipinjam)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Buku berhasil dihapus"})
}