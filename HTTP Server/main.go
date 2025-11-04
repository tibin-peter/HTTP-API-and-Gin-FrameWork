package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User is the data model for incoming JSON at /api/user
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// helloHandler handles GET /api/hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// prepare the JSON response
	response := map[string]string{"message": "Hello, World!"}

	// set the Content-Type header so the client knows we are sending JSON
	w.Header().Set("Content-Type", "application/json")

	// encode the response map as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// if encoding fails, return an internal server error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// userHandler handles POST /api/user
func userHandler(w http.ResponseWriter, r *http.Request) {
	// only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// ensure request body will be closed after we are done reading it
	defer r.Body.Close()

	// decode JSON body into a User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// print received user to the server terminal (for debugging / logging)
	fmt.Printf("Received user: %+v\n", user)

	// prepare a success response
	response := map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("User %s (age %d) received successfully!", user.Name, user.Age),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	// register URL paths with handler functions
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/user", userHandler)

	// log a startup message
	fmt.Println("Server is running on http://localhost:8080")

	// start the HTTP server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// ListenAndServe only returns an error if it fails to start or the server stops
		fmt.Printf("Server failed: %v\n", err)
	}
}
