package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"time"
	"fmt"
	"net/http"
	"github.com/rs/cors"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	UserID      int `json:"user_id"`
	Creator     string `json:"creator"`
	Column      string `json:"column"`
	LastUpdated string `json:"lastUpdated"`
}

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal(err)
	}

	todos := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER,
		title TEXT,
		column TEXT,
		last_updated TEXT
	)`
	_, err = db.Exec(todos)
	if err != nil {
		log.Fatal(err)
	}

	todolists := `CREATE TABLE IF NOT EXISTS todolists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		slug TEXT UNIQUE
	)`
	_, err = db.Exec(todolists)
	if err != nil {
		log.Fatal(err)
	}

	users := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT UNIQUE
	)`
	_, err = db.Exec(users)
	if err != nil {
		log.Fatal(err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT
			todos.id,
			todos.title,
			todos.user_id,
			users.username AS creator,
			todos.column,
			todos.last_updated
		FROM todos
		LEFT JOIN users ON todos.user_id = users.id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.UserID, &todo.Creator, &todo.Column, &todo.LastUpdated); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("SELECT username FROM users WHERE id = ?", todo.UserID).Scan(&todo.Creator)
	if err != nil {

		http.Error(w, "Failed to fetch username: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and execute the INSERT statement
	stmt, err := db.Prepare("INSERT INTO todos (title, user_id, column, last_updated) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(todo.Title, todo.UserID, todo.Column, todo.LastUpdated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	todo.ID = int(id)


	// Return the newly created todo with the username included
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodoColumn(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE todos SET column = ?, last_updated = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.Column, todo.LastUpdated, todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := generateUsername()

	stmt, err := db.Prepare("INSERT INTO users (username) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user := User{ID: int(id), Username: username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// creates "nice" looking strings, similar to Figma's usernames or GitHub's repo names
func generateUsername() string {
	var adjectives = []string{"bright", "calm", "cool", "dark", "fast", "happy", "kind", "lucky", "quick", "shiny"}
	var nouns = []string{"cat", "dog", "fox", "lion", "panda", "tiger", "wolf", "zebra", "whale", "koala"}

	rand.Seed(time.Now().UnixNano())
	for {
		// add a random number to the tail to avoid collisions...
		username := fmt.Sprintf("%s-%s-%d", adjectives[rand.Intn(len(adjectives))], nouns[rand.Intn(len(nouns))], rand.Intn(100))

		// ...but still check for them!
		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
		if err != nil {
			log.Printf("Error checking username existence: %v", err)
			continue
		}

		if !exists {
			return username
		}
	}
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodos(w, r)
		} else if r.Method == http.MethodPost {
			addTodo(w, r)
		}
	})

	http.HandleFunc("/todos/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			updateTodoColumn(w, r)
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createUser(w, r)
		}
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
