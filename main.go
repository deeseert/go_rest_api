package main

import (
	"encoding/json" // For creating json
	"log"           // Loging errors
	"math/rand"     // For creating ids
	"net/http"      // creating web server
	"strconv"       // Stringify the ids

	// Gorilla Mux for creating Routers. See GitHub for installation
	"github.com/gorilla/mux"
)

// Create Book Struct here (Model - it's similar to class)
type Book struct {
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
	//Author has it's own struct
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Initialize books var as a slice Book struct
var books []Book

// Function to get all books - it's similar to nodejs where you have res and req
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Need the id. Get it using mux.Vars() and pass in the request "r"
	params := mux.Vars(r)
	// Then loop through books and find with id
	for _, item := range books {
		// if you find it
		if item.ID == params["id"] {
			// then do this for that particular
			json.NewEncoder(w).Encode(item)
			// Then return from this
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Create a variable called book and set it to the Book struct
	var book Book
	// When you need to pass in the body, you need to use Decode and not Encode
	_ = json.NewDecoder(r.Body).Decode(&book)
	// Then you need to create an id for your new book using math/rand
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	// Then append the new created book to the global variable books you created earlier
	books = append(books, book) //First argument is where, second argument is what
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	// it's a combination of delete and create functions
	// Part of delete function
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Slice it out and add a new one
			books = append(books[:index], books[index+1:]...)
			// and here you add the new one ()
			// Part of create function
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Initialize Router
	router := mux.NewRouter()

	// Mock data - @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "448743", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "543534", Title: "Book two", Author: &Author{Firstname: "Gio", Lastname: "Dina"}})

	// Create Router Handelers to establish Endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Run Server -log.Fatal() logs errors
	log.Fatal(http.ListenAndServe(":8000", router))
}
