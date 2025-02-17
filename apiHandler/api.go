package apiHandler

import (
	"appscode/fr-sarker/golang-chi-crud-api/authHandler"
	"appscode/fr-sarker/golang-chi-crud-api/dataHandler"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// encode the list of books as json data.
func ListBook(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(dataHandler.ListBooks())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// The GetBooks handler reads the requested book ID from the URL using the chi.URLParam function
// This ID is passed to the getBook function.
func GetBooks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book := dataHandler.GetBook(id)
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
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book dataHandler.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dataHandler.StoreBook(book)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// update book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var book dataHandler.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedBook := dataHandler.UpdateBook(id, book)

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
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deletedBook := dataHandler.DeleteBook(id)
	if deletedBook == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func BookRoutes() chi.Router {
	r := chi.NewRouter()
	//bookHandler := dataHandler.Book{}
	r.Get("/", ListBook)
	r.Post("/", CreateBook)
	r.Get("/{id}", GetBooks)
	r.Put("/{id}", UpdateBook)
	r.Delete("/{id}", DeleteBook)
	return r
}
func RunServer(Port int) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Group(func(r chi.Router) {
		r.Use(authHandler.Authenticate) // Apply the Authenticate middleware to this group
		r.Get("/protected", authHandler.ProtectedHandler)
		r.Mount("/books", BookRoutes())
	})
	http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}
