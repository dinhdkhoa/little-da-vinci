package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	const (
		port    = "8080"
		timeout = 30 * time.Second
	)

	logger := log.New(os.Stdout, "[perma-app] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting Permacomputing Web Application")

	db, err := initDb(logger)
	if err != nil {
		logger.Fatalf("cant connect to db", err)
	}

	migrationsql := getMigrationSql()
	db.Exec(migrationsql)
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("POST /sign-up", signUp(db))
	mux.HandleFunc("GET /sync-db/{id}", syncDb(db))

	handler := loggingMiddleware(logger)(panicRecoveryMiddleware(logger)(mux))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}

	logger.Printf("Server starting on port :%s", port)
	logger.Println("Press Ctrl+C to stop")

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func panicRecoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("PANIC: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Printf(
				"%s %s %s %d %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		})
	}
}

func signUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func syncDb(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
