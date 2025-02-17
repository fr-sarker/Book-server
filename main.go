package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// ***************************  Credentials  *******************************************************************
var secretKey = []byte("your-256-bit-secret")

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

// Mock database of users
var users = map[string]string{
	"admin": "password123",
}

// Login function

func Login(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	storedPassword, ok := users[cred.Username]
	if !ok || storedPassword != cred.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: cred.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unable to retrieve claims", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome, %s!", claims.Username)})
}

// ******************************************************  Model of Books   ***************************************

// create a Book struct and a slice of pointer to the Book struct:
type Book struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Author           string `json:"author"`
	PublishedDate    string `json:"published_date"`
	OriginalLanguage string `json:"original_language"`
}

var books = []*Book{
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

// The listBooks function returns the books as a function instead of directly using the books variable.
func listBooks() []*Book {
	return books
}

// Given an ID as input, this function searches for a book with a matching ID in the books collection. If found, it returns a pointer to that book. Otherwise, it returns nil.

func getBook(id string) *Book {
	for _, book := range books {
		if book.ID == id {
			return book
		}
	}
	return nil
}

// This function adds a new book to the books collection. It takes a Book struct as input and appends a pointer to that struct to the books slice
func storeBook(book Book) {
	books = append(books, &book)
}

// This function removes abook from the books collection based on the provided id
func deleteBook(id string) *Book {
	for i, book := range books {
		if book.ID == id {
			deleteBook := books[i]
			books = append(books[:i], books[i+1:]...)
			return deleteBook
		}
	}
	return nil
}

// This function updates a book in the books collection
func updateBook(id string, bookUpdate Book) *Book {
	for i, book := range books {
		if book.ID == id {
			books[i] = &bookUpdate
			return book
		}
	}
	return nil
}

// **********************************   BookHandler      *******************************

// encode the list of books as json data.
func (b Book) ListBooks(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(listBooks())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// The GetBooks handler reads the requested book ID from the URL using the chi.URLParam function
// This ID is passed to the getBook function.
func (b Book) GetBooks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book := getBook(id)
	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// added new book
func (b Book) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	storeBook(book)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// update book
func (b Book) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedBook := updateBook(id, book)

	if updatedBook == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(updatedBook)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// delete book
func (b Book) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deletedBook := deleteBook(id)
	if deletedBook == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func BookRoutes() chi.Router {
	r := chi.NewRouter()
	bookHandler := Book{}
	r.Get("/", bookHandler.ListBooks)
	r.Post("/", bookHandler.CreateBook)
	r.Get("/{id}", bookHandler.GetBooks)
	r.Put("/{id}", bookHandler.UpdateBook)
	r.Delete("/{id}", bookHandler.DeleteBook)
	return r
}

//*****************************************  Main Function  ***********************************************

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/login", Login)
	r.Post("/logout", Logout)
	r.Group(func(r chi.Router) {
		r.Use(Authenticate) // Apply the Authenticate middleware to this group
		r.Get("/protected", ProtectedHandler)
		r.Mount("/books", BookRoutes())
	})
	http.ListenAndServe(":3000", r)
}
