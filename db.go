package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDb(log *log.Logger) (*sql.DB, error) {

	connStr := os.Getenv("DB_URI")
	if connStr == "" {
		log.Fatal("DB_URI environment variable not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB handle: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func getMigrationSql() string {
	return `CREATE TABLE IF NOT EXISTS registrations (
    id SERIAL PRIMARY KEY,
    title VARCHAR(20) NOT NULL,
    parent_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    birth_year INT NOT NULL,
    gender VARCHAR(20) NOT NULL,
    source VARCHAR(100),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (timezone('Asia/Ho_Chi_Minh', now()))
	);
`
}

func checkDuplicate(db *sql.DB, title, parentName, phone, studentName string, birthYear int, gender string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM registrations
        WHERE title = $1
          AND parent_name = $2
          AND phone = $3
          AND student_name = $4
          AND birth_year = $5
          AND gender = $6
    )`
	var exists bool
	err := db.QueryRow(query, title, parentName, phone, studentName, birthYear, gender).Scan(&exists)
	return exists, err
}

