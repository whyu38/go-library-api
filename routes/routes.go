package routes

import (
	"go-library-api/controllers"
	"net/http"
)

func InitRoutes(mux *http.ServeMux) {
	// Routes Buku
	mux.HandleFunc("GET /books", controllers.GetBooks)
	mux.HandleFunc("POST /books", controllers.CreateBook)
	mux.HandleFunc("PUT /books/{id}", controllers.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", controllers.DeleteBook)

	// Routes Customer
	mux.HandleFunc("GET /customers", controllers.GetCustomers)
	mux.HandleFunc("POST /customers", controllers.CreateCustomer)
	mux.HandleFunc("PUT /customers/{id}", controllers.UpdateCustomer)
	mux.HandleFunc("DELETE /customers/{id}", controllers.DeleteCustomer)

	// Routes Peminjaman
	mux.HandleFunc("GET /borrowings", controllers.GetBorrowings)
	mux.HandleFunc("POST /borrowings", controllers.BorrowBook)
	mux.HandleFunc("PUT /borrowings/{id}", controllers.ReturnBook)
	mux.HandleFunc("DELETE /borrowings/{id}", controllers.DeleteBorrowing)
}