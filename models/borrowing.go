package models

type Borrowing struct {
	ID     			int       `json:"id"`
	BorrowCode 		string    `json:"borrow_code"`
	BookID 			int       `json:"book_id"`
	CustomerID 		int       `json:"customer_id"`
	BorrowDate 		string `json:"borrow_date"`
	ReturnDate 		string `json:"return_date"`
}