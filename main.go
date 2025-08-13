package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	TimeCreate string `json:"time_create"`
}

var users []User
var currentID int = 1

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func isEmailValid(email string) bool {
	emailCheck := regexp.MustCompile(`^[a-zA-Z0-9]+(?:[@][a-zA-Z0-9]+)(?:[.][a-zA-Z0-9]+)+$`)
	return emailCheck.MatchString(email)
}


func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if !isEmailValid(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.Username == user.Username {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if u.Email == user.Email {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}
	user.ID = currentID
	currentID++
	user.TimeCreate = time.Now().Format(time.RFC3339)

	users = append(users, user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func findUserByID(id int) (*User, int) {
	for i, u := range users {
		if u.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/users/put/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, _ := findUserByID(id)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var input User
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Email != "" && !isEmailValid(input.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	
	for _, u := range users {
		if u.ID != user.ID && input.Username != "" && u.Username == input.Username {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if u.ID != user.ID && input.Email != "" && u.Email == input.Email {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/users/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, idx := findUserByID(id)
	if idx == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	users = append(users[:idx], users[idx+1:]...)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func main() {
	http.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getUsers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			createUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/put/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			updateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			deleteUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
