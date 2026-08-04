package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func initDb(log *log.Logger) (*sql.DB, error) {

	connStr := "postgresql://little_da_vinci_trial_signup_db_user:N9VqR3OdzZC2Sv3H2GppA3jltDtdmgUn@dpg-d9oauue417fc73eu88m0-a.singapore-postgres.render.com/little_da_vinci_trial_signup_db"

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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
`
}
