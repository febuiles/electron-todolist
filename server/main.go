package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type TodoList struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Slug   string `json:"slug"`
}

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	UserID      int    `json:"user_id"`
	TodoListID  int    `json:"todolist_id"`
	Creator     string `json:"creator"`
	Column      string `json:"column"`
	LastUpdated string `json:"lastUpdated"`
}

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	LastUsedTodolistID int    `json:"lastUsedTodolistId"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	_, err = db.Exec(TodosSchema)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(TodoListsSchema)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(UsersSchema)
	if err != nil {
		log.Fatal(err)
	}
}

func getTodolist(w http.ResponseWriter, r *http.Request) {
	todolistID := r.URL.Path[len("/todolists/"):]
	if todolistID == "" {
		http.Error(w, "Missing todolist ID", http.StatusBadRequest)
		return
	}

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
		WHERE todos.todolist_id = ?
	`, todolistID)
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

func createTodolist(w http.ResponseWriter, r *http.Request) {
	var todolist TodoList
	if err := json.NewDecoder(r.Body).Decode(&todolist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)", todolist.UserID).Scan(&userExists)
	if err != nil || !userExists {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	todolist.Slug = generateTodolistSlug()

	stmt, err := db.Prepare("INSERT INTO todolists (user_id, slug) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(todolist.UserID, todolist.Slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	todolist.ID = int(id)

	// update the user to set this as the last used list
	_, err = db.Exec("UPDATE users SET last_used_todolist_id = ? WHERE id = ?", todolist.ID, todolist.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todolist)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow(`
        SELECT id, username, last_used_todolist_id
        FROM users
        WHERE id = ?`, id).Scan(&user.ID, &user.Username, &user.LastUsedTodolistID)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var todolistExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM todolists WHERE id = ?)", todo.TodoListID).Scan(&todolistExists)
	if err != nil || !todolistExists {
		http.Error(w, "Failed to validated todolist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT username FROM users WHERE id = ?", todo.UserID).Scan(&todo.Creator)
	if err != nil {
		http.Error(w, "Failed to fetch username: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO todos (title, user_id, todolist_id, column, last_updated) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(todo.Title, todo.UserID, todo.TodoListID, todo.Column, todo.LastUpdated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	todo.ID = int(id)

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

func createUser(w http.ResponseWriter, _ *http.Request) {
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

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]
	if id == "" {
		http.Error(w, "Missing TODO ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare delete statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Failed to delete TODO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// creates "nice" looking strings, similar to Figma's usernames or GitHub's repo names
func generateUsername() string {
	var adjectives = []string{"bright", "calm", "cool", "dark", "fast", "happy", "kind", "lucky", "quick", "shiny"}
	var nouns = []string{"cat", "dog", "fox", "lion", "panda", "tiger", "wolf", "zebra", "whale", "koala"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		username := fmt.Sprintf("%s-%s-%d",
			adjectives[r.Intn(len(adjectives))],
			nouns[r.Intn(len(nouns))],
			r.Intn(100))

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

// generates sharing slugs inspired by Google Meet codes
func generateTodolistSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var segmentLengths = [3]int{3, 4, 3}

	for {
		var code string
		for i, length := range segmentLengths {
			for j := 0; j < length; j++ {
				num := rand.Intn(len(charset))
				code += string(charset[num])
			}
			if i < len(segmentLengths)-1 {
				code += "-"
			}
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM todolists WHERE slug = ?)", code).Scan(&exists)
		if err != nil {
			log.Printf("Error checking slug existence: %v", err)
			continue
		}

		if !exists {
			return code
		}
	}
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/_ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "PONG")
	})

	http.HandleFunc("/todolists/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTodolist(w, r)
		} else if r.Method == http.MethodGet {
			getTodolist(w, r)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTodo(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTodo(w, r)
		}
	})

	http.HandleFunc("/todos/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			updateTodoColumn(w, r)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createUser(w, r)
		} else if r.Method == http.MethodGet {
			getUser(w, r)
		}
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	log.Println("Server is running.")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
