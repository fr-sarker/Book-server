package dataHandler

import "github.com/golang-jwt/jwt/v5"

// Credentials struct to store username and password
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct to store JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Book struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Author           string `json:"author"`
	PublishedDate    string `json:"published_date"`
	OriginalLanguage string `json:"original_language"`
}

var Books = []*Book{
	{
		ID:               "1",
		Title:            "Atomic Habits",
		Author:           "John Doe",
		PublishedDate:    "2022-01-01",
		OriginalLanguage: "English",
	},
	{
		ID:               "2",
		Title:            "The Things You Can See Only When You Slow Down",
		Author:           "Haemin Sunim",
		PublishedDate:    "2022-04-01",
		OriginalLanguage: "English",
	},
}

// The listBooks function returns the Books as a function instead of directly using the Books variable.
func ListBooks() []*Book {
	return Books
}

// Given an ID as input, this function searches for a book with a matching ID in the Books collection. If found, it returns a pointer to that book. Otherwise, it returns nil.

func GetBook(id string) *Book {
	for _, book := range Books {
		if book.ID == id {
			return book
		}
	}
	return nil
}

// This function adds a new book to the Books collection. It takes a Book struct as input and appends a pointer to that struct to the Books slice
func StoreBook(book Book) {
	Books = append(Books, &book)
}

// This function removes abook from the Books collection based on the provided id
func DeleteBook(id string) *Book {
	for i, book := range Books {
		if book.ID == id {
			deleteBook := Books[i]
			Books = append(Books[:i], Books[i+1:]...)
			return deleteBook
		}
	}
	return nil
}

// This function updates a book in the Books collection
func UpdateBook(id string, bookUpdate Book) *Book {
	for i, book := range Books {
		if book.ID == id {
			Books[i] = &bookUpdate
			return book
		}
	}
	return nil
}
