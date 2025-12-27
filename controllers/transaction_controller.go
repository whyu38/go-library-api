package controllers

import (
	"encoding/json"
	"go-library-api/config"
	"go-library-api/models"
	"net/http"
	"time"
)

// 1. PINJAM BUKU (Create)
func BorrowBook(w http.ResponseWriter, r *http.Request) {
	var req models.Borrowing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.BorrowDate == "" {
		req.BorrowDate = time.Now().Format("2006-01-02")
	} else {
		// Parse BorrowDate if provided
		parsed, err := time.Parse("2006-01-02", req.BorrowDate)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		req.BorrowDate = parsed.Format("2006-01-02")
	}

	// Transaksi Database (Safety)
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Cek Stok
	var stock int
	if err := tx.QueryRow("SELECT stock FROM books WHERE id = ?", req.BookID).Scan(&stock); err != nil {
		tx.Rollback()
		http.Error(w, "Buku tidak ditemukan", http.StatusNotFound)
		return
	}
	if stock <= 0 {
		tx.Rollback()
		http.Error(w, "Stok habis", http.StatusBadRequest)
		return
	}

	// Kurangi Stok
	if _, err := tx.Exec("UPDATE books SET stock = stock - 1 WHERE id = ?", req.BookID); err != nil {
		tx.Rollback()
		return
	}

	// Catat Pinjam
	req.BorrowCode = "BRW-" + time.Now().Format("20060102150405")
	res, err := tx.Exec("INSERT INTO borrowings (borrow_code, book_id, customer_id, borrow_date) VALUES (?, ?, ?, ?)", 
		req.BorrowCode, req.BookID, req.CustomerID, req.BorrowDate)
	
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()
	id, _ := res.LastInsertId()
	req.ID = int(id)
	json.NewEncoder(w).Encode(req)
}

// 2. KEMBALIKAN BUKU (Update/Edit special case)
func ReturnBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // ID Peminjaman (Borrowing ID)
	returnDate := time.Now().Format("2006-01-02")

	tx, _ := config.DB.Begin()

	// 1. Ambil BookID dari data peminjaman
	var bookID int
	err := tx.QueryRow("SELECT book_id FROM borrowings WHERE id = ?", id).Scan(&bookID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Data peminjaman tidak ditemukan", http.StatusNotFound)
		return
	}

	// 2. Update tanggal kembali
	_, err = tx.Exec("UPDATE borrowings SET return_date = ? WHERE id = ?", returnDate, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Gagal update status kembali", http.StatusInternalServerError)
		return
	}

	// 3. Kembalikan Stok Buku (+1)
	_, err = tx.Exec("UPDATE books SET stock = stock + 1 WHERE id = ?", bookID)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	json.NewEncoder(w).Encode(map[string]string{"message": "Buku berhasil dikembalikan"})
}

// 3. LIHAT SEMUA PEMINJAMAN (Read)
func GetBorrowings(w http.ResponseWriter, r *http.Request) {
	// Kita gunakan COALESCE(return_date, '') agar jika tanggal kembali masih NULL (belum kembali),
	// database mengembalikannya sebagai string kosong, bukan error.
	query := `SELECT id, borrow_code, book_id, customer_id, borrow_date, COALESCE(return_date, '') 
	          FROM borrowings`

	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var borrowings []models.Borrowing
	for rows.Next() {
		var b models.Borrowing
		// Scan data ke struct
		if err := rows.Scan(&b.ID, &b.BorrowCode, &b.BookID, &b.CustomerID, &b.BorrowDate, &b.ReturnDate); err != nil {
			continue
		}
		borrowings = append(borrowings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrowings)
}

// 4. HAPUS PEMINJAMAN (Delete)
func DeleteBorrowing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Mulai Transaksi (Penting!)
	// Kita harus cek: Jika buku SEDANG DIPINJAM (belum dikembalikan) lalu datanya dihapus,
	// maka stok buku harus dikembalikan (+1) agar tidak hilang misterius.
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// 1. Cek status peminjaman sebelum dihapus
	var bookID int
	var returnDate string // Bisa kosong kalau null
	
	// Gunakan COALESCE agar tidak error saat scan NULL
	err = tx.QueryRow("SELECT book_id, COALESCE(return_date, '') FROM borrowings WHERE id = ?", id).Scan(&bookID, &returnDate)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	// 2. Logika Stok: Jika 'returnDate' kosong, artinya buku masih di luar.
	// Karena transaksinya mau dihapus paksa, stoknya kita balikin dulu.
	if returnDate == "" {
		_, err = tx.Exec("UPDATE books SET stock = stock + 1 WHERE id = ?", bookID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Gagal mengembalikan stok", http.StatusInternalServerError)
			return
		}
	}

	// 3. Hapus data peminjaman
	_, err = tx.Exec("DELETE FROM borrowings WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Gagal menghapus data", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data peminjaman berhasil dihapus"})
}