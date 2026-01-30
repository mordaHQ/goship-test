package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite" // Драйвер базы данных
)

// Структура пользователя
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	// 1. Подключаем базу данных (создаст файл data.db)
	db, err = sql.Open("sqlite", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Создаем таблицу, если её нет
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		created_at DATETIME
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Ошибка создания таблицы: %q: %s\n", err, sqlStmt)
	}
	fmt.Println("💾 База данных готова!")

	// 3. Маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html") // Отдаем админку
	})
	http.HandleFunc("/api/users", handleUsers) // API для данных

	fmt.Println("🚀 Сервер слушает :80")
	http.ListenAndServe(":80", nil)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET: Получить список
	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, name, email, created_at FROM users ORDER BY id DESC")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
			users = append(users, u)
		}
		json.NewEncoder(w).Encode(users)
		return
	}

	// POST: Создать нового
	if r.Method == "POST" {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		
		_, err := db.Exec("INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)", u.Name, u.Email, time.Now())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
